package handler

import (
	"application/internal/service"
	svchandler "application/internal/service/handler"

	"github.com/google/wire"
)

// ProviderSet wires the media handler and this service's handler list.
var ProviderSet = wire.NewSet(
	NewMedia,
	NewServiceList,
)

// NewServiceList assembles the []service.Handler served by the media binary:
// the shared healthz handler plus the media handler.
func NewServiceList(healthz *svchandler.HealthzHandler, m *media) []service.Handler {
	return []service.Handler{
		healthz,
		m,
	}
}
