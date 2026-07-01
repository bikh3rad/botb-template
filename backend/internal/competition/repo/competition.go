package repo

import (
	"application/internal/competition/biz"
	"application/internal/competition/entity"
	"application/internal/datasource"
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// uniqueViolation is the Postgres SQLSTATE for a unique_violation.
const uniqueViolation = "23505"

type competition struct {
	logger *slog.Logger
	tracer trace.Tracer
	db     *datasource.PostgresDB
}

var _ biz.Repository = (*competition)(nil)

// NewCompetition constructs the pgx-backed competition repository.
func NewCompetition(logger *slog.Logger, db *datasource.PostgresDB) *competition {
	return &competition{
		logger: logger.With("layer", "CompetitionRepo"),
		tracer: otel.Tracer("CompetitionRepo"),
		db:     db,
	}
}

const competitionColumns = `id, title, slug, description, prize, ticket_price_pence,
	tickets_total, tickets_sold, status, starts_at, ends_at, created_at, updated_at`

// Create inserts a competition (id pre-generated) and returns the stored row.
func (r *competition) Create(ctx context.Context, c entity.Competition) (entity.Competition, error) {
	logger := r.logger.With("method", "Create")

	query := `INSERT INTO competitions
		(id, title, slug, description, prize, ticket_price_pence, tickets_total,
		 tickets_sold, status, starts_at, ends_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		RETURNING created_at, updated_at`

	row := r.db.QueryRowContext(ctx, query,
		c.ID, c.Title, c.Slug, c.Description, c.Prize, c.TicketPricePence,
		c.TicketsTotal, c.TicketsSold, string(c.Status), c.StartsAt, c.EndsAt,
	)

	if err := row.Scan(&c.CreatedAt, &c.UpdatedAt); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return entity.Competition{}, biz.ErrResourceExists
		}

		logger.WarnContext(ctx, "failed to insert competition", "error", err)

		return entity.Competition{}, err
	}

	c.Media = []entity.MediaRef{}

	return c, nil
}

// Get returns a competition by id with its media attached.
func (r *competition) Get(ctx context.Context, id uuid.UUID) (entity.Competition, error) {
	logger := r.logger.With("method", "Get")

	query := `SELECT ` + competitionColumns + ` FROM competitions WHERE id = $1`

	c, err := scanCompetition(r.db.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Competition{}, biz.ErrResourceNotFound
		}

		logger.WarnContext(ctx, "failed to scan competition", "error", err)

		return entity.Competition{}, err
	}

	media, err := r.mediaByOwner(ctx, c.ID)
	if err != nil {
		return entity.Competition{}, err
	}

	c.Media = media

	return c, nil
}

// List returns competitions (optionally filtered by status), newest first, each
// with its media attached.
func (r *competition) List(ctx context.Context, filter biz.ListFilter) ([]entity.Competition, error) {
	logger := r.logger.With("method", "List")

	query := `SELECT ` + competitionColumns + ` FROM competitions`

	args := []any{}

	if filter.Status != nil {
		query += ` WHERE status = $1`

		args = append(args, string(*filter.Status))
	}

	query += ` ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		logger.WarnContext(ctx, "failed to query competitions", "error", err)

		return nil, err
	}
	defer rows.Close()

	var out []entity.Competition

	for rows.Next() {
		c, err := scanCompetition(rows)
		if err != nil {
			logger.WarnContext(ctx, "failed to scan competition row", "error", err)

			continue
		}

		out = append(out, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Attach media per competition. This is N+1 by competition; acceptable for
	// the modest listing sizes here and keeps the query ramsql-testable. Swap
	// for a single `owner_id = ANY($1)` batch query if listings grow large.
	for i := range out {
		media, err := r.mediaByOwner(ctx, out[i].ID)
		if err != nil {
			return nil, err
		}

		out[i].Media = media
	}

	return out, nil
}

// Update writes the editable fields and returns the refreshed row with media.
func (r *competition) Update(ctx context.Context, c entity.Competition) (entity.Competition, error) {
	logger := r.logger.With("method", "Update")

	query := `UPDATE competitions SET
		title = $2, description = $3, prize = $4, ticket_price_pence = $5,
		tickets_total = $6, status = $7, starts_at = $8, ends_at = $9,
		updated_at = NOW()
		WHERE id = $1
		RETURNING ` + competitionColumns

	updated, err := scanCompetition(r.db.QueryRowContext(ctx, query,
		c.ID, c.Title, c.Description, c.Prize, c.TicketPricePence,
		c.TicketsTotal, string(c.Status), c.StartsAt, c.EndsAt,
	))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Competition{}, biz.ErrResourceNotFound
		}

		logger.WarnContext(ctx, "failed to update competition", "error", err)

		return entity.Competition{}, err
	}

	media, err := r.mediaByOwner(ctx, updated.ID)
	if err != nil {
		return entity.Competition{}, err
	}

	updated.Media = media

	return updated, nil
}

// Delete removes a competition, returning ErrResourceNotFound if absent.
func (r *competition) Delete(ctx context.Context, id uuid.UUID) error {
	logger := r.logger.With("method", "Delete")

	result, err := r.db.ExecContext(ctx, `DELETE FROM competitions WHERE id = $1`, id)
	if err != nil {
		logger.WarnContext(ctx, "failed to delete competition", "error", err)

		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return biz.ErrResourceNotFound
	}

	return nil
}

// mediaByOwner reads the shared `media` table for this competition's media.
// Design note: we query the shared Postgres directly rather than calling the
// media service over HTTP — the template uses a single pgx datasource, so a
// read query is the most consistent, lowest-latency choice. It can be swapped
// for an HTTP media client later without changing the biz/handler layers.
func (r *competition) mediaByOwner(ctx context.Context, ownerID uuid.UUID) ([]entity.MediaRef, error) {
	query := `SELECT id, kind, bucket, object_key, content_type, position
		FROM media WHERE owner_type = 'competition' AND owner_id = $1
		ORDER BY position ASC, created_at ASC`

	rows, err := r.db.QueryContext(ctx, query, ownerID)
	if err != nil {
		r.logger.WarnContext(ctx, "failed to query media", "error", err)

		return nil, err
	}
	defer rows.Close()

	media := []entity.MediaRef{}

	for rows.Next() {
		var m entity.MediaRef
		if err := rows.Scan(&m.ID, &m.Kind, &m.Bucket, &m.ObjectKey, &m.ContentType, &m.Position); err != nil {
			r.logger.WarnContext(ctx, "failed to scan media ref", "error", err)

			continue
		}

		media = append(media, m)
	}

	return media, rows.Err()
}

type scanner interface {
	Scan(dest ...any) error
}

func scanCompetition(s scanner) (entity.Competition, error) {
	var (
		c      entity.Competition
		status string
	)

	if err := s.Scan(
		&c.ID, &c.Title, &c.Slug, &c.Description, &c.Prize, &c.TicketPricePence,
		&c.TicketsTotal, &c.TicketsSold, &status, &c.StartsAt, &c.EndsAt,
		&c.CreatedAt, &c.UpdatedAt,
	); err != nil {
		return entity.Competition{}, err
	}

	c.Status = entity.Status(status)
	c.Media = []entity.MediaRef{}

	return c, nil
}
