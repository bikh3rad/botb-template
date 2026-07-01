//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"application/app"
	flatbiz "application/internal/biz"
	"application/internal/datasource"
	"application/internal/service"
	svchandler "application/internal/service/handler"
	"application/pkg/middlewares"
	userbiz "application/internal/user/biz"
	userhandler "application/internal/user/handler"
	userrepo "application/internal/user/repo"
	"context"

	"github.com/google/wire"
)

// wireApp is the composition root for the user service binary. It reuses the
// shared app + HTTP infrastructure and Postgres datasource, adds the shared
// healthz endpoints, and the user + ticket domain stack.
func wireApp(
	ctx context.Context,
) (app.Application, error) {
	panic(wire.Build(
		app.AppProviderSet,
		service.ServerProviderSet,

		datasource.PostgresProviderSet,

		// Shared healthz endpoints.
		flatbiz.HealthzProviderSet,
		svchandler.NewMuxHealthzHandler,

		// Shared JWT auth (admin route guard — defense in depth).
		middlewares.JWTProviderSet,

		// User + ticket domain.
		userbiz.ProviderSet,
		userrepo.ProviderSet,
		userhandler.ProviderSet,
	))
}
