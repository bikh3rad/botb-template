//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"application/app"
	flatbiz "application/internal/biz"
	"application/internal/datasource"
	drawbiz "application/internal/draw/biz"
	drawhandler "application/internal/draw/handler"
	drawrepo "application/internal/draw/repo"
	"application/internal/service"
	svchandler "application/internal/service/handler"
	"context"

	"github.com/google/wire"
)

// wireApp is the composition root for the draw service binary. It reuses the
// shared app + HTTP infrastructure and Postgres datasource, adds the shared
// healthz endpoints, and the draw domain stack.
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

		// Draw domain.
		drawbiz.ProviderSet,
		drawrepo.ProviderSet,
		drawhandler.ProviderSet,
	))
}
