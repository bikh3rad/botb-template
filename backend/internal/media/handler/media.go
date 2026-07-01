package handler

import (
	"application/internal/media/biz"
	"application/internal/media/dto"
	"application/internal/service"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Upload guards. maxUploadBytes caps the request body defensively; the biz layer
// enforces the finer per-kind limits (10 MiB images / 200 MiB videos).
const (
	maxUploadBytes    int64 = 210 << 20 // 200 MiB video limit + headroom
	multipartMemBytes int64 = 32 << 20  // buffer up to 32 MiB in memory, rest to disk
)

type media struct {
	logger *slog.Logger
	tracer trace.Tracer
	mux    *http.ServeMux
	uc     biz.UsecaseMedia
}

var _ service.Handler = (*media)(nil)

// NewMedia creates the media HTTP handler.
func NewMedia(logger *slog.Logger, mux *http.ServeMux, uc biz.UsecaseMedia) *media {
	return &media{
		logger: logger.With("layer", "MediaHandler"),
		tracer: otel.Tracer("MediaHandler"),
		mux:    mux,
		uc:     uc,
	}
}

// RegisterHandler mounts the media routes under /apis/media/v1.
func (h *media) RegisterHandler(_ context.Context) error {
	h.mux.HandleFunc("POST /apis/media/v1/uploads", h.upload)
	h.mux.HandleFunc("GET /apis/media/v1/media/{id}", h.get)
	h.mux.HandleFunc("GET /apis/media/v1/media", h.listByOwner)

	return nil
}

// upload handles a multipart file upload.
//
//	@Summary		Upload a media file
//	@Description	Upload an image (jpeg/png/webp ≤10MB) or video (mp4/webm ≤200MB) and store it in object storage.
//	@Tags			Media
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file		formData	file	true	"Media file"
//	@Param			owner_type	formData	string	true	"Owner object type, e.g. competition"
//	@Param			owner_id	formData	string	true	"Owner object UUID"
//	@Param			position	formData	int		false	"Ordering position (default 0)"
//	@Success		201			{object}	dto.MediaResp
//	@Failure		400			{object}	dto.ErrorResponse	"Bad Request"
//	@Failure		413			{object}	dto.ErrorResponse	"Payload Too Large"
//	@Failure		415			{object}	dto.ErrorResponse	"Unsupported Media Type"
//	@Failure		500			{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router			/apis/media/v1/uploads [post]
func (h *media) upload(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "Upload")
	ctx := r.Context()

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadBytes)

	if err := r.ParseMultipartForm(multipartMemBytes); err != nil {
		logger.WarnContext(ctx, "failed to parse multipart form", "error", err)
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	ownerID, err := uuid.Parse(r.FormValue("owner_id"))
	if err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, errors.New("invalid owner_id")), w)

		return
	}

	position, _ := strconv.Atoi(r.FormValue("position"))

	file, header, err := r.FormFile("file")
	if err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, errors.New("missing file")), w)

		return
	}
	defer file.Close()

	in := biz.UploadInput{
		OwnerType:   r.FormValue("owner_type"),
		OwnerID:     ownerID,
		ContentType: header.Header.Get("Content-Type"),
		Size:        header.Size,
		Position:    position,
		Reader:      file,
	}

	stored, err := h.uc.Upload(ctx, in)
	if err != nil {
		logger.ErrorContext(ctx, "upload failed", "error", err)
		dto.HandleError(err, w)

		return
	}

	writeJSON(w, http.StatusCreated, dto.ToMediaResp(stored), logger, ctx)
}

// get returns media metadata plus a presigned read URL.
//
//	@Summary		Get media metadata
//	@Description	Fetch a media record by ID, including a time-limited presigned URL.
//	@Tags			Media
//	@Produce		json
//	@Param			id	path		string	true	"Media UUID"
//	@Success		200	{object}	dto.MediaResp
//	@Failure		400	{object}	dto.ErrorResponse	"Bad Request"
//	@Failure		404	{object}	dto.ErrorResponse	"Not Found"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router			/apis/media/v1/media/{id} [get]
func (h *media) get(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "Get")
	ctx := r.Context()

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	res, err := h.uc.Get(ctx, id)
	if err != nil {
		logger.ErrorContext(ctx, "get failed", "error", err)
		dto.HandleError(err, w)

		return
	}

	resp := dto.ToMediaResp(res.Media)
	resp.URL = res.URL

	writeJSON(w, http.StatusOK, resp, logger, ctx)
}

// listByOwner returns all media for an owner object.
//
//	@Summary		List media for an owner
//	@Description	List every media object attached to an owner (owner_type + owner_id query params).
//	@Tags			Media
//	@Produce		json
//	@Param			owner_type	query		string	true	"Owner object type, e.g. competition"
//	@Param			owner_id	query		string	true	"Owner object UUID"
//	@Success		200			{object}	dto.MediaListResp
//	@Failure		400			{object}	dto.ErrorResponse	"Bad Request"
//	@Failure		500			{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router			/apis/media/v1/media [get]
func (h *media) listByOwner(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "ListByOwner")
	ctx := r.Context()

	ownerID, err := uuid.Parse(r.URL.Query().Get("owner_id"))
	if err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, errors.New("invalid owner_id")), w)

		return
	}

	items, err := h.uc.ListByOwner(ctx, r.URL.Query().Get("owner_type"), ownerID)
	if err != nil {
		logger.ErrorContext(ctx, "list failed", "error", err)
		dto.HandleError(err, w)

		return
	}

	writeJSON(w, http.StatusOK, dto.ToMediaListResp(items), logger, ctx)
}

func writeJSON(w http.ResponseWriter, status int, body any, logger *slog.Logger, ctx context.Context) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		logger.ErrorContext(ctx, "failed to encode response", "error", err)
	}
}
