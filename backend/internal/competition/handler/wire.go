package handler

import (
	"application/internal/service"
	svchandler "application/internal/service/handler"

	"github.com/google/wire"
)

// ProviderSet wires the competition handler and this service's handler list.
var ProviderSet = wire.NewSet(
	NewCompetition,
	NewServiceList,
)

// NewServiceList assembles the []service.Handler served by the competition
// binary: the shared healthz handler plus the competition handler.
func NewServiceList(healthz *svchandler.HealthzHandler, c *competition) []service.Handler {
	return []service.Handler{
		healthz,
		c,
	}
}
