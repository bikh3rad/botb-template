package biz

import (
	"github.com/google/wire"
)

// HealthzProviderSet wires only the healthz use case. Every service binary
// includes this so the /healthz endpoints work without pulling in the
// placeholder example. The bind lives here because *healthz is unexported.
var HealthzProviderSet = wire.NewSet(
	NewHealthz,
	wire.Bind(new(UsecaseHealthzer), new(*healthz)),
)

var BizProviderSet = wire.NewSet(
	HealthzProviderSet,

	NewPlaceholder,
	wire.Bind(new(UsecasePlaceholder), new(*placeholder)),
)
