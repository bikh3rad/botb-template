package repo

import (
	"context"
	"log/slog"
	"strconv"

	"application/internal/competition/biz"
	"application/internal/competition/entity"
	"application/internal/datasource"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type auditRepo struct {
	logger *slog.Logger
	tracer trace.Tracer
	db     *datasource.PostgresDB
}

var _ biz.RepositoryAudit = (*auditRepo)(nil)

// NewAuditRepo constructs the pgx-backed audit-log reader.
func NewAuditRepo(logger *slog.Logger, db *datasource.PostgresDB) *auditRepo {
	return &auditRepo{
		logger: logger.With("layer", "AuditRepo"),
		tracer: otel.Tracer("AuditRepo"),
		db:     db,
	}
}

func (r *auditRepo) Recent(ctx context.Context, limit int) ([]entity.AuditEntry, error) {
	logger := r.logger.With("method", "Recent")

	// limit is bounded by the use case; inlined as a literal for portability.
	query := `SELECT id, actor_id, actor_email, action, entity_type, entity_id, reason, created_at
		FROM admin_audit_log ORDER BY created_at DESC LIMIT ` + strconv.Itoa(limit)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		logger.WarnContext(ctx, "failed to query audit log", "error", err)

		return nil, err
	}
	defer rows.Close()

	out := []entity.AuditEntry{}

	for rows.Next() {
		var e entity.AuditEntry
		if err := rows.Scan(&e.ID, &e.ActorID, &e.ActorEmail, &e.Action,
			&e.EntityType, &e.EntityID, &e.Reason, &e.CreatedAt); err != nil {
			logger.WarnContext(ctx, "failed to scan audit row", "error", err)

			continue
		}

		out = append(out, e)
	}

	return out, rows.Err()
}
