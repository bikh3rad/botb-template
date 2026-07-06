package biz

import (
	"context"
	"time"

	"application/internal/adminauth/entity"

	"github.com/google/uuid"
)

// LoginResult carries a freshly issued token pair plus the account it belongs
// to. The refresh token is the raw (unhashed) value — shown once, never stored.
type LoginResult struct {
	AccessToken  string
	ExpiresIn    int64 // seconds
	RefreshToken string
	Admin        entity.AdminAccount
}

// CreateAccountInput is the superadmin account-creation input.
type CreateAccountInput struct {
	Name     string
	Email    string
	Password string
	Role     entity.Role
}

// UpdateAccountInput is a partial account edit; nil fields are unchanged.
// There is intentionally no hard delete — accounts are disabled instead.
type UpdateAccountInput struct {
	Name     *string
	Email    *string
	Password *string
	Role     *entity.Role
	IsActive *bool
}

// UsecaseAuth is the login/session use case consumed by the HTTP handler.
type UsecaseAuth interface {
	Login(ctx context.Context, email, password, clientIP string) (LoginResult, error)
	Refresh(ctx context.Context, refreshToken string) (LoginResult, error)
	Logout(ctx context.Context, refreshToken string) error
	Me(ctx context.Context, adminID string) (entity.AdminAccount, error)
}

// UsecaseAccounts is the superadmin account-management use case.
type UsecaseAccounts interface {
	List(ctx context.Context) ([]entity.AdminAccount, error)
	Create(ctx context.Context, input CreateAccountInput) (entity.AdminAccount, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateAccountInput) (entity.AdminAccount, error)
}

// Repository persists admin accounts + refresh tokens (schema adminauth.*).
type Repository interface {
	CreateAccount(ctx context.Context, a entity.AdminAccount) (entity.AdminAccount, error)
	GetAccount(ctx context.Context, id uuid.UUID) (entity.AdminAccount, error)
	GetAccountByEmail(ctx context.Context, email string) (entity.AdminAccount, error)
	ListAccounts(ctx context.Context) ([]entity.AdminAccount, error)
	UpdateAccount(ctx context.Context, a entity.AdminAccount) (entity.AdminAccount, error)
	CountActiveSuperadmins(ctx context.Context) (int, error)
	TouchLastLogin(ctx context.Context, id uuid.UUID, at time.Time) error

	CreateRefreshToken(ctx context.Context, t entity.RefreshToken) error
	GetRefreshToken(ctx context.Context, tokenHash string) (entity.RefreshToken, error)
	MarkRefreshRotated(ctx context.Context, id uuid.UUID, at time.Time) error
	RevokeRefreshToken(ctx context.Context, tokenHash string, at time.Time) error
	RevokeAllForAdmin(ctx context.Context, adminID uuid.UUID, at time.Time) error
}
