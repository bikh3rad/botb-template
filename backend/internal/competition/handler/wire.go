package handler

import (
	"log/slog"

	"application/internal/datasource"
	"application/internal/service"
	svchandler "application/internal/service/handler"
	"application/pkg/audit"

	"github.com/google/wire"
)

// ProviderSet wires the competition, category and content handlers, the audit
// recorder, and this service's handler list.
var ProviderSet = wire.NewSet(
	NewCompetition,
	NewCategoryHandler,
	NewContentHandler,
	NewAuditHandler,
	NewAuditRecorder,
	NewServiceList,
)

// NewAuditRecorder wires the shared audit recorder against this service's DB.
func NewAuditRecorder(logger *slog.Logger, db *datasource.PostgresDB) *audit.Recorder {
	return audit.NewRecorder(logger, db)
}

// NewServiceList assembles the []service.Handler served by the competition
// binary: shared healthz plus the competition, category and content handlers.
func NewServiceList(healthz *svchandler.HealthzHandler, c *competition, cat *category, sc *content, au *auditHandler) []service.Handler {
	return []service.Handler{
		healthz,
		c,
		cat,
		sc,
		au,
	}
}
