package handler

import (
	"application/internal/service"
	svchandler "application/internal/service/handler"

	"github.com/google/wire"
)

// ProviderSet wires the draw handler and this service's handler list.
var ProviderSet = wire.NewSet(
	NewDraw,
	NewServiceList,
)

// NewServiceList assembles the []service.Handler served by the draw binary: the
// shared healthz handler plus the draw handler.
func NewServiceList(healthz *svchandler.HealthzHandler, d *draw) []service.Handler {
	return []service.Handler{
		healthz,
		d,
	}
}
