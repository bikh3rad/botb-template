package biz

import "github.com/google/wire"

// ProviderSet wires the draw use case.
var ProviderSet = wire.NewSet(
	NewDraw,
	wire.Bind(new(UsecaseDraw), new(*draw)),
)
