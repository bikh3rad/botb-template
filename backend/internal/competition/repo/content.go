package repo

import (
	"context"
	"log/slog"

	"application/internal/competition/biz"
	"application/internal/competition/entity"
	"application/internal/datasource"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type content struct {
	logger *slog.Logger
	tracer trace.Tracer
	db     *datasource.PostgresDB
}

var _ biz.RepositoryContent = (*content)(nil)

// NewContent constructs the pgx-backed site-content repository.
func NewContent(logger *slog.Logger, db *datasource.PostgresDB) *content {
	return &content{
		logger: logger.With("layer", "SiteContentRepo"),
		tracer: otel.Tracer("SiteContentRepo"),
		db:     db,
	}
}

func (r *content) GetAll(ctx context.Context) ([]entity.SiteContent, error) {
	logger := r.logger.With("method", "GetAll")

	rows, err := r.db.QueryContext(ctx,
		`SELECT key, value, updated_at FROM site_content ORDER BY key`)
	if err != nil {
		logger.WarnContext(ctx, "failed to query site content", "error", err)

		return nil, err
	}
	defer rows.Close()

	out := []entity.SiteContent{}

	for rows.Next() {
		var c entity.SiteContent
		if err := rows.Scan(&c.Key, &c.Value, &c.UpdatedAt); err != nil {
			logger.WarnContext(ctx, "failed to scan site content", "error", err)

			continue
		}

		out = append(out, c)
	}

	return out, rows.Err()
}

func (r *content) Upsert(ctx context.Context, key, value string) (entity.SiteContent, error) {
	logger := r.logger.With("method", "Upsert")

	c := entity.SiteContent{Key: key, Value: value}

	row := r.db.QueryRowContext(ctx,
		`INSERT INTO site_content (key, value, updated_at) VALUES ($1, $2, NOW())
		 ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW()
		 RETURNING updated_at`, key, value)
	if err := row.Scan(&c.UpdatedAt); err != nil {
		logger.WarnContext(ctx, "failed to upsert site content", "error", err)

		return entity.SiteContent{}, err
	}

	return c, nil
}
