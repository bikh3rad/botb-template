package biz

import "github.com/google/wire"

// ProviderSet wires the media use case and binds it to UsecaseMedia.
var ProviderSet = wire.NewSet(
	NewMedia,
	wire.Bind(new(UsecaseMedia), new(*media)),
)
