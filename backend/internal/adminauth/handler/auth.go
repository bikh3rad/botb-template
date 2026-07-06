package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"strings"

	"application/internal/adminauth/biz"
	"application/internal/adminauth/dto"
	"application/internal/service"
	"application/pkg/middlewares"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type auth struct {
	logger *slog.Logger
	tracer trace.Tracer
	mux    *http.ServeMux
	uc     biz.UsecaseAuth
	jwt    *middlewares.JWTAuth
}

var _ service.Handler = (*auth)(nil)

// NewAuth creates the adminauth session HTTP handler.
func NewAuth(logger *slog.Logger, mux *http.ServeMux, uc biz.UsecaseAuth, jwt *middlewares.JWTAuth) *auth {
	return &auth{
		logger: logger.With("layer", "AdminAuthHandler"),
		tracer: otel.Tracer("AdminAuthHandler"),
		mux:    mux,
		uc:     uc,
		jwt:    jwt,
	}
}

// RegisterHandler mounts the session endpoints. login/refresh/logout are
// public (you cannot hold a token before logging in); /me requires a valid
// admin token (defense in depth — the gateway does not guard non-/admin/
// paths of this service).
func (h *auth) RegisterHandler(_ context.Context) error {
	h.mux.HandleFunc("POST /apis/adminauth/v1/login", h.login)
	h.mux.HandleFunc("POST /apis/adminauth/v1/refresh", h.refresh)
	h.mux.HandleFunc("POST /apis/adminauth/v1/logout", h.logout)
	h.mux.HandleFunc("GET /apis/adminauth/v1/me",
		middlewares.MultipleMiddleware(h.me, h.jwt.RequireAdmin))

	return nil
}

// login verifies credentials and issues an access+refresh token pair.
//
//	@Summary		Admin login
//	@Description	Verify admin credentials; returns a short-lived access JWT and a rotating refresh token.
//	@Tags			AdminAuth
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body		dto.LoginReq	true	"Credentials"
//	@Success		200			{object}	dto.TokenResp
//	@Failure		401			{object}	dto.ErrorResponse	"Invalid credentials"
//	@Failure		429			{object}	dto.ErrorResponse	"Rate limited"
//	@Router			/apis/adminauth/v1/login [post]
func (h *auth) login(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "Login")
	ctx := r.Context()

	req := new(dto.LoginReq)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		dto.HandleError(errors.Join(biz.ErrResourceInvalid, err), w)

		return
	}

	result, err := h.uc.Login(ctx, req.Email, req.Password, clientIP(r))
	if err != nil {
		// Generic log line: never log the password, and failures are expected.
		logger.WarnContext(ctx, "login rejected")
		dto.HandleError(err, w)

		return
	}

	writeJSON(ctx, w, http.StatusOK, dto.ToTokenResp(result), logger)
}

// refresh rotates a refresh token into a new access+refresh pair.
//
//	@Summary		Refresh admin session
//	@Description	Exchange a valid refresh token for a new pair (rotation; reuse revokes all sessions).
//	@Tags			AdminAuth
//	@Accept			json
//	@Produce		json
//	@Param			token	body		dto.RefreshReq	true	"Refresh token"
//	@Success		200		{object}	dto.TokenResp
//	@Failure		401		{object}	dto.ErrorResponse	"Invalid refresh token"
//	@Router			/apis/adminauth/v1/refresh [post]
func (h *auth) refresh(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "Refresh")
	ctx := r.Context()

	req := new(dto.RefreshReq)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil || req.RefreshToken == "" {
		dto.HandleError(biz.ErrInvalidRefresh, w)

		return
	}

	result, err := h.uc.Refresh(ctx, req.RefreshToken)
	if err != nil {
		dto.HandleError(err, w)

		return
	}

	writeJSON(ctx, w, http.StatusOK, dto.ToTokenResp(result), logger)
}

// logout revokes a refresh token (idempotent).
//
//	@Summary		Admin logout
//	@Description	Revoke a refresh token.
//	@Tags			AdminAuth
//	@Accept			json
//	@Success		204	"No Content"
//	@Router			/apis/adminauth/v1/logout [post]
func (h *auth) logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := new(dto.RefreshReq)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil || req.RefreshToken == "" {
		w.WriteHeader(http.StatusNoContent)

		return
	}

	if err := h.uc.Logout(ctx, req.RefreshToken); err != nil {
		dto.HandleError(err, w)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// me returns the profile behind the presented access token.
//
//	@Summary		Current admin
//	@Description	Profile + role of the authenticated admin.
//	@Tags			AdminAuth
//	@Produce		json
//	@Success		200	{object}	dto.AdminResp
//	@Failure		401	{object}	dto.ErrorResponse	"Unauthorized"
//	@Router			/apis/adminauth/v1/me [get]
func (h *auth) me(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "Me")
	ctx := r.Context()

	claims, ok := middlewares.ClaimsFromContext(ctx)
	if !ok {
		dto.HandleError(biz.ErrInvalidCredentials, w)

		return
	}

	account, err := h.uc.Me(ctx, claims.Subject)
	if err != nil {
		dto.HandleError(err, w)

		return
	}

	writeJSON(ctx, w, http.StatusOK, dto.ToAdminResp(account), logger)
}

// clientIP prefers the first X-Forwarded-For hop (set by the gateway's reverse
// proxy) and falls back to the socket peer.
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if first, _, ok := strings.Cut(xff, ","); ok {
			return strings.TrimSpace(first)
		}

		return strings.TrimSpace(xff)
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return host
}

func writeJSON(ctx context.Context, w http.ResponseWriter, status int, body any, logger *slog.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		logger.ErrorContext(ctx, "failed to encode response", "error", err)
	}
}
