package repo

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"application/internal/competition/biz"
	"application/internal/competition/entity"
	"application/internal/datasource"

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
	tickets_total, tickets_sold, category_id, status, starts_at, ends_at, created_at, updated_at`

// Create inserts a competition (id pre-generated) and returns the stored row.
func (r *competition) Create(ctx context.Context, c entity.Competition) (entity.Competition, error) {
	logger := r.logger.With("method", "Create")

	query := `INSERT INTO competitions
		(id, title, slug, description, prize, ticket_price_pence, tickets_total,
		 tickets_sold, category_id, status, starts_at, ends_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		RETURNING created_at, updated_at`

	row := r.db.QueryRowContext(
		ctx, query,
		c.ID, c.Title, c.Slug, c.Description, c.Prize, c.TicketPricePence,
		c.TicketsTotal, c.TicketsSold, nullableUUID(c.CategoryID), string(c.Status), c.StartsAt, c.EndsAt,
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

	cs := []entity.Competition{c}
	if err := r.fillCategoryNames(ctx, cs); err != nil {
		return entity.Competition{}, err
	}

	return cs[0], nil
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

	if err := r.fillCategoryNames(ctx, out); err != nil {
		return nil, err
	}

	return out, nil
}

// Update writes the editable fields and returns the refreshed row with media.
func (r *competition) Update(ctx context.Context, c entity.Competition) (entity.Competition, error) {
	logger := r.logger.With("method", "Update")

	// tickets_sold is deliberately NOT in this SET list — it is derived from
	// purchases and no admin endpoint may write it. slug IS updatable (full-
	// field editability); uniqueness maps to ErrResourceExists below.
	query := `UPDATE competitions SET
		title = $2, slug = $3, description = $4, prize = $5, ticket_price_pence = $6,
		tickets_total = $7, category_id = $8, status = $9, starts_at = $10, ends_at = $11,
		updated_at = NOW()
		WHERE id = $1
		RETURNING ` + competitionColumns

	updated, err := scanCompetition(r.db.QueryRowContext(
		ctx, query,
		c.ID, c.Title, c.Slug, c.Description, c.Prize, c.TicketPricePence,
		c.TicketsTotal, nullableUUID(c.CategoryID), string(c.Status), c.StartsAt, c.EndsAt,
	))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Competition{}, biz.ErrResourceNotFound
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return entity.Competition{}, biz.ErrResourceExists
		}

		logger.WarnContext(ctx, "failed to update competition", "error", err)

		return entity.Competition{}, err
	}

	media, err := r.mediaByOwner(ctx, updated.ID)
	if err != nil {
		return entity.Competition{}, err
	}

	updated.Media = media

	cs := []entity.Competition{updated}
	if err := r.fillCategoryNames(ctx, cs); err != nil {
		return entity.Competition{}, err
	}

	return cs[0], nil
}

// Delete removes a competition only when it has no entrants (zero sold tickets
// AND no draws), plus its media rows, all in one transaction. It returns the
// deleted media object keys so the use case can purge them from object storage.
// The guards live inside the transaction so a concurrent purchase/draw can't
// slip a competition out from under a paying entrant.
func (r *competition) Delete(ctx context.Context, id uuid.UUID) ([]string, error) {
	logger := r.logger.With("method", "Delete")

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() { _ = tx.Rollback() }()

	// The competition must exist.
	var soldOnRow int64
	if err := tx.QueryRowContext(ctx,
		`SELECT tickets_sold FROM competitions WHERE id = $1`, id).Scan(&soldOnRow); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, biz.ErrResourceNotFound
		}

		return nil, err
	}

	// Real tickets are authoritative (tickets_sold is a denormalized counter);
	// block on either. A single draw referencing it also blocks the delete.
	var ticketCount, drawCount int64
	if err := tx.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM tickets WHERE competition_id = $1`, id).Scan(&ticketCount); err != nil {
		return nil, err
	}

	if err := tx.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM draws WHERE competition_id = $1`, id).Scan(&drawCount); err != nil {
		return nil, err
	}

	if soldOnRow > 0 || ticketCount > 0 || drawCount > 0 {
		return nil, biz.ErrCompetitionHasEntrants
	}

	// Collect the media object keys, then delete the media rows and the
	// competition. Media has no FK to competitions, so we remove its rows
	// explicitly to avoid orphaned records.
	objectKeys, err := deleteCompetitionMedia(ctx, tx, id)
	if err != nil {
		logger.WarnContext(ctx, "failed to delete competition media rows", "error", err)

		return nil, err
	}

	result, err := tx.ExecContext(ctx, `DELETE FROM competitions WHERE id = $1`, id)
	if err != nil {
		logger.WarnContext(ctx, "failed to delete competition", "error", err)

		return nil, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if affected == 0 {
		return nil, biz.ErrResourceNotFound
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return objectKeys, nil
}

// deleteCompetitionMedia removes a competition's media rows within tx and
// returns their object keys for downstream object-storage purge. The SELECT is
// fully drained + closed (in mediaObjectKeys) before the DELETE runs, so the
// single transaction connection is never busy across statements.
func deleteCompetitionMedia(ctx context.Context, tx *sql.Tx, id uuid.UUID) ([]string, error) {
	objectKeys, err := mediaObjectKeys(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	if _, err := tx.ExecContext(ctx,
		`DELETE FROM media WHERE owner_type = 'competition' AND owner_id = $1`, id); err != nil {
		return nil, err
	}

	return objectKeys, nil
}

// mediaObjectKeys returns the object keys of a competition's media rows. It
// fully drains + closes the result set (defer) so the caller can run further
// statements on the same transaction.
func mediaObjectKeys(ctx context.Context, tx *sql.Tx, id uuid.UUID) ([]string, error) {
	rows, err := tx.QueryContext(ctx,
		`SELECT object_key FROM media WHERE owner_type = 'competition' AND owner_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var objectKeys []string

	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, err
		}

		objectKeys = append(objectKeys, key)
	}

	return objectKeys, rows.Err()
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
		c          entity.Competition
		status     string
		categoryID uuid.NullUUID
	)

	if err := s.Scan(
		&c.ID, &c.Title, &c.Slug, &c.Description, &c.Prize, &c.TicketPricePence,
		&c.TicketsTotal, &c.TicketsSold, &categoryID, &status, &c.StartsAt, &c.EndsAt,
		&c.CreatedAt, &c.UpdatedAt,
	); err != nil {
		return entity.Competition{}, err
	}

	c.Status = entity.Status(status)
	c.Media = []entity.MediaRef{}

	if categoryID.Valid {
		id := categoryID.UUID
		c.CategoryID = &id
	}

	return c, nil
}

// nullableUUID maps a *uuid.UUID to a driver-friendly NULL when unset.
func nullableUUID(id *uuid.UUID) any {
	if id == nil {
		return nil
	}

	return *id
}

// fillCategoryNames resolves category display names in one small query (the
// categories table is tiny) instead of a JOIN, keeping the main queries
// ramsql-testable.
func (r *competition) fillCategoryNames(ctx context.Context, cs []entity.Competition) error {
	needed := false

	for i := range cs {
		if cs[i].CategoryID != nil {
			needed = true

			break
		}
	}

	if !needed {
		return nil
	}

	rows, err := r.db.QueryContext(ctx, `SELECT id, name FROM categories`)
	if err != nil {
		r.logger.WarnContext(ctx, "failed to query categories", "error", err)

		return err
	}
	defer rows.Close()

	names := map[uuid.UUID]string{}

	for rows.Next() {
		var (
			id   uuid.UUID
			name string
		)

		if err := rows.Scan(&id, &name); err != nil {
			continue
		}

		names[id] = name
	}

	if err := rows.Err(); err != nil {
		return err
	}

	for i := range cs {
		if cs[i].CategoryID != nil {
			cs[i].CategoryName = names[*cs[i].CategoryID]
		}
	}

	return nil
}
