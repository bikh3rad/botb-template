package repo

import (
	"application/internal/user/biz"

	"github.com/google/wire"
)

// ProviderSet wires the user and ticket repositories and binds them to their
// biz repository interfaces.
var ProviderSet = wire.NewSet(
	NewUser,
	wire.Bind(new(biz.RepositoryUser), new(*user)),

	NewTicket,
	wire.Bind(new(biz.RepositoryTicket), new(*ticket)),
)
