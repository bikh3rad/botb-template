//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"application/app"
	flatbiz "application/internal/biz"
	"application/internal/gateway"
	"application/internal/service"
	svchandler "application/internal/service/handler"
	"application/pkg/middlewares"
	"context"

	"github.com/google/wire"
)

// wireApp is the composition root for the gateway binary. It reuses the shared
// app + HTTP infrastructure and shared healthz, adds the JWT auth middleware
// (admin route guard) and the reverse-proxy dispatcher. No datasource is wired —
// the gateway holds no state, it only proxies.
func wireApp(
	ctx context.Context,
) (app.Application, error) {
	panic(wire.Build(
		app.AppProviderSet,
		service.ServerProviderSet,

		// Shared healthz endpoints.
		flatbiz.HealthzProviderSet,
		svchandler.NewMuxHealthzHandler,

		// Shared JWT auth (reads jwt.secret from config).
		middlewares.JWTProviderSet,

		// Reverse-proxy dispatcher.
		gateway.ProviderSet,
	))
}
