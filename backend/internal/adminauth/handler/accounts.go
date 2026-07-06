package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"application/internal/adminauth/biz"
	"application/internal/adminauth/dto"
	"application/internal/adminauth/entity"
	"application/internal/service"
	"application/pkg/audit"
	"application/pkg/middlewares"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type accounts struct {
	logger *slog.Logger
	tracer trace.Tracer
	mux    *http.ServeMux
	uc     biz.UsecaseAccounts
	jwt    *middlewares.JWTAuth
	audit  *audit.Recorder
}

var _ service.Handler = (*accounts)(nil)

// NewAccounts creates the superadmin account-management HTTP handler.
func NewAccounts(
	logger *slog.Logger,
	mux *http.ServeMux,
	uc biz.UsecaseAccounts,
	jwt *middlewares.JWTAuth,
	recorder *audit.Recorder,
) *accounts {
	return &accounts{
		logger: logger.With("layer", "AdminAccountsHandler"),
		tracer: otel.Tracer("AdminAccountsHandler"),
		mux:    mux,
		uc:     uc,
		jwt:    jwt,
		audit:  recorder,
	}
}

// RegisterHandler mounts account management. Superadmin-only, guarded here as
// well as at the gateway (defense in depth). Accounts are disabled, never
// hard-deleted, so there is no DELETE route.
func (h *accounts) RegisterHandler(_ context.Context) error {
	super := func(fn http.HandlerFunc) http.HandlerFunc {
		return middlewares.MultipleMiddleware(fn, h.jwt.RequireSuperadmin)
	}

	h.mux.HandleFunc("GET /apis/adminauth/v1/admin/accounts", super(h.list))
	h.mux.HandleFunc("POST /apis/adminauth/v1/admin/accounts", super(h.create))
	h.mux.HandleFunc("PUT /apis/adminauth/v1/admin/accounts/{id}", super(h.update))

	return nil
}

// list returns all admin accounts.
//
//	@Summary		List admin accounts
//	@Description	Superadmin: list every admin account (no password material).
//	@Tags			AdminAuth (Superadmin)
//	@Produce		json
//	@Success		200	{object}	dto.AccountListResp
//	@Failure		401	{object}	dto.ErrorResponse	"Unauthorized"
//	@Failure		403	{object}	dto.ErrorResponse	"Forbidden (not superadmin)"
//	@Router			/apis/adminauth/v1/admin/accounts [get]
func (h *accounts) list(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "List")
	ctx := r.Context()

	all, err := h.uc.List(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "list failed", "error", err)
		dto.HandleError(err, w)

		return
	}

	writeJSON(ctx, w, http.StatusOK, dto.ToAccountListResp(all), logger)
}

// create adds an admin account.
//
//	@Summary		Create an admin account
//	@Description	Superadmin: create an admin or superadmin account.
//	@Tags			AdminAuth (Superadmin)
//	@Accept			json
//	@Produce		json
//	@Param			account	body		dto.AccountCreateReq	true	"Account"
//	@Success		201		{object}	dto.AdminResp
//	@Failure		400		{object}	dto.ErrorResponse	"Bad Request"
//	@Failure		409		{object}	dto.ErrorResponse	"Email already in use"
//	@Router			/apis/adminauth/v1/admin/accounts [post]
func (h *accounts) create(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "Create")
	ctx := r.Context()

	req := new(dto.AccountCreateReq)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	account, err := h.uc.Create(ctx, biz.CreateAccountInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     entity.Role(req.Role),
	})
	if err != nil {
		logger.ErrorContext(ctx, "create failed", "error", err)
		dto.HandleError(err, w)

		return
	}

	h.audit.Record(ctx, audit.Entry{
		Action:     "adminauth.account.create",
		EntityType: "admin_account",
		EntityID:   account.ID.String(),
		Reason:     "created " + string(account.Role) + " " + account.Email,
	})

	writeJSON(ctx, w, http.StatusCreated, dto.ToAdminResp(account), logger)
}

// update edits or disables an admin account.
//
//	@Summary		Update an admin account
//	@Description	Superadmin: partial edit (name/email/password/role/is_active). Disabling or demoting the last active superadmin is refused.
//	@Tags			AdminAuth (Superadmin)
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"Account UUID"
//	@Param			account	body		dto.AccountUpdateReq	true	"Partial edit"
//	@Success		200		{object}	dto.AdminResp
//	@Failure		400		{object}	dto.ErrorResponse	"Bad Request"
//	@Failure		404		{object}	dto.ErrorResponse	"Not Found"
//	@Failure		409		{object}	dto.ErrorResponse	"Conflict (duplicate email / last superadmin)"
//	@Router			/apis/adminauth/v1/admin/accounts/{id} [put]
func (h *accounts) update(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "Update")
	ctx := r.Context()

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	req := new(dto.AccountUpdateReq)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	account, err := h.uc.Update(ctx, id, req.ToUpdateInput())
	if err != nil {
		logger.ErrorContext(ctx, "update failed", "error", err)
		dto.HandleError(err, w)

		return
	}

	h.audit.Record(ctx, audit.Entry{
		Action:     "adminauth.account.update",
		EntityType: "admin_account",
		EntityID:   account.ID.String(),
	})

	writeJSON(ctx, w, http.StatusOK, dto.ToAdminResp(account), logger)
}
