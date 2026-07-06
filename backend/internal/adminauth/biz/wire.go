package biz

import "github.com/google/wire"

// ProviderSet wires the adminauth use cases + config.
var ProviderSet = wire.NewSet(
	NewConfig,

	NewAuth,
	wire.Bind(new(UsecaseAuth), new(*auth)),

	NewAccounts,
	wire.Bind(new(UsecaseAccounts), new(*accounts)),
)
