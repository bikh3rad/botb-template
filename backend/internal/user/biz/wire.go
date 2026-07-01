package biz

import "github.com/google/wire"

// ProviderSet wires the user and ticket use cases.
var ProviderSet = wire.NewSet(
	NewUser,
	wire.Bind(new(UsecaseUser), new(*user)),

	NewTicket,
	wire.Bind(new(UsecaseTicket), new(*ticket)),
)
