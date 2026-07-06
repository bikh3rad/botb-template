//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"context"

	"application/app"
	adminauthbiz "application/internal/adminauth/biz"
	adminauthhandler "application/internal/adminauth/handler"
	adminauthrepo "application/internal/adminauth/repo"
	flatbiz "application/internal/biz"
	"application/internal/datasource"
	"application/internal/service"
	svchandler "application/internal/service/handler"
	"application/pkg/middlewares"

	"github.com/google/wire"
)

// wireApp is the composition root for the adminauth service binary. It reuses
// the shared app + HTTP infrastructure and Postgres datasource, adds the
// shared healthz endpoints, and the adminauth domain stack.
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

		// Shared JWT secret + guards (adminauth SIGNS with the same secret).
		middlewares.JWTProviderSet,

		// Adminauth domain.
		adminauthbiz.ProviderSet,
		adminauthrepo.ProviderSet,
		adminauthhandler.ProviderSet,
	))
}
