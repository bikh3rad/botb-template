package handler

import (
	"context"
	"log/slog"

	"application/internal/adminauth/biz"
	"application/internal/datasource"
	"application/internal/service"
	svchandler "application/internal/service/handler"
	"application/pkg/audit"

	"github.com/google/wire"
)

// ProviderSet wires the adminauth handlers, the audit recorder, the bootstrap
// step and this service's handler list.
var ProviderSet = wire.NewSet(
	NewAuth,
	NewAccounts,
	NewAuditRecorder,
	RunBootstrap,
	NewServiceList,
)

// NewAuditRecorder wires the shared audit recorder against this service's DB.
func NewAuditRecorder(logger *slog.Logger, db *datasource.PostgresDB) *audit.Recorder {
	return audit.NewRecorder(logger, db)
}

// BootstrapDone is a wire marker proving the first-superadmin bootstrap ran
// before the HTTP handlers were assembled.
type BootstrapDone struct{}

// RunBootstrap seeds the first superadmin from config/env (idempotent; see
// biz.EnsureBootstrap). Failing the seed fails service startup deliberately —
// a half-configured bootstrap should be loud, not silent.
func RunBootstrap(ctx context.Context, logger *slog.Logger, cfg *biz.Config, repo biz.Repository) (BootstrapDone, error) {
	if err := biz.EnsureBootstrap(ctx, logger, cfg, repo); err != nil {
		return BootstrapDone{}, err
	}

	return BootstrapDone{}, nil
}

// NewServiceList assembles the []service.Handler served by the adminauth
// binary: shared healthz plus the session and account-management handlers.
func NewServiceList(healthz *svchandler.HealthzHandler, a *auth, acc *accounts, _ BootstrapDone) []service.Handler {
	return []service.Handler{
		healthz,
		a,
		acc,
	}
}
