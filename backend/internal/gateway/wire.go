package gateway

import (
	"application/internal/service"
	svchandler "application/internal/service/handler"

	"github.com/google/wire"
)

// ProviderSet wires the gateway config, handler, and this binary's handler list.
var ProviderSet = wire.NewSet(
	NewGatewayConfig,
	NewGateway,
	NewServiceList,
)

// NewServiceList assembles the []service.Handler served by the gateway binary:
// the shared healthz handler plus the reverse-proxy dispatcher.
func NewServiceList(healthz *svchandler.HealthzHandler, g *gateway) []service.Handler {
	return []service.Handler{
		healthz,
		g,
	}
}
