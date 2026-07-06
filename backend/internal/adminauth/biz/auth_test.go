package biz_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log/slog"
	"testing"
	"time"

	"application/internal/adminauth/biz"
	"application/internal/adminauth/entity"
	"application/internal/adminauth/mocks"
	"application/pkg/middlewares"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

const testSecret = "adminauth-test-secret"

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func testConfig() *biz.Config {
	return &biz.Config{AccessTTL: 15 * time.Minute, RefreshTTL: 168 * time.Hour}
}

func newAuth(t *testing.T) (biz.UsecaseAuth, *mocks.MockRepository) {
	t.Helper()

	repo := mocks.NewMockRepository(t)

	return biz.NewAuth(discardLogger(), repo, middlewares.JWTSecret(testSecret), testConfig()), repo
}

func account(t *testing.T, password string, role entity.Role, active bool) entity.AdminAccount {
	t.Helper()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	require.NoError(t, err)

	return entity.AdminAccount{
		ID:           uuid.New(),
		Name:         "Test Admin",
		Email:        "admin@example.com",
		PasswordHash: string(hash),
		Role:         role,
		IsActive:     active,
	}
}

func TestLogin_Valid(t *testing.T) {
	uc, repo := newAuth(t)
	acc := account(t, "correct-password", entity.RoleSuperadmin, true)

	repo.EXPECT().GetAccountByEmail(mock.Anything, "admin@example.com").Return(acc, nil)
	repo.EXPECT().TouchLastLogin(mock.Anything, acc.ID, mock.Anything).Return(nil)
	repo.EXPECT().CreateRefreshToken(mock.Anything, mock.Anything).Return(nil)

	result, err := uc.Login(context.Background(), "Admin@Example.com", "correct-password", "1.2.3.4")
	require.NoError(t, err)
	require.NotEmpty(t, result.AccessToken)
	require.NotEmpty(t, result.RefreshToken)
	require.Equal(t, int64(900), result.ExpiresIn)

	// The access token must be verifiable with the shared secret and carry the
	// role claim the guards check.
	token, err := jwt.Parse(result.AccessToken, func(*jwt.Token) (any, error) {
		return []byte(testSecret), nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	require.NoError(t, err)

	claims := token.Claims.(jwt.MapClaims)
	require.Equal(t, "superadmin", claims["role"])
	require.Equal(t, acc.ID.String(), claims["sub"])
	require.Equal(t, acc.Email, claims["email"])
}

func TestLogin_WrongPassword(t *testing.T) {
	uc, repo := newAuth(t)
	acc := account(t, "correct-password", entity.RoleAdmin, true)
	repo.EXPECT().GetAccountByEmail(mock.Anything, "admin@example.com").Return(acc, nil)

	_, err := uc.Login(context.Background(), "admin@example.com", "wrong", "1.2.3.4")
	require.ErrorIs(t, err, biz.ErrInvalidCredentials)
}

func TestLogin_UnknownEmail(t *testing.T) {
	uc, repo := newAuth(t)
	repo.EXPECT().GetAccountByEmail(mock.Anything, "nobody@example.com").
		Return(entity.AdminAccount{}, biz.ErrResourceNotFound)

	_, err := uc.Login(context.Background(), "nobody@example.com", "whatever", "1.2.3.4")
	// Same generic error as a wrong password — no account enumeration.
	require.ErrorIs(t, err, biz.ErrInvalidCredentials)
}

func TestLogin_DisabledAccount(t *testing.T) {
	uc, repo := newAuth(t)
	acc := account(t, "correct-password", entity.RoleAdmin, false)
	repo.EXPECT().GetAccountByEmail(mock.Anything, "admin@example.com").Return(acc, nil)

	_, err := uc.Login(context.Background(), "admin@example.com", "correct-password", "1.2.3.4")
	require.ErrorIs(t, err, biz.ErrInvalidCredentials)
}

func TestLogin_RateLimited(t *testing.T) {
	uc, repo := newAuth(t)
	repo.EXPECT().GetAccountByEmail(mock.Anything, mock.Anything).
		Return(entity.AdminAccount{}, biz.ErrResourceNotFound).Times(5)

	for range 5 {
		_, err := uc.Login(context.Background(), "nobody@example.com", "x", "1.2.3.4")
		require.ErrorIs(t, err, biz.ErrInvalidCredentials)
	}

	// 6th attempt within the window: limiter trips before the repo is hit.
	_, err := uc.Login(context.Background(), "nobody@example.com", "x", "1.2.3.4")
	require.ErrorIs(t, err, biz.ErrRateLimited)
}

func hashOf(raw string) string {
	sum := sha256.Sum256([]byte(raw))

	return hex.EncodeToString(sum[:])
}

func TestRefresh_ValidRotates(t *testing.T) {
	uc, repo := newAuth(t)
	acc := account(t, "pw12345678", entity.RoleAdmin, true)
	raw := "some-refresh-token"

	stored := entity.RefreshToken{
		ID:        uuid.New(),
		AdminID:   acc.ID,
		TokenHash: hashOf(raw),
		ExpiresAt: time.Now().Add(time.Hour),
	}

	repo.EXPECT().GetRefreshToken(mock.Anything, hashOf(raw)).Return(stored, nil)
	repo.EXPECT().GetAccount(mock.Anything, acc.ID).Return(acc, nil)
	repo.EXPECT().MarkRefreshRotated(mock.Anything, stored.ID, mock.Anything).Return(nil)
	repo.EXPECT().CreateRefreshToken(mock.Anything, mock.Anything).Return(nil)

	result, err := uc.Refresh(context.Background(), raw)
	require.NoError(t, err)
	require.NotEmpty(t, result.AccessToken)
	require.NotEqual(t, raw, result.RefreshToken)
}

func TestRefresh_Expired(t *testing.T) {
	uc, repo := newAuth(t)
	raw := "expired-token"
	stored := entity.RefreshToken{
		ID:        uuid.New(),
		AdminID:   uuid.New(),
		TokenHash: hashOf(raw),
		ExpiresAt: time.Now().Add(-time.Hour),
	}
	repo.EXPECT().GetRefreshToken(mock.Anything, hashOf(raw)).Return(stored, nil)

	_, err := uc.Refresh(context.Background(), raw)
	require.ErrorIs(t, err, biz.ErrInvalidRefresh)
}

func TestRefresh_Revoked(t *testing.T) {
	uc, repo := newAuth(t)
	raw := "revoked-token"
	at := time.Now().Add(-time.Minute)
	stored := entity.RefreshToken{
		ID:        uuid.New(),
		AdminID:   uuid.New(),
		TokenHash: hashOf(raw),
		ExpiresAt: time.Now().Add(time.Hour),
		RevokedAt: &at,
	}
	repo.EXPECT().GetRefreshToken(mock.Anything, hashOf(raw)).Return(stored, nil)

	_, err := uc.Refresh(context.Background(), raw)
	require.ErrorIs(t, err, biz.ErrInvalidRefresh)
}

// Reusing an already-rotated token is treated as theft: every session for
// that admin is revoked.
func TestRefresh_ReuseRevokesAll(t *testing.T) {
	uc, repo := newAuth(t)
	raw := "reused-token"
	adminID := uuid.New()
	at := time.Now().Add(-time.Minute)
	stored := entity.RefreshToken{
		ID:        uuid.New(),
		AdminID:   adminID,
		TokenHash: hashOf(raw),
		ExpiresAt: time.Now().Add(time.Hour),
		RotatedAt: &at,
	}
	repo.EXPECT().GetRefreshToken(mock.Anything, hashOf(raw)).Return(stored, nil)
	repo.EXPECT().RevokeAllForAdmin(mock.Anything, adminID, mock.Anything).Return(nil)

	_, err := uc.Refresh(context.Background(), raw)
	require.ErrorIs(t, err, biz.ErrInvalidRefresh)
}

func TestRefresh_Unknown(t *testing.T) {
	uc, repo := newAuth(t)
	repo.EXPECT().GetRefreshToken(mock.Anything, mock.Anything).
		Return(entity.RefreshToken{}, biz.ErrResourceNotFound)

	_, err := uc.Refresh(context.Background(), "never-issued")
	require.ErrorIs(t, err, biz.ErrInvalidRefresh)
}

func TestLogout_RevokesToken(t *testing.T) {
	uc, repo := newAuth(t)
	raw := "live-token"
	repo.EXPECT().RevokeRefreshToken(mock.Anything, hashOf(raw), mock.Anything).Return(nil)

	require.NoError(t, uc.Logout(context.Background(), raw))
}

func TestLogout_UnknownIsIdempotent(t *testing.T) {
	uc, repo := newAuth(t)
	repo.EXPECT().RevokeRefreshToken(mock.Anything, mock.Anything, mock.Anything).
		Return(biz.ErrResourceNotFound)

	require.NoError(t, uc.Logout(context.Background(), "unknown"))
}

func TestMe_DisabledAccountRejected(t *testing.T) {
	uc, repo := newAuth(t)
	acc := account(t, "pw12345678", entity.RoleAdmin, false)
	repo.EXPECT().GetAccount(mock.Anything, acc.ID).Return(acc, nil)

	_, err := uc.Me(context.Background(), acc.ID.String())
	require.ErrorIs(t, err, biz.ErrInvalidCredentials)
}
