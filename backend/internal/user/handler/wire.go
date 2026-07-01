package handler

import (
	"application/internal/service"
	svchandler "application/internal/service/handler"

	"github.com/google/wire"
)

// ProviderSet wires the user + ticket handlers and this service's handler list.
var ProviderSet = wire.NewSet(
	NewUser,
	NewTicket,
	NewServiceList,
)

// NewServiceList assembles the []service.Handler served by the user binary: the
// shared healthz handler plus the user and ticket handlers.
func NewServiceList(healthz *svchandler.HealthzHandler, u *user, t *ticket) []service.Handler {
	return []service.Handler{
		healthz,
		u,
		t,
	}
}
