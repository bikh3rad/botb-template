package biz

import "github.com/google/wire"

// ProviderSet wires the competition use case.
var ProviderSet = wire.NewSet(
	NewCompetition,
	wire.Bind(new(UsecaseCompetition), new(*competition)),
)
