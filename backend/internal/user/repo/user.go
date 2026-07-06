package repo

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"strconv"

	"application/internal/datasource"
	"application/internal/user/biz"
	"application/internal/user/entity"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const uniqueViolation = "23505"

type user struct {
	logger *slog.Logger
	tracer trace.Tracer
	db     *datasource.PostgresDB
}

var _ biz.RepositoryUser = (*user)(nil)

// NewUser constructs the pgx-backed user repository.
func NewUser(logger *slog.Logger, db *datasource.PostgresDB) *user {
	return &user{
		logger: logger.With("layer", "UserRepo"),
		tracer: otel.Tracer("UserRepo"),
		db:     db,
	}
}

const userColumns = `id, name, email, tickets_owned, total_spent_pence, is_active, created_at`

// Create inserts a user (id pre-generated) and returns the stored row.
func (r *user) Create(ctx context.Context, u entity.User) (entity.User, error) {
	logger := r.logger.With("method", "Create")

	query := `INSERT INTO users (id, name, email) VALUES ($1, $2, $3)
		RETURNING tickets_owned, total_spent_pence, is_active, created_at`

	row := r.db.QueryRowContext(ctx, query, u.ID, u.Name, u.Email)
	if err := row.Scan(&u.TicketsOwned, &u.TotalSpentPence, &u.IsActive, &u.CreatedAt); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return entity.User{}, biz.ErrResourceExists
		}

		logger.WarnContext(ctx, "failed to insert user", "error", err)

		return entity.User{}, err
	}

	return u, nil
}

// Get returns a user by id, mapping a missing row to ErrResourceNotFound.
func (r *user) Get(ctx context.Context, id uuid.UUID) (entity.User, error) {
	logger := r.logger.With("method", "Get")

	query := `SELECT ` + userColumns + ` FROM users WHERE id = $1`

	u, err := scanUser(r.db.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, biz.ErrResourceNotFound
		}

		logger.WarnContext(ctx, "failed to scan user", "error", err)

		return entity.User{}, err
	}

	return u, nil
}

// List returns a page of users (optionally filtered by a name/email substring)
// plus the total match count.
func (r *user) List(ctx context.Context, filter biz.UserListFilter) (biz.UserPage, error) {
	logger := r.logger.With("method", "List")

	where := ""
	args := []any{}

	if filter.Query != "" {
		where = ` WHERE name ILIKE $1 OR email ILIKE $1`

		args = append(args, "%"+filter.Query+"%")
	}

	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`+where, args...).Scan(&total); err != nil {
		logger.WarnContext(ctx, "failed to count users", "error", err)

		return biz.UserPage{}, err
	}

	// Limit/Offset are bounded ints (capped in the use case), so they are inlined
	// as literals rather than bound parameters — portable across drivers.
	query := `SELECT ` + userColumns + ` FROM users` + where +
		` ORDER BY created_at DESC LIMIT ` + strconv.Itoa(filter.Limit) +
		` OFFSET ` + strconv.Itoa(filter.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		logger.WarnContext(ctx, "failed to query users", "error", err)

		return biz.UserPage{}, err
	}
	defer rows.Close()

	users := []entity.User{}

	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			logger.WarnContext(ctx, "failed to scan user row", "error", err)

			continue
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return biz.UserPage{}, err
	}

	return biz.UserPage{Users: users, Total: total}, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanUser(s scanner) (entity.User, error) {
	var u entity.User
	if err := s.Scan(&u.ID, &u.Name, &u.Email, &u.TicketsOwned, &u.TotalSpentPence, &u.IsActive, &u.CreatedAt); err != nil {
		return entity.User{}, err
	}

	return u, nil
}

// Update writes the admin-editable profile fields (name, email). The derived
// counters (tickets_owned, total_spent_pence) are intentionally NOT part of
// this statement — they are owned by the purchase transaction. UPDATE-then-
// SELECT (no RETURNING) keeps the statement ramsql-testable.
func (r *user) Update(ctx context.Context, id uuid.UUID, name, email string) (entity.User, error) {
	logger := r.logger.With("method", "Update")

	res, err := r.db.ExecContext(ctx,
		`UPDATE users SET name = $1, email = $2 WHERE id = $3`, name, email, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return entity.User{}, biz.ErrResourceExists
		}

		logger.WarnContext(ctx, "failed to update user", "error", err)

		return entity.User{}, err
	}

	if affected, err := res.RowsAffected(); err != nil {
		return entity.User{}, err
	} else if affected == 0 {
		return entity.User{}, biz.ErrResourceNotFound
	}

	return r.Get(ctx, id)
}

// SetActive flips the suspend flag.
func (r *user) SetActive(ctx context.Context, id uuid.UUID, active bool) (entity.User, error) {
	logger := r.logger.With("method", "SetActive")

	res, err := r.db.ExecContext(ctx,
		`UPDATE users SET is_active = $1 WHERE id = $2`, active, id)
	if err != nil {
		logger.WarnContext(ctx, "failed to set user active flag", "error", err)

		return entity.User{}, err
	}

	if affected, err := res.RowsAffected(); err != nil {
		return entity.User{}, err
	} else if affected == 0 {
		return entity.User{}, biz.ErrResourceNotFound
	}

	return r.Get(ctx, id)
}
