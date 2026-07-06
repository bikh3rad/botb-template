package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"application/internal/competition/biz"
	"application/internal/competition/dto"
	"application/internal/competition/entity"
	"application/internal/service"
	"application/pkg/audit"
	"application/pkg/middlewares"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type competition struct {
	logger *slog.Logger
	tracer trace.Tracer
	mux    *http.ServeMux
	uc     biz.UsecaseCompetition
	auth   *middlewares.JWTAuth
	audit  *audit.Recorder
}

var _ service.Handler = (*competition)(nil)

// NewCompetition creates the competition HTTP handler.
func NewCompetition(
	logger *slog.Logger,
	mux *http.ServeMux,
	uc biz.UsecaseCompetition,
	auth *middlewares.JWTAuth,
	recorder *audit.Recorder,
) *competition {
	return &competition{
		logger: logger.With("layer", "CompetitionHandler"),
		tracer: otel.Tracer("CompetitionHandler"),
		mux:    mux,
		uc:     uc,
		auth:   auth,
		audit:  recorder,
	}
}

// RegisterHandler mounts public reads and admin mutations. The admin group is
// JWT-guarded here too (defense in depth) — the gateway also guards it, but a
// service reached directly on its own port must not be unprotected.
func (h *competition) RegisterHandler(_ context.Context) error {
	admin := func(fn http.HandlerFunc) http.HandlerFunc {
		return middlewares.MultipleMiddleware(fn, h.auth.RequireAdmin)
	}

	// Public.
	h.mux.HandleFunc("GET /apis/competition/v1/competitions", h.list)
	h.mux.HandleFunc("GET /apis/competition/v1/competitions/{id}", h.get)
	// Admin.
	h.mux.HandleFunc("POST /apis/competition/v1/admin/competitions", admin(h.create))
	h.mux.HandleFunc("PUT /apis/competition/v1/admin/competitions/{id}", admin(h.update))
	h.mux.HandleFunc("DELETE /apis/competition/v1/admin/competitions/{id}", admin(h.delete))

	return nil
}

// list returns competitions, optionally filtered by status.
//
//	@Summary		List competitions
//	@Description	Public list of competitions, optionally filtered by status.
//	@Tags			Competitions
//	@Produce		json
//	@Param			status	query		string	false	"Filter by status: draft, live or closed"
//	@Success		200		{object}	dto.CompetitionListResp
//	@Failure		400		{object}	dto.ErrorResponse	"Bad Request"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router			/apis/competition/v1/competitions [get]
func (h *competition) list(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "List")
	ctx := r.Context()

	var filter biz.ListFilter

	if s := r.URL.Query().Get("status"); s != "" {
		status := entity.Status(s)
		filter.Status = &status
	}

	items, err := h.uc.List(ctx, filter)
	if err != nil {
		logger.ErrorContext(ctx, "list failed", "error", err)
		dto.HandleError(err, w)

		return
	}

	writeJSON(ctx, w, http.StatusOK, dto.ToCompetitionListResp(items), logger)
}

// get returns a single competition with its media.
//
//	@Summary		Get a competition
//	@Description	Public fetch of a competition by ID, including associated media.
//	@Tags			Competitions
//	@Produce		json
//	@Param			id	path		string	true	"Competition UUID"
//	@Success		200	{object}	dto.CompetitionResp
//	@Failure		400	{object}	dto.ErrorResponse	"Bad Request"
//	@Failure		404	{object}	dto.ErrorResponse	"Not Found"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router			/apis/competition/v1/competitions/{id} [get]
func (h *competition) get(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "Get")
	ctx := r.Context()

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	c, err := h.uc.Get(ctx, id)
	if err != nil {
		logger.ErrorContext(ctx, "get failed", "error", err)
		dto.HandleError(err, w)

		return
	}

	writeJSON(ctx, w, http.StatusOK, dto.ToCompetitionResp(c), logger)
}

// create adds a competition (admin).
//
//	@Summary		Create a competition
//	@Description	Admin: create a new competition.
//	@Tags			Competitions (Admin)
//	@Accept			json
//	@Produce		json
//	@Param			competition	body		dto.CompetitionReq	true	"Competition"
//	@Success		201			{object}	dto.CompetitionResp
//	@Failure		400			{object}	dto.ErrorResponse	"Bad Request"
//	@Failure		409			{object}	dto.ErrorResponse	"Conflict"
//	@Failure		500			{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router			/apis/competition/v1/admin/competitions [post]
func (h *competition) create(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "Create")
	ctx := r.Context()

	req := new(dto.CompetitionReq)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	in, err := req.ToCreateInput()
	if err != nil {
		dto.HandleError(err, w)

		return
	}

	c, err := h.uc.Create(ctx, in)
	if err != nil {
		logger.ErrorContext(ctx, "create failed", "error", err)
		dto.HandleError(err, w)

		return
	}

	h.audit.Record(ctx, audit.Entry{
		Action: "competition.create", EntityType: "competition", EntityID: c.ID.String(), Reason: c.Title,
	})

	writeJSON(ctx, w, http.StatusCreated, dto.ToCompetitionResp(c), logger)
}

// update replaces a competition's editable fields (admin).
//
//	@Summary		Update a competition
//	@Description	Admin: update an existing competition.
//	@Tags			Competitions (Admin)
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string				true	"Competition UUID"
//	@Param			competition	body		dto.CompetitionReq	true	"Competition"
//	@Success		200			{object}	dto.CompetitionResp
//	@Failure		400			{object}	dto.ErrorResponse	"Bad Request"
//	@Failure		404			{object}	dto.ErrorResponse	"Not Found"
//	@Failure		500			{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router			/apis/competition/v1/admin/competitions/{id} [put]
func (h *competition) update(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "Update")
	ctx := r.Context()

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	req := new(dto.CompetitionReq)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	in, err := req.ToUpdateInput()
	if err != nil {
		dto.HandleError(err, w)

		return
	}

	c, err := h.uc.Update(ctx, id, in)
	if err != nil {
		logger.ErrorContext(ctx, "update failed", "error", err)
		dto.HandleError(err, w)

		return
	}

	h.audit.Record(ctx, audit.Entry{
		Action: "competition.update", EntityType: "competition", EntityID: c.ID.String(), Reason: c.Title,
	})

	writeJSON(ctx, w, http.StatusOK, dto.ToCompetitionResp(c), logger)
}

// delete removes a competition (admin).
//
//	@Summary		Delete a competition
//	@Description	Admin: delete a competition by ID.
//	@Tags			Competitions (Admin)
//	@Produce		json
//	@Param			id	path	string	true	"Competition UUID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	dto.ErrorResponse	"Bad Request"
//	@Failure		404	{object}	dto.ErrorResponse	"Not Found"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router			/apis/competition/v1/admin/competitions/{id} [delete]
func (h *competition) delete(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "Delete")
	ctx := r.Context()

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	if err := h.uc.Delete(ctx, id); err != nil {
		logger.ErrorContext(ctx, "delete failed", "error", err)
		dto.HandleError(err, w)

		return
	}

	h.audit.Record(ctx, audit.Entry{
		Action: "competition.delete", EntityType: "competition", EntityID: id.String(),
	})

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(ctx context.Context, w http.ResponseWriter, status int, body any, logger *slog.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		logger.ErrorContext(ctx, "failed to encode response", "error", err)
	}
}
