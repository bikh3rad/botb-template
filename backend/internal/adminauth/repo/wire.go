package repo

import (
	"application/internal/adminauth/biz"

	"github.com/google/wire"
)

// ProviderSet wires the adminauth repository to its biz interface.
var ProviderSet = wire.NewSet(
	NewAdmin,
	wire.Bind(new(biz.Repository), new(*admin)),
)
