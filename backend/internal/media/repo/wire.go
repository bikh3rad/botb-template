package repo

import (
	"application/internal/media/biz"

	"github.com/google/wire"
)

// ProviderSet wires the media repository and binds it to biz.Repository.
var ProviderSet = wire.NewSet(
	NewMedia,
	wire.Bind(new(biz.Repository), new(*media)),
)
