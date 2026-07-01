package repo

import (
	"application/internal/competition/biz"

	"github.com/google/wire"
)

// ProviderSet wires the competition repository and binds it to biz.Repository.
var ProviderSet = wire.NewSet(
	NewCompetition,
	wire.Bind(new(biz.Repository), new(*competition)),
)
