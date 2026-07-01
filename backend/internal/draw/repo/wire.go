package repo

import (
	"application/internal/draw/biz"

	"github.com/google/wire"
)

// ProviderSet wires the draw repository and binds it to biz.Repository.
var ProviderSet = wire.NewSet(
	NewDraw,
	wire.Bind(new(biz.Repository), new(*draw)),
)
