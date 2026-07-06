package biz

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"log/slog"
	"strings"
	"time"

	"application/internal/adminauth/entity"
	"application/pkg/middlewares"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
)

// dummyBcryptHash is compared against for UNKNOWN emails so a login attempt
// takes roughly the same time whether or not the account exists (no user
// enumeration via timing). It is the hash of an unguessable random string and
// grants nothing.
const dummyBcryptHash = "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"

type auth struct {
	logger  *slog.Logger
	tracer  trace.Tracer
	repo    Repository
	secret  []byte
	cfg     *Config
	limiter *loginLimiter
	now     func() time.Time
}

var _ UsecaseAuth = (*auth)(nil)

// NewAuth constructs the auth use case. It signs access tokens with the SAME
// shared HS256 secret every service's guard verifies (`jwt.secret`).
func NewAuth(logger *slog.Logger, repo Repository, secret middlewares.JWTSecret, cfg *Config) *auth {
	return &auth{
		logger:  logger.With("layer", "AdminAuth"),
		tracer:  otel.Tracer("AdminAuthUseCase"),
		repo:    repo,
		secret:  []byte(secret),
		cfg:     cfg,
		limiter: newLoginLimiter(),
		now:     time.Now,
	}
}

// Login verifies credentials and issues an access+refresh pair. All failure
// modes (unknown email, wrong password, disabled account) return the same
// generic ErrInvalidCredentials. Passwords and tokens are never logged.
func (uc *auth) Login(ctx context.Context, email, password, clientIP string) (LoginResult, error) {
	ctx, span := uc.tracer.Start(ctx, "Login")
	defer span.End()

	email = strings.TrimSpace(strings.ToLower(email))

	if !uc.limiter.allow("email:"+email, "ip:"+clientIP) {
		uc.logger.WarnContext(ctx, "login rate limited", "email", email)

		return LoginResult{}, ErrRateLimited
	}

	account, err := uc.repo.GetAccountByEmail(ctx, email)
	if err != nil {
		// Burn a bcrypt compare anyway to keep timing flat for unknown emails.
		_ = bcrypt.CompareHashAndPassword([]byte(dummyBcryptHash), []byte(password))

		return LoginResult{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(password)); err != nil {
		return LoginResult{}, ErrInvalidCredentials
	}

	if !account.IsActive {
		return LoginResult{}, ErrInvalidCredentials
	}

	now := uc.now().UTC()
	if err := uc.repo.TouchLastLogin(ctx, account.ID, now); err != nil {
		uc.logger.WarnContext(ctx, "failed to update last_login_at", "error", err)
	}

	account.LastLoginAt = &now

	return uc.issue(ctx, account)
}

// Refresh rotates a refresh token: the presented token is marked rotated and a
// new pair is issued. Presenting an ALREADY-ROTATED token is treated as theft
// (the legitimate holder rotated it) — every session for that admin is revoked.
func (uc *auth) Refresh(ctx context.Context, refreshToken string) (LoginResult, error) {
	ctx, span := uc.tracer.Start(ctx, "Refresh")
	defer span.End()

	stored, err := uc.repo.GetRefreshToken(ctx, hashToken(refreshToken))
	if err != nil {
		return LoginResult{}, ErrInvalidRefresh
	}

	now := uc.now().UTC()

	if stored.RotatedAt != nil {
		uc.logger.WarnContext(ctx, "reused refresh token — revoking all sessions", "admin_id", stored.AdminID)
		_ = uc.repo.RevokeAllForAdmin(ctx, stored.AdminID, now)

		return LoginResult{}, ErrInvalidRefresh
	}

	if stored.RevokedAt != nil || now.After(stored.ExpiresAt) {
		return LoginResult{}, ErrInvalidRefresh
	}

	account, err := uc.repo.GetAccount(ctx, stored.AdminID)
	if err != nil || !account.IsActive {
		return LoginResult{}, ErrInvalidRefresh
	}

	if err := uc.repo.MarkRefreshRotated(ctx, stored.ID, now); err != nil {
		return LoginResult{}, err
	}

	return uc.issue(ctx, account)
}

// Logout revokes the presented refresh token. Idempotent: unknown tokens are
// not an error (nothing to revoke).
func (uc *auth) Logout(ctx context.Context, refreshToken string) error {
	ctx, span := uc.tracer.Start(ctx, "Logout")
	defer span.End()

	err := uc.repo.RevokeRefreshToken(ctx, hashToken(refreshToken), uc.now().UTC())
	if err != nil && !errors.Is(err, ErrResourceNotFound) {
		return err
	}

	return nil
}

// Me returns the account behind a verified token subject. A missing or
// disabled account maps to ErrInvalidCredentials so the handler returns 401
// (the token outlived the account).
func (uc *auth) Me(ctx context.Context, adminID string) (entity.AdminAccount, error) {
	id, err := uuid.Parse(adminID)
	if err != nil {
		return entity.AdminAccount{}, ErrInvalidCredentials
	}

	account, err := uc.repo.GetAccount(ctx, id)
	if err != nil || !account.IsActive {
		return entity.AdminAccount{}, ErrInvalidCredentials
	}

	return account, nil
}

// issue creates the access JWT + a fresh stored refresh token.
func (uc *auth) issue(ctx context.Context, account entity.AdminAccount) (LoginResult, error) {
	now := uc.now().UTC()
	expiresAt := now.Add(uc.cfg.AccessTTL)

	claims := jwt.MapClaims{
		"sub":   account.ID.String(),
		"email": account.Email,
		"role":  string(account.Role),
		"iat":   now.Unix(),
		"exp":   expiresAt.Unix(),
	}

	access, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(uc.secret)
	if err != nil {
		return LoginResult{}, err
	}

	raw, err := newRefreshToken()
	if err != nil {
		return LoginResult{}, err
	}

	stored := entity.RefreshToken{
		ID:        uuid.New(),
		AdminID:   account.ID,
		TokenHash: hashToken(raw),
		ExpiresAt: now.Add(uc.cfg.RefreshTTL),
		CreatedAt: now,
	}

	if err := uc.repo.CreateRefreshToken(ctx, stored); err != nil {
		return LoginResult{}, err
	}

	return LoginResult{
		AccessToken:  access,
		ExpiresIn:    int64(uc.cfg.AccessTTL.Seconds()),
		RefreshToken: raw,
		Admin:        account,
	}, nil
}

// newRefreshToken returns 32 bytes of crypto/rand as base64url.
func newRefreshToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(buf), nil
}

// hashToken stores refresh tokens as SHA-256 hex — a DB leak must not leak
// usable tokens.
func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))

	return hex.EncodeToString(sum[:])
}
