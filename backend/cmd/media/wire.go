//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"context"

	"application/app"
	flatbiz "application/internal/biz"
	"application/internal/datasource"
	mediabiz "application/internal/media/biz"
	mediahandler "application/internal/media/handler"
	mediarepo "application/internal/media/repo"
	"application/internal/service"
	svchandler "application/internal/service/handler"
	"application/pkg/middlewares"

	"github.com/google/wire"
)

// wireApp is the composition root for the media service binary. It reuses the
// shared app + HTTP infrastructure, adds Postgres + MinIO datasources, the
// shared healthz endpoints, and the media domain stack.
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

		// Shared JWT auth (media mutations moved under /admin/ — the old
		// unauthenticated upload is gone).
		middlewares.JWTProviderSet,

		// Media domain.
		mediabiz.ProviderSet,
		mediarepo.ProviderSet,
		mediahandler.ProviderSet,

		// Bind the MinIO client to the media biz storage port.
		wire.Bind(new(mediabiz.ObjectStorage), new(*datasource.MinioStorage)),
	))
}
