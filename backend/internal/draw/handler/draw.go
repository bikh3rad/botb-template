package handler

import (
	"application/internal/draw/biz"
	"application/internal/draw/dto"
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

type draw struct {
	logger *slog.Logger
	tracer trace.Tracer
	mux    *http.ServeMux
	uc     biz.UsecaseDraw
}

var _ service.Handler = (*draw)(nil)

// NewDraw creates the draw HTTP handler.
func NewDraw(logger *slog.Logger, mux *http.ServeMux, uc biz.UsecaseDraw) *draw {
	return &draw{
		logger: logger.With("layer", "DrawHandler"),
		tracer: otel.Tracer("DrawHandler"),
		mux:    mux,
		uc:     uc,
	}
}

// RegisterHandler mounts admin management + the public draw-result read. Admin
// routes sit under /admin/ so the gateway can guard that group with JWT.
func (h *draw) RegisterHandler(_ context.Context) error {
	// Admin.
	h.mux.HandleFunc("GET /apis/draw/v1/admin/draws", h.list)
	h.mux.HandleFunc("GET /apis/draw/v1/admin/draws/{id}", h.get)
	h.mux.HandleFunc("POST /apis/draw/v1/admin/draws", h.create)
	h.mux.HandleFunc("POST /apis/draw/v1/admin/draws/{id}/run", h.run)
	// Public.
	h.mux.HandleFunc("GET /apis/draw/v1/draws/{id}", h.getPublic)

	return nil
}

// list returns a paginated, searchable list of draws.
//
//	@Summary		List draws
//	@Description	Admin: paginated, searchable draw list (by prize).
//	@Tags			Draws (Admin)
//	@Produce		json
//	@Param			q		query		string	false	"Search prize"
//	@Param			limit	query		int		false	"Page size (default 20, max 100)"
//	@Param			offset	query		int		false	"Offset"
//	@Success		200		{object}	dto.DrawListResp
//	@Failure		500		{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router			/apis/draw/v1/admin/draws [get]
func (h *draw) list(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "List")
	ctx := r.Context()

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	filter := biz.ListFilter{
		Query:  r.URL.Query().Get("q"),
		Limit:  limit,
		Offset: offset,
	}

	page, err := h.uc.List(ctx, filter)
	if err != nil {
		logger.ErrorContext(ctx, "list failed", "error", err)
		dto.HandleError(err, w)

		return
	}

	resp := dto.ToDrawListResp(page.Draws, page.Total, filter.Limit, filter.Offset)
	writeJSON(ctx, w, http.StatusOK, resp, logger)
}

// get returns a single draw (admin — includes pending internals).
//
//	@Summary		Get a draw
//	@Description	Admin: fetch a draw by ID.
//	@Tags			Draws (Admin)
//	@Produce		json
//	@Param			id	path		string	true	"Draw UUID"
//	@Success		200	{object}	dto.DrawResp
//	@Failure		400	{object}	dto.ErrorResponse	"Bad Request"
//	@Failure		404	{object}	dto.ErrorResponse	"Not Found"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router			/apis/draw/v1/admin/draws/{id} [get]
func (h *draw) get(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "Get")
	ctx := r.Context()

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	d, err := h.uc.Get(ctx, id)
	if err != nil {
		logger.ErrorContext(ctx, "get failed", "error", err)
		dto.HandleError(err, w)

		return
	}

	writeJSON(ctx, w, http.StatusOK, dto.ToDrawResp(d), logger)
}

// create adds a pending draw for a competition (admin).
//
//	@Summary		Create a draw
//	@Description	Admin: create a pending draw for a competition.
//	@Tags			Draws (Admin)
//	@Accept			json
//	@Produce		json
//	@Param			draw	body		dto.CreateDrawReq	true	"Draw"
//	@Success		201		{object}	dto.DrawResp
//	@Failure		400		{object}	dto.ErrorResponse	"Bad Request"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router			/apis/draw/v1/admin/draws [post]
func (h *draw) create(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "Create")
	ctx := r.Context()

	req := new(dto.CreateDrawReq)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	competitionID, err := uuid.Parse(req.CompetitionID)
	if err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, errors.New("invalid competition_id")), w)

		return
	}

	d, err := h.uc.Create(ctx, biz.CreateInput{CompetitionID: competitionID, Prize: req.Prize})
	if err != nil {
		logger.ErrorContext(ctx, "create failed", "error", err)
		dto.HandleError(err, w)

		return
	}

	writeJSON(ctx, w, http.StatusCreated, dto.ToDrawResp(d), logger)
}

// run executes the draw, selecting a winner.
//
//	@Summary		Run a draw
//	@Description	Admin: run a pending draw — pick a winning ticket and record the winner. Rejects a re-run.
//	@Tags			Draws (Admin)
//	@Produce		json
//	@Param			id	path		string	true	"Draw UUID"
//	@Success		200	{object}	dto.DrawResp
//	@Failure		400	{object}	dto.ErrorResponse	"Bad Request"
//	@Failure		404	{object}	dto.ErrorResponse	"Not Found"
//	@Failure		409	{object}	dto.ErrorResponse	"Already Drawn"
//	@Failure		422	{object}	dto.ErrorResponse	"No Tickets"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router			/apis/draw/v1/admin/draws/{id}/run [post]
func (h *draw) run(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "Run")
	ctx := r.Context()

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	d, err := h.uc.Run(ctx, id)
	if err != nil {
		logger.ErrorContext(ctx, "run failed", "error", err)
		dto.HandleError(err, w)

		return
	}

	writeJSON(ctx, w, http.StatusOK, dto.ToDrawResp(d), logger)
}

// getPublic returns a completed draw result, hiding pending draws.
//
//	@Summary		Get a draw result (public)
//	@Description	Public: read a completed (drawn/void) draw result. Pending draws return 404.
//	@Tags			Draws
//	@Produce		json
//	@Param			id	path		string	true	"Draw UUID"
//	@Success		200	{object}	dto.DrawResp
//	@Failure		400	{object}	dto.ErrorResponse	"Bad Request"
//	@Failure		404	{object}	dto.ErrorResponse	"Not Found"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router			/apis/draw/v1/draws/{id} [get]
func (h *draw) getPublic(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "GetPublic")
	ctx := r.Context()

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	d, err := h.uc.GetPublic(ctx, id)
	if err != nil {
		logger.ErrorContext(ctx, "public get failed", "error", err)
		dto.HandleError(err, w)

		return
	}

	writeJSON(ctx, w, http.StatusOK, dto.ToDrawResp(d), logger)
}

func writeJSON(ctx context.Context, w http.ResponseWriter, status int, body any, logger *slog.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		logger.ErrorContext(ctx, "failed to encode response", "error", err)
	}
}
