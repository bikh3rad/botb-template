package biz

import "github.com/google/wire"

// ProviderSet wires the competition, category and site-content use cases.
var ProviderSet = wire.NewSet(
	NewCompetition,
	wire.Bind(new(UsecaseCompetition), new(*competition)),

	NewCategory,
	wire.Bind(new(UsecaseCategory), new(*category)),

	NewContent,
	wire.Bind(new(UsecaseContent), new(*content)),
)
