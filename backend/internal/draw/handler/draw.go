package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"application/internal/draw/biz"
	"application/internal/draw/dto"
	"application/internal/service"
	"application/pkg/audit"
	"application/pkg/middlewares"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type draw struct {
	logger *slog.Logger
	tracer trace.Tracer
	mux    *http.ServeMux
	uc     biz.UsecaseDraw
	auth   *middlewares.JWTAuth
	audit  *audit.Recorder
}

var _ service.Handler = (*draw)(nil)

// NewDraw creates the draw HTTP handler. Draws are the most sensitive admin
// surface, so EVERY admin mutation here writes an admin_audit_log entry.
func NewDraw(
	logger *slog.Logger,
	mux *http.ServeMux,
	uc biz.UsecaseDraw,
	auth *middlewares.JWTAuth,
	recorder *audit.Recorder,
) *draw {
	return &draw{
		logger: logger.With("layer", "DrawHandler"),
		tracer: otel.Tracer("DrawHandler"),
		mux:    mux,
		uc:     uc,
		auth:   auth,
		audit:  recorder,
	}
}

// RegisterHandler mounts admin management + the public draw-result read. The
// admin group is JWT-guarded here too (defense in depth) — the gateway also
// guards it, but a service reached directly on its own port must not be
// unprotected.
func (h *draw) RegisterHandler(_ context.Context) error {
	admin := func(fn http.HandlerFunc) http.HandlerFunc {
		return middlewares.MultipleMiddleware(fn, h.auth.RequireAdmin)
	}

	// Admin.
	h.mux.HandleFunc("GET /apis/draw/v1/admin/draws", admin(h.list))
	h.mux.HandleFunc("GET /apis/draw/v1/admin/draws/{id}", admin(h.get))
	h.mux.HandleFunc("POST /apis/draw/v1/admin/draws", admin(h.create))
	h.mux.HandleFunc("POST /apis/draw/v1/admin/draws/{id}/run", admin(h.run))
	h.mux.HandleFunc("PUT /apis/draw/v1/admin/draws/{id}", admin(h.update))
	h.mux.HandleFunc("POST /apis/draw/v1/admin/draws/{id}/void", admin(h.void))
	h.mux.HandleFunc("POST /apis/draw/v1/admin/draws/{id}/reassign", admin(h.reassign))
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

	h.audit.Record(ctx, audit.Entry{
		Action: "draw.create", EntityType: "draw", EntityID: d.ID.String(), Reason: d.Prize,
	})

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

	winner := ""
	if d.WinnerUserID != nil {
		winner = d.WinnerUserID.String()
	}

	h.audit.Record(ctx, audit.Entry{
		Action: "draw.run", EntityType: "draw", EntityID: d.ID.String(),
		Reason: "winner user " + winner,
	})

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

// update edits a draw's prize text (admin, audited).
//
//	@Summary		Update a draw
//	@Description	Admin: edit the prize text. Winner fields are not editable here — void+re-run or use the audited reassign endpoint.
//	@Tags			Draws (Admin)
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"Draw UUID"
//	@Param			draw	body		dto.UpdateDrawReq	true	"Prize"
//	@Success		200		{object}	dto.DrawResp
//	@Failure		400		{object}	dto.ErrorResponse	"Bad Request"
//	@Failure		404		{object}	dto.ErrorResponse	"Not Found"
//	@Failure		409		{object}	dto.ErrorResponse	"Void draws are frozen"
//	@Router			/apis/draw/v1/admin/draws/{id} [put]
func (h *draw) update(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "Update")
	ctx := r.Context()

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	req := new(dto.UpdateDrawReq)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	d, err := h.uc.UpdatePrize(ctx, id, req.Prize)
	if err != nil {
		logger.ErrorContext(ctx, "update failed", "error", err)
		dto.HandleError(err, w)

		return
	}

	h.audit.Record(ctx, audit.Entry{
		Action: "draw.update", EntityType: "draw", EntityID: d.ID.String(), Reason: "prize -> " + d.Prize,
	})

	writeJSON(ctx, w, http.StatusOK, dto.ToDrawResp(d), logger)
}

// void voids a draw with a required reason (admin, audited). The safe path to
// CHANGE a winner is: void this draw, create a new one, run it.
//
//	@Summary		Void a draw
//	@Description	Admin: mark a pending/drawn draw void with a required reason. Voided draws disappear from the public winners feed.
//	@Tags			Draws (Admin)
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string			true	"Draw UUID"
//	@Param			void	body		dto.VoidDrawReq	true	"Reason (required)"
//	@Success		200		{object}	dto.DrawResp
//	@Failure		400		{object}	dto.ErrorResponse	"Reason missing"
//	@Failure		404		{object}	dto.ErrorResponse	"Not Found"
//	@Failure		409		{object}	dto.ErrorResponse	"Already void"
//	@Router			/apis/draw/v1/admin/draws/{id}/void [post]
func (h *draw) void(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "Void")
	ctx := r.Context()

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	req := new(dto.VoidDrawReq)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		dto.HandleError(errors.Join(biz.ErrReasonRequired, err), w)

		return
	}

	d, err := h.uc.Void(ctx, id, req.Reason)
	if err != nil {
		logger.ErrorContext(ctx, "void failed", "error", err)
		dto.HandleError(err, w)

		return
	}

	h.audit.Record(ctx, audit.Entry{
		Action: "draw.void", EntityType: "draw", EntityID: d.ID.String(), Reason: req.Reason,
	})

	writeJSON(ctx, w, http.StatusOK, dto.ToDrawResp(d), logger)
}

// reassign directly changes a drawn draw's winner (admin, audited). It exists
// precisely so this mutation is explicit, validated and attributable instead
// of a hand-edited row.
//
//	@Summary		Reassign a draw's winner
//	@Description	Admin: move a DRAWN draw's winner to another ticket of the same competition. Requires a reason; writes an audit entry.
//	@Tags			Draws (Admin)
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string				true	"Draw UUID"
//	@Param			reassign	body		dto.ReassignDrawReq	true	"New winner ticket + reason"
//	@Success		200			{object}	dto.DrawResp
//	@Failure		400			{object}	dto.ErrorResponse	"Bad Request / reason missing"
//	@Failure		404			{object}	dto.ErrorResponse	"Not Found"
//	@Failure		409			{object}	dto.ErrorResponse	"Draw not in drawn state"
//	@Failure		422			{object}	dto.ErrorResponse	"Ticket not in this competition"
//	@Router			/apis/draw/v1/admin/draws/{id}/reassign [post]
func (h *draw) reassign(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "Reassign")
	ctx := r.Context()

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	req := new(dto.ReassignDrawReq)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	ticketID, err := uuid.Parse(req.WinnerTicketID)
	if err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	d, err := h.uc.Reassign(ctx, id, ticketID, req.Reason)
	if err != nil {
		logger.ErrorContext(ctx, "reassign failed", "error", err)
		dto.HandleError(err, w)

		return
	}

	winner := ""
	if d.WinnerUserID != nil {
		winner = d.WinnerUserID.String()
	}

	h.audit.Record(ctx, audit.Entry{
		Action: "draw.reassign", EntityType: "draw", EntityID: d.ID.String(),
		Reason: req.Reason + " (new winner ticket " + ticketID.String() + ", user " + winner + ")",
	})

	writeJSON(ctx, w, http.StatusOK, dto.ToDrawResp(d), logger)
}

func writeJSON(ctx context.Context, w http.ResponseWriter, status int, body any, logger *slog.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		logger.ErrorContext(ctx, "failed to encode response", "error", err)
	}
}
