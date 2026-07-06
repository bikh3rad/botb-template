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

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type category struct {
	logger *slog.Logger
	tracer trace.Tracer
	mux    *http.ServeMux
	uc     biz.UsecaseCategory
	auth   *middlewares.JWTAuth
	audit  *audit.Recorder
}

var _ service.Handler = (*category)(nil)

// NewCategoryHandler creates the category HTTP handler.
func NewCategoryHandler(
	logger *slog.Logger,
	mux *http.ServeMux,
	uc biz.UsecaseCategory,
	auth *middlewares.JWTAuth,
	recorder *audit.Recorder,
) *category {
	return &category{
		logger: logger.With("layer", "CategoryHandler"),
		tracer: otel.Tracer("CategoryHandler"),
		mux:    mux,
		uc:     uc,
		auth:   auth,
		audit:  recorder,
	}
}

// RegisterHandler mounts the public category list and the admin CRUD (role
// guarded here as well as at the gateway).
func (h *category) RegisterHandler(_ context.Context) error {
	admin := func(fn http.HandlerFunc) http.HandlerFunc {
		return middlewares.MultipleMiddleware(fn, h.auth.RequireAdmin)
	}

	h.mux.HandleFunc("GET /apis/competition/v1/categories", h.list)
	h.mux.HandleFunc("POST /apis/competition/v1/admin/categories", admin(h.create))
	h.mux.HandleFunc("PUT /apis/competition/v1/admin/categories/{id}", admin(h.update))
	h.mux.HandleFunc("DELETE /apis/competition/v1/admin/categories/{id}", admin(h.delete))

	return nil
}

// list returns all categories.
//
//	@Summary		List categories
//	@Description	Public list of competition categories.
//	@Tags			Categories
//	@Produce		json
//	@Success		200	{object}	dto.CategoryListResp
//	@Failure		500	{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router			/apis/competition/v1/categories [get]
func (h *category) list(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "List")
	ctx := r.Context()

	categories, err := h.uc.List(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "list failed", "error", err)
		dto.HandleError(err, w)

		return
	}

	writeJSON(ctx, w, http.StatusOK, dto.ToCategoryListResp(categories), logger)
}

// create adds a category.
//
//	@Summary		Create a category
//	@Description	Admin: create a category (name unique; slug derived when omitted).
//	@Tags			Categories (Admin)
//	@Accept			json
//	@Produce		json
//	@Param			category	body		dto.CategoryReq	true	"Category"
//	@Success		201			{object}	dto.CategoryResp
//	@Failure		400			{object}	dto.ErrorResponse	"Bad Request"
//	@Failure		409			{object}	dto.ErrorResponse	"Conflict"
//	@Router			/apis/competition/v1/admin/categories [post]
func (h *category) create(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "Create")
	ctx := r.Context()

	req := new(dto.CategoryReq)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	c, err := h.uc.Create(ctx, biz.CategoryInput{Name: req.Name, Slug: req.Slug})
	if err != nil {
		logger.ErrorContext(ctx, "create failed", "error", err)
		dto.HandleError(err, w)

		return
	}

	h.audit.Record(ctx, audit.Entry{
		Action: "category.create", EntityType: "category", EntityID: c.ID.String(), Reason: c.Name,
	})

	writeJSON(ctx, w, http.StatusCreated, dto.ToCategoryResp(c), logger)
}

// update renames a category.
//
//	@Summary		Update a category
//	@Description	Admin: rename a category / change its slug.
//	@Tags			Categories (Admin)
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string			true	"Category UUID"
//	@Param			category	body		dto.CategoryReq	true	"Category"
//	@Success		200			{object}	dto.CategoryResp
//	@Failure		400			{object}	dto.ErrorResponse	"Bad Request"
//	@Failure		404			{object}	dto.ErrorResponse	"Not Found"
//	@Failure		409			{object}	dto.ErrorResponse	"Conflict"
//	@Router			/apis/competition/v1/admin/categories/{id} [put]
func (h *category) update(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "Update")
	ctx := r.Context()

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	req := new(dto.CategoryReq)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	c, err := h.uc.Update(ctx, id, biz.CategoryInput{Name: req.Name, Slug: req.Slug})
	if err != nil {
		logger.ErrorContext(ctx, "update failed", "error", err)
		dto.HandleError(err, w)

		return
	}

	h.audit.Record(ctx, audit.Entry{
		Action: "category.update", EntityType: "category", EntityID: c.ID.String(), Reason: c.Name,
	})

	writeJSON(ctx, w, http.StatusOK, dto.ToCategoryResp(c), logger)
}

// delete removes a category (blocked or reassigned when still in use).
//
//	@Summary		Delete a category
//	@Description	Admin: delete a category. 409 while competitions still use it unless ?reassign_to=<uuid> moves them first (atomic).
//	@Tags			Categories (Admin)
//	@Param			id			path	string	true	"Category UUID"
//	@Param			reassign_to	query	string	false	"Category UUID to move competitions to"
//	@Success		204	"No Content"
//	@Failure		400	{object}	dto.ErrorResponse	"Bad Request"
//	@Failure		404	{object}	dto.ErrorResponse	"Not Found"
//	@Failure		409	{object}	dto.ErrorResponse	"Category in use"
//	@Router			/apis/competition/v1/admin/categories/{id} [delete]
func (h *category) delete(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "Delete")
	ctx := r.Context()

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	var reassignTo *uuid.UUID

	if raw := r.URL.Query().Get("reassign_to"); raw != "" {
		target, err := uuid.Parse(raw)
		if err != nil {
			dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

			return
		}

		reassignTo = &target
	}

	if err := h.uc.Delete(ctx, id, reassignTo); err != nil {
		logger.ErrorContext(ctx, "delete failed", "error", err)
		dto.HandleError(err, w)

		return
	}

	reason := "deleted"
	if reassignTo != nil {
		reason = "deleted; competitions reassigned to " + reassignTo.String()
	}

	h.audit.Record(ctx, audit.Entry{
		Action: "category.delete", EntityType: "category", EntityID: id.String(), Reason: reason,
	})

	w.WriteHeader(http.StatusNoContent)
}
