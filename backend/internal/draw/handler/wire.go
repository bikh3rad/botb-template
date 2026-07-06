package handler

import (
	"log/slog"

	"application/internal/datasource"
	"application/internal/service"
	svchandler "application/internal/service/handler"
	"application/pkg/audit"

	"github.com/google/wire"
)

// ProviderSet wires the draw handler, the audit recorder, and this service's
// handler list.
var ProviderSet = wire.NewSet(
	NewDraw,
	NewAuditRecorder,
	NewServiceList,
)

// NewAuditRecorder wires the shared audit recorder against this service's DB.
func NewAuditRecorder(logger *slog.Logger, db *datasource.PostgresDB) *audit.Recorder {
	return audit.NewRecorder(logger, db)
}

// NewServiceList assembles the []service.Handler served by the draw binary: the
// shared healthz handler plus the draw handler.
func NewServiceList(healthz *svchandler.HealthzHandler, d *draw) []service.Handler {
	return []service.Handler{
		healthz,
		d,
	}
}
