package repo

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"application/internal/adminauth/biz"
	"application/internal/adminauth/entity"
	"application/internal/datasource"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const uniqueViolation = "23505"

type admin struct {
	logger *slog.Logger
	tracer trace.Tracer
	db     *datasource.PostgresDB
}

var _ biz.Repository = (*admin)(nil)

// NewAdmin constructs the pgx-backed adminauth repository. All tables live in
// the dedicated `adminauth` schema — never in public.
func NewAdmin(logger *slog.Logger, db *datasource.PostgresDB) *admin {
	return &admin{
		logger: logger.With("layer", "AdminAuthRepo"),
		tracer: otel.Tracer("AdminAuthRepo"),
		db:     db,
	}
}

const accountColumns = `id, name, email, password_hash, role, is_active, created_at, last_login_at`

func (r *admin) CreateAccount(ctx context.Context, a entity.AdminAccount) (entity.AdminAccount, error) {
	logger := r.logger.With("method", "CreateAccount")

	query := `INSERT INTO adminauth.admin_accounts (id, name, email, password_hash, role, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at`

	row := r.db.QueryRowContext(ctx, query,
		a.ID, a.Name, a.Email, a.PasswordHash, string(a.Role), a.IsActive)
	if err := row.Scan(&a.CreatedAt); err != nil {
		if isUniqueViolation(err) {
			return entity.AdminAccount{}, biz.ErrResourceExists
		}

		logger.WarnContext(ctx, "failed to insert admin account", "error", err)

		return entity.AdminAccount{}, err
	}

	return a, nil
}

func (r *admin) GetAccount(ctx context.Context, id uuid.UUID) (entity.AdminAccount, error) {
	query := `SELECT ` + accountColumns + ` FROM adminauth.admin_accounts WHERE id = $1`

	return r.scanAccount(ctx, r.db.QueryRowContext(ctx, query, id))
}

func (r *admin) GetAccountByEmail(ctx context.Context, email string) (entity.AdminAccount, error) {
	query := `SELECT ` + accountColumns + ` FROM adminauth.admin_accounts WHERE email = $1`

	return r.scanAccount(ctx, r.db.QueryRowContext(ctx, query, email))
}

func (r *admin) ListAccounts(ctx context.Context) ([]entity.AdminAccount, error) {
	logger := r.logger.With("method", "ListAccounts")

	rows, err := r.db.QueryContext(ctx,
		`SELECT `+accountColumns+` FROM adminauth.admin_accounts ORDER BY created_at`)
	if err != nil {
		logger.WarnContext(ctx, "failed to query admin accounts", "error", err)

		return nil, err
	}
	defer rows.Close()

	accounts := []entity.AdminAccount{}

	for rows.Next() {
		a, err := scanAccountRow(rows)
		if err != nil {
			logger.WarnContext(ctx, "failed to scan admin account", "error", err)

			continue
		}

		accounts = append(accounts, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (r *admin) UpdateAccount(ctx context.Context, a entity.AdminAccount) (entity.AdminAccount, error) {
	logger := r.logger.With("method", "UpdateAccount")

	query := `UPDATE adminauth.admin_accounts
		SET name = $1, email = $2, password_hash = $3, role = $4, is_active = $5
		WHERE id = $6`

	res, err := r.db.ExecContext(ctx, query,
		a.Name, a.Email, a.PasswordHash, string(a.Role), a.IsActive, a.ID)
	if err != nil {
		if isUniqueViolation(err) {
			return entity.AdminAccount{}, biz.ErrResourceExists
		}

		logger.WarnContext(ctx, "failed to update admin account", "error", err)

		return entity.AdminAccount{}, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return entity.AdminAccount{}, err
	}

	if affected == 0 {
		return entity.AdminAccount{}, biz.ErrResourceNotFound
	}

	return a, nil
}

func (r *admin) CountActiveSuperadmins(ctx context.Context) (int, error) {
	var count int

	err := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM adminauth.admin_accounts WHERE role = 'superadmin' AND is_active`,
	).Scan(&count)

	return count, err
}

func (r *admin) TouchLastLogin(ctx context.Context, id uuid.UUID, at time.Time) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE adminauth.admin_accounts SET last_login_at = $1 WHERE id = $2`, at, id)

	return err
}

func (r *admin) CreateRefreshToken(ctx context.Context, t entity.RefreshToken) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO adminauth.admin_refresh_tokens (id, admin_id, token_hash, expires_at, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		t.ID, t.AdminID, t.TokenHash, t.ExpiresAt, t.CreatedAt)

	return err
}

func (r *admin) GetRefreshToken(ctx context.Context, tokenHash string) (entity.RefreshToken, error) {
	var (
		t         entity.RefreshToken
		rotatedAt sql.NullTime
		revokedAt sql.NullTime
	)

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, admin_id, token_hash, expires_at, rotated_at, revoked_at, created_at
		 FROM adminauth.admin_refresh_tokens WHERE token_hash = $1`, tokenHash,
	).Scan(&t.ID, &t.AdminID, &t.TokenHash, &t.ExpiresAt, &rotatedAt, &revokedAt, &t.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.RefreshToken{}, biz.ErrResourceNotFound
		}

		return entity.RefreshToken{}, err
	}

	if rotatedAt.Valid {
		at := rotatedAt.Time
		t.RotatedAt = &at
	}

	if revokedAt.Valid {
		at := revokedAt.Time
		t.RevokedAt = &at
	}

	return t, nil
}

func (r *admin) MarkRefreshRotated(ctx context.Context, id uuid.UUID, at time.Time) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE adminauth.admin_refresh_tokens SET rotated_at = $1 WHERE id = $2`, at, id)

	return err
}

func (r *admin) RevokeRefreshToken(ctx context.Context, tokenHash string, at time.Time) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE adminauth.admin_refresh_tokens SET revoked_at = $1
		 WHERE token_hash = $2 AND revoked_at IS NULL`, at, tokenHash)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return biz.ErrResourceNotFound
	}

	return nil
}

func (r *admin) RevokeAllForAdmin(ctx context.Context, adminID uuid.UUID, at time.Time) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE adminauth.admin_refresh_tokens SET revoked_at = $1
		 WHERE admin_id = $2 AND revoked_at IS NULL`, at, adminID)

	return err
}

func (r *admin) scanAccount(ctx context.Context, row *sql.Row) (entity.AdminAccount, error) {
	a, err := scanAccountRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.AdminAccount{}, biz.ErrResourceNotFound
		}

		r.logger.WarnContext(ctx, "failed to scan admin account", "error", err)

		return entity.AdminAccount{}, err
	}

	return a, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanAccountRow(s scanner) (entity.AdminAccount, error) {
	var (
		a         entity.AdminAccount
		role      string
		lastLogin sql.NullTime
	)

	if err := s.Scan(&a.ID, &a.Name, &a.Email, &a.PasswordHash, &role,
		&a.IsActive, &a.CreatedAt, &lastLogin); err != nil {
		return entity.AdminAccount{}, err
	}

	a.Role = entity.Role(role)

	if lastLogin.Valid {
		at := lastLogin.Time
		a.LastLoginAt = &at
	}

	return a, nil
}

// isUniqueViolation matches Postgres unique_violation (SQLSTATE 23505),
// mirroring the user repo's pgconn error mapping.
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError

	return errors.As(err, &pgErr) && pgErr.Code == uniqueViolation
}
