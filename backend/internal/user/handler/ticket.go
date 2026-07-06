package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"application/internal/service"
	"application/internal/user/biz"
	"application/internal/user/dto"
	"application/pkg/middlewares"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type ticket struct {
	logger *slog.Logger
	tracer trace.Tracer
	mux    *http.ServeMux
	uc     biz.UsecaseTicket
	auth   *middlewares.JWTAuth
}

var _ service.Handler = (*ticket)(nil)

// NewTicket creates the ticket HTTP handler.
func NewTicket(logger *slog.Logger, mux *http.ServeMux, uc biz.UsecaseTicket, auth *middlewares.JWTAuth) *ticket {
	return &ticket{
		logger: logger.With("layer", "TicketHandler"),
		tracer: otel.Tracer("TicketHandler"),
		mux:    mux,
		uc:     uc,
		auth:   auth,
	}
}

// RegisterHandler mounts ticket purchase (public) and per-user listing (admin,
// JWT-guarded here as well as at the gateway — defense in depth).
func (h *ticket) RegisterHandler(_ context.Context) error {
	h.mux.HandleFunc("POST /apis/user/v1/tickets", h.purchase)
	h.mux.HandleFunc(
		"GET /apis/user/v1/admin/users/{id}/tickets",
		middlewares.MultipleMiddleware(h.listByUser, h.auth.RequireAdmin),
	)

	return nil
}

// purchase buys one or more tickets for a competition.
//
//	@Summary		Purchase tickets
//	@Description	Buy `quantity` tickets for a competition; charges the user and records the entries.
//	@Tags			Tickets
//	@Accept			json
//	@Produce		json
//	@Param			purchase	body		dto.PurchaseReq	true	"Purchase"
//	@Success		201			{object}	dto.PurchaseResp
//	@Failure		400			{object}	dto.ErrorResponse	"Bad Request"
//	@Failure		404			{object}	dto.ErrorResponse	"Not Found"
//	@Failure		500			{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router			/apis/user/v1/tickets [post]
func (h *ticket) purchase(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "Purchase")
	ctx := r.Context()

	req := new(dto.PurchaseReq)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	competitionID, err := uuid.Parse(req.CompetitionID)
	if err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, errors.New("invalid competition_id")), w)

		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, errors.New("invalid user_id")), w)

		return
	}

	result, err := h.uc.Purchase(ctx, biz.PurchaseInput{
		CompetitionID: competitionID,
		UserID:        userID,
		Quantity:      req.Quantity,
	})
	if err != nil {
		logger.ErrorContext(ctx, "purchase failed", "error", err)
		dto.HandleError(err, w)

		return
	}

	resp := dto.PurchaseResp{
		Tickets:        make([]dto.TicketResp, 0, len(result.Tickets)),
		User:           dto.ToUserResp(result.User),
		TotalCostPence: result.TotalCostPence,
		Count:          len(result.Tickets),
	}
	for _, t := range result.Tickets {
		resp.Tickets = append(resp.Tickets, dto.ToTicketResp(t))
	}

	writeJSON(ctx, w, http.StatusCreated, resp, logger)
}

// listByUser returns a user's tickets.
//
//	@Summary		List a user's tickets
//	@Description	Admin: list every ticket a user has bought.
//	@Tags			Tickets (Admin)
//	@Produce		json
//	@Param			id	path		string	true	"User UUID"
//	@Success		200	{object}	dto.TicketListResp
//	@Failure		400	{object}	dto.ErrorResponse	"Bad Request"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router			/apis/user/v1/admin/users/{id}/tickets [get]
func (h *ticket) listByUser(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "ListByUser")
	ctx := r.Context()

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	tickets, err := h.uc.ListByUser(ctx, id)
	if err != nil {
		logger.ErrorContext(ctx, "list failed", "error", err)
		dto.HandleError(err, w)

		return
	}

	writeJSON(ctx, w, http.StatusOK, dto.ToTicketListResp(tickets), logger)
}
