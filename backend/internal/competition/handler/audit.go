package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"application/internal/competition/biz"
	"application/internal/competition/dto"
	"application/internal/service"
	"application/pkg/middlewares"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type auditHandler struct {
	logger *slog.Logger
	tracer trace.Tracer
	mux    *http.ServeMux
	uc     biz.UsecaseAudit
	auth   *middlewares.JWTAuth
}

var _ service.Handler = (*auditHandler)(nil)

// NewAuditHandler creates the admin audit-log read handler.
func NewAuditHandler(logger *slog.Logger, mux *http.ServeMux, uc biz.UsecaseAudit, auth *middlewares.JWTAuth) *auditHandler {
	return &auditHandler{
		logger: logger.With("layer", "AuditHandler"),
		tracer: otel.Tracer("AuditHandler"),
		mux:    mux,
		uc:     uc,
		auth:   auth,
	}
}

// RegisterHandler mounts the admin-only recent-audit read.
func (h *auditHandler) RegisterHandler(_ context.Context) error {
	h.mux.HandleFunc("GET /apis/competition/v1/admin/audit",
		middlewares.MultipleMiddleware(h.recent, h.auth.RequireAdmin))

	return nil
}

// recent returns the most recent admin actions across all services.
//
//	@Summary		Recent admin audit log
//	@Description	Admin: recent admin_audit_log entries (all services write here), newest first.
//	@Tags			Audit (Admin)
//	@Produce		json
//	@Param			limit	query		int	false	"Max rows (default 20, max 100)"
//	@Success		200		{object}	dto.AuditListResp
//	@Failure		500		{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router			/apis/competition/v1/admin/audit [get]
func (h *auditHandler) recent(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "Recent")
	ctx := r.Context()

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	entries, err := h.uc.Recent(ctx, limit)
	if err != nil {
		logger.ErrorContext(ctx, "recent audit failed", "error", err)
		dto.HandleError(err, w)

		return
	}

	writeJSON(ctx, w, http.StatusOK, dto.ToAuditListResp(entries), logger)
}
