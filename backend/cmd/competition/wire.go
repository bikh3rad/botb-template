//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"context"

	"application/app"
	flatbiz "application/internal/biz"
	compbiz "application/internal/competition/biz"
	comphandler "application/internal/competition/handler"
	comprepo "application/internal/competition/repo"
	"application/internal/datasource"
	"application/internal/service"
	svchandler "application/internal/service/handler"
	"application/pkg/middlewares"

	"github.com/google/wire"
)

// wireApp is the composition root for the competition service binary. It reuses
// the shared app + HTTP infrastructure and Postgres datasource, adds the shared
// healthz endpoints, and the competition domain stack. It also wires the MinIO
// object-storage client so a permitted competition delete can purge that
// competition's media objects (the media rows are cleaned up in the same DB
// transaction; the objects are removed via this port).
func wireApp(
	ctx context.Context,
) (app.Application, error) {
	panic(wire.Build(
		app.AppProviderSet,
		service.ServerProviderSet,

		datasource.PostgresProviderSet,
		datasource.MinioProviderSet,

		// Shared healthz endpoints.
		flatbiz.HealthzProviderSet,
		svchandler.NewMuxHealthzHandler,

		// Shared JWT auth (admin route guard — defense in depth).
		middlewares.JWTProviderSet,

		// Competition domain.
		compbiz.ProviderSet,
		comprepo.ProviderSet,
		comphandler.ProviderSet,

		// The competition delete path purges media objects via MinIO.
		wire.Bind(new(compbiz.ObjectStorage), new(*datasource.MinioStorage)),
	))
}
