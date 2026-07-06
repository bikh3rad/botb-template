package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"application/internal/competition/biz"
	"application/internal/competition/dto"
	"application/internal/service"
	"application/pkg/audit"
	"application/pkg/middlewares"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type content struct {
	logger *slog.Logger
	tracer trace.Tracer
	mux    *http.ServeMux
	uc     biz.UsecaseContent
	auth   *middlewares.JWTAuth
	audit  *audit.Recorder
}

var _ service.Handler = (*content)(nil)

// NewContentHandler creates the site-content HTTP handler.
func NewContentHandler(
	logger *slog.Logger,
	mux *http.ServeMux,
	uc biz.UsecaseContent,
	auth *middlewares.JWTAuth,
	recorder *audit.Recorder,
) *content {
	return &content{
		logger: logger.With("layer", "SiteContentHandler"),
		tracer: otel.Tracer("SiteContentHandler"),
		mux:    mux,
		uc:     uc,
		auth:   auth,
		audit:  recorder,
	}
}

// RegisterHandler mounts the public read and the admin write.
func (h *content) RegisterHandler(_ context.Context) error {
	h.mux.HandleFunc("GET /apis/competition/v1/content", h.getAll)
	h.mux.HandleFunc("PUT /apis/competition/v1/admin/content/{key}",
		middlewares.MultipleMiddleware(h.upsert, h.auth.RequireAdmin))

	return nil
}

// getAll returns every site-copy value.
//
//	@Summary		Site content
//	@Description	Public read of all editable site copy (key/value).
//	@Tags			SiteContent
//	@Produce		json
//	@Success		200	{object}	dto.ContentResp
//	@Failure		500	{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router			/apis/competition/v1/content [get]
func (h *content) getAll(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "GetAll")
	ctx := r.Context()

	rows, err := h.uc.GetAll(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "get all failed", "error", err)
		dto.HandleError(err, w)

		return
	}

	writeJSON(ctx, w, http.StatusOK, dto.ToContentResp(rows), logger)
}

// upsert writes one site-copy value.
//
//	@Summary		Update site content
//	@Description	Admin: upsert one site-copy value by key.
//	@Tags			SiteContent (Admin)
//	@Accept			json
//	@Produce		json
//	@Param			key		path		string					true	"Content key"
//	@Param			content	body		dto.ContentUpsertReq	true	"Value"
//	@Success		200		{object}	dto.ContentItemResp
//	@Failure		400		{object}	dto.ErrorResponse	"Bad Request"
//	@Router			/apis/competition/v1/admin/content/{key} [put]
func (h *content) upsert(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "Upsert")
	ctx := r.Context()

	req := new(dto.ContentUpsertReq)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	row, err := h.uc.Upsert(ctx, r.PathValue("key"), req.Value)
	if err != nil {
		logger.ErrorContext(ctx, "upsert failed", "error", err)
		dto.HandleError(err, w)

		return
	}

	h.audit.Record(ctx, audit.Entry{
		Action: "content.update", EntityType: "site_content", EntityID: row.Key,
	})

	writeJSON(ctx, w, http.StatusOK, dto.ContentItemResp{
		Key: row.Key, Value: row.Value, UpdatedAt: row.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
	}, logger)
}
