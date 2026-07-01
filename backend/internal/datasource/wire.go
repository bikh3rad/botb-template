package datasource

import "github.com/google/wire"

// PostgresProviderSet wires the config-driven Postgres client.
var PostgresProviderSet = wire.NewSet(
	NewPostgresConfig,
	NewPostgresDB,
)

// MinioProviderSet wires the S3-compatible MinIO object storage client.
var MinioProviderSet = wire.NewSet(
	NewMinioConfig,
	NewMinioStorage,
)

// DataProviderSet is the original template set (in-memory + Postgres) used by
// cmd/app. Individual services compose the leaner sets above as needed.
var DataProviderSet = wire.NewSet(
	NewInmemoryDB,
	PostgresProviderSet,
)
