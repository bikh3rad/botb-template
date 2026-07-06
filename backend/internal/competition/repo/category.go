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

// isUniqueViolation matches Postgres unique_violation (SQLSTATE 23505).
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError

	return errors.As(err, &pgErr) && pgErr.Code == uniqueViolation
}

type category struct {
	logger *slog.Logger
	tracer trace.Tracer
	db     *datasource.PostgresDB
}

var _ biz.RepositoryCategory = (*category)(nil)

// NewCategory constructs the pgx-backed category repository.
func NewCategory(logger *slog.Logger, db *datasource.PostgresDB) *category {
	return &category{
		logger: logger.With("layer", "CategoryRepo"),
		tracer: otel.Tracer("CategoryRepo"),
		db:     db,
	}
}

const categoryColumns = `id, name, slug, created_at`

func (r *category) List(ctx context.Context) ([]entity.Category, error) {
	logger := r.logger.With("method", "List")

	rows, err := r.db.QueryContext(ctx,
		`SELECT `+categoryColumns+` FROM categories ORDER BY name`)
	if err != nil {
		logger.WarnContext(ctx, "failed to query categories", "error", err)

		return nil, err
	}
	defer rows.Close()

	out := []entity.Category{}

	for rows.Next() {
		var c entity.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.CreatedAt); err != nil {
			logger.WarnContext(ctx, "failed to scan category", "error", err)

			continue
		}

		out = append(out, c)
	}

	return out, rows.Err()
}

func (r *category) Get(ctx context.Context, id uuid.UUID) (entity.Category, error) {
	var c entity.Category

	err := r.db.QueryRowContext(
		ctx,
		`SELECT `+categoryColumns+` FROM categories WHERE id = $1`, id,
	).Scan(&c.ID, &c.Name, &c.Slug, &c.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Category{}, biz.ErrCategoryNotFound
		}

		return entity.Category{}, err
	}

	return c, nil
}

func (r *category) Create(ctx context.Context, c entity.Category) (entity.Category, error) {
	logger := r.logger.With("method", "Create")

	row := r.db.QueryRowContext(ctx,
		`INSERT INTO categories (id, name, slug) VALUES ($1, $2, $3) RETURNING created_at`,
		c.ID, c.Name, c.Slug)
	if err := row.Scan(&c.CreatedAt); err != nil {
		if isUniqueViolation(err) {
			return entity.Category{}, biz.ErrResourceExists
		}

		logger.WarnContext(ctx, "failed to insert category", "error", err)

		return entity.Category{}, err
	}

	return c, nil
}

func (r *category) Update(ctx context.Context, c entity.Category) (entity.Category, error) {
	logger := r.logger.With("method", "Update")

	res, err := r.db.ExecContext(ctx,
		`UPDATE categories SET name = $1, slug = $2 WHERE id = $3`, c.Name, c.Slug, c.ID)
	if err != nil {
		if isUniqueViolation(err) {
			return entity.Category{}, biz.ErrResourceExists
		}

		logger.WarnContext(ctx, "failed to update category", "error", err)

		return entity.Category{}, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return entity.Category{}, err
	}

	if affected == 0 {
		return entity.Category{}, biz.ErrCategoryNotFound
	}

	return c, nil
}

// Delete removes a category. The in-use guard and the optional reassignment
// happen inside ONE transaction so competitions can never be orphaned: either
// they are all moved to reassignTo and the category goes away, or nothing
// changes at all.
func (r *category) Delete(ctx context.Context, id uuid.UUID, reassignTo *uuid.UUID) error {
	logger := r.logger.With("method", "Delete")

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() { _ = tx.Rollback() }()

	var inUse int
	if err := tx.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM competitions WHERE category_id = $1`, id).Scan(&inUse); err != nil {
		return err
	}

	if inUse > 0 {
		if reassignTo == nil {
			return biz.ErrCategoryInUse
		}

		if _, err := tx.ExecContext(ctx,
			`UPDATE competitions SET category_id = $1 WHERE category_id = $2`, *reassignTo, id); err != nil {
			logger.WarnContext(ctx, "failed to reassign competitions", "error", err)

			return err
		}
	}

	res, err := tx.ExecContext(ctx, `DELETE FROM categories WHERE id = $1`, id)
	if err != nil {
		logger.WarnContext(ctx, "failed to delete category", "error", err)

		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return biz.ErrCategoryNotFound
	}

	return tx.Commit()
}
