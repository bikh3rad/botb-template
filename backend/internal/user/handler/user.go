package handler

import (
	"application/internal/service"
	"application/internal/user/biz"
	"application/internal/user/dto"
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

type user struct {
	logger *slog.Logger
	tracer trace.Tracer
	mux    *http.ServeMux
	uc     biz.UsecaseUser
}

var _ service.Handler = (*user)(nil)

// NewUser creates the user HTTP handler.
func NewUser(logger *slog.Logger, mux *http.ServeMux, uc biz.UsecaseUser) *user {
	return &user{
		logger: logger.With("layer", "UserHandler"),
		tracer: otel.Tracer("UserHandler"),
		mux:    mux,
		uc:     uc,
	}
}

// RegisterHandler mounts registration (public) and user management (admin).
func (h *user) RegisterHandler(_ context.Context) error {
	h.mux.HandleFunc("POST /apis/user/v1/users", h.register)
	h.mux.HandleFunc("GET /apis/user/v1/admin/users", h.list)
	h.mux.HandleFunc("GET /apis/user/v1/admin/users/{id}", h.get)

	return nil
}

// register creates a new user account.
//
//	@Summary		Register a user
//	@Description	Public self-registration with a name and email.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dto.RegisterReq	true	"Registration"
//	@Success		201		{object}	dto.UserResp
//	@Failure		400		{object}	dto.ErrorResponse	"Bad Request"
//	@Failure		409		{object}	dto.ErrorResponse	"Conflict"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router			/apis/user/v1/users [post]
func (h *user) register(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "Register")
	ctx := r.Context()

	req := new(dto.RegisterReq)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	u, err := h.uc.Register(ctx, req.Name, req.Email)
	if err != nil {
		logger.ErrorContext(ctx, "register failed", "error", err)
		dto.HandleError(err, w)

		return
	}

	writeJSON(ctx, w, http.StatusCreated, dto.ToUserResp(u), logger)
}

// list returns a paginated, searchable list of users.
//
//	@Summary		List users
//	@Description	Admin: paginated, searchable user list (name/email).
//	@Tags			Users (Admin)
//	@Produce		json
//	@Param			q		query		string	false	"Search name or email"
//	@Param			limit	query		int		false	"Page size (default 20, max 100)"
//	@Param			offset	query		int		false	"Offset"
//	@Success		200		{object}	dto.UserListResp
//	@Failure		500		{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router			/apis/user/v1/admin/users [get]
func (h *user) list(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "List")
	ctx := r.Context()

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	filter := biz.UserListFilter{
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

	resp := dto.ToUserListResp(page.Users, page.Total, filter.Limit, filter.Offset)
	writeJSON(ctx, w, http.StatusOK, resp, logger)
}

// get returns a single user.
//
//	@Summary		Get a user
//	@Description	Admin: fetch a user by ID.
//	@Tags			Users (Admin)
//	@Produce		json
//	@Param			id	path		string	true	"User UUID"
//	@Success		200	{object}	dto.UserResp
//	@Failure		400	{object}	dto.ErrorResponse	"Bad Request"
//	@Failure		404	{object}	dto.ErrorResponse	"Not Found"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router			/apis/user/v1/admin/users/{id} [get]
func (h *user) get(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "Get")
	ctx := r.Context()

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	u, err := h.uc.Get(ctx, id)
	if err != nil {
		logger.ErrorContext(ctx, "get failed", "error", err)
		dto.HandleError(err, w)

		return
	}

	writeJSON(ctx, w, http.StatusOK, dto.ToUserResp(u), logger)
}

func writeJSON(ctx context.Context, w http.ResponseWriter, status int, body any, logger *slog.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		logger.ErrorContext(ctx, "failed to encode response", "error", err)
	}
}
