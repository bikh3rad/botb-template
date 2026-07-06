package handler_test

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"application/internal/adminauth/biz"
	"application/internal/adminauth/entity"
	adminhandler "application/internal/adminauth/handler"
	"application/internal/adminauth/mocks"
	"application/pkg/audit"
	"application/pkg/middlewares"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const testSecret = "test-secret"

func tokenWithRole(role string) string {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   uuid.NewString(),
		"email": "admin@example.com",
		"role":  role,
		"exp":   time.Now().Add(time.Hour).Unix(),
	})
	s, _ := tok.SignedString([]byte(testSecret))

	return s
}

func newTestHandlers(t *testing.T) (*http.ServeMux, *mocks.MockUsecaseAuth, *mocks.MockUsecaseAccounts) {
	t.Helper()

	ucAuth := mocks.NewMockUsecaseAuth(t)
	ucAccounts := mocks.NewMockUsecaseAccounts(t)
	mux := http.NewServeMux()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	jwtAuth := middlewares.NewJWTAuth(middlewares.JWTSecret(testSecret))
	// nil DB recorder is a no-op — handler audit calls must not explode.
	recorder := audit.NewRecorder(logger, nil)

	a := adminhandler.NewAuth(logger, mux, ucAuth, jwtAuth)
	require.NoError(t, a.RegisterHandler(context.Background()))

	acc := adminhandler.NewAccounts(logger, mux, ucAccounts, jwtAuth, recorder)
	require.NoError(t, acc.RegisterHandler(context.Background()))

	return mux, ucAuth, ucAccounts
}

func do(mux *http.ServeMux, method, target, token, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequestWithContext(context.Background(), method, target, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	return rec
}

func TestLogin_OK(t *testing.T) {
	mux, ucAuth, _ := newTestHandlers(t)
	adminID := uuid.New()

	ucAuth.EXPECT().Login(mock.Anything, "a@b.co", "pw", mock.Anything).
		Return(biz.LoginResult{
			AccessToken:  "access",
			ExpiresIn:    900,
			RefreshToken: "refresh",
			Admin:        entity.AdminAccount{ID: adminID, Role: entity.RoleAdmin, Email: "a@b.co"},
		}, nil)

	rec := do(mux, http.MethodPost, "/apis/adminauth/v1/login", "", `{"email":"a@b.co","password":"pw"}`)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]any

	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	require.Equal(t, "access", resp["access_token"])
	require.Equal(t, "Bearer", resp["token_type"])
	require.Equal(t, "refresh", resp["refresh_token"])
}

func TestLogin_GenericFailure(t *testing.T) {
	mux, ucAuth, _ := newTestHandlers(t)
	ucAuth.EXPECT().Login(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(biz.LoginResult{}, biz.ErrInvalidCredentials)

	rec := do(mux, http.MethodPost, "/apis/adminauth/v1/login", "", `{"email":"a@b.co","password":"bad"}`)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
	// The body must not leak WHICH part failed.
	require.NotContains(t, rec.Body.String(), "password")
	require.NotContains(t, rec.Body.String(), "email")
}

func TestLogin_RateLimited429(t *testing.T) {
	mux, ucAuth, _ := newTestHandlers(t)
	ucAuth.EXPECT().Login(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(biz.LoginResult{}, biz.ErrRateLimited)

	rec := do(mux, http.MethodPost, "/apis/adminauth/v1/login", "", `{"email":"a@b.co","password":"x"}`)
	require.Equal(t, http.StatusTooManyRequests, rec.Code)
}

func TestRefresh_Invalid401(t *testing.T) {
	mux, ucAuth, _ := newTestHandlers(t)
	ucAuth.EXPECT().Refresh(mock.Anything, "stale").Return(biz.LoginResult{}, biz.ErrInvalidRefresh)

	rec := do(mux, http.MethodPost, "/apis/adminauth/v1/refresh", "", `{"refresh_token":"stale"}`)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestLogout_NoContent(t *testing.T) {
	mux, ucAuth, _ := newTestHandlers(t)
	ucAuth.EXPECT().Logout(mock.Anything, "live").Return(nil)

	rec := do(mux, http.MethodPost, "/apis/adminauth/v1/logout", "", `{"refresh_token":"live"}`)
	require.Equal(t, http.StatusNoContent, rec.Code)
}

func TestMe_RequiresToken(t *testing.T) {
	mux, _, _ := newTestHandlers(t)

	rec := do(mux, http.MethodGet, "/apis/adminauth/v1/me", "", "")
	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestMe_OK(t *testing.T) {
	mux, ucAuth, _ := newTestHandlers(t)
	ucAuth.EXPECT().Me(mock.Anything, mock.Anything).
		Return(entity.AdminAccount{ID: uuid.New(), Name: "Ops", Email: "a@b.co", Role: entity.RoleAdmin, IsActive: true}, nil)

	rec := do(mux, http.MethodGet, "/apis/adminauth/v1/me", tokenWithRole("admin"), "")
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), `"role":"admin"`)
	require.NotContains(t, rec.Body.String(), "password")
}

// Account management is superadmin-only IN the service (not just the
// gateway): no token → 401, plain admin → 403, superadmin → 200.
func TestAccounts_RoleLadder(t *testing.T) {
	mux, _, ucAccounts := newTestHandlers(t)

	rec := do(mux, http.MethodGet, "/apis/adminauth/v1/admin/accounts", "", "")
	require.Equal(t, http.StatusUnauthorized, rec.Code)

	rec = do(mux, http.MethodGet, "/apis/adminauth/v1/admin/accounts", tokenWithRole("admin"), "")
	require.Equal(t, http.StatusForbidden, rec.Code)

	ucAccounts.EXPECT().List(mock.Anything).Return([]entity.AdminAccount{}, nil)
	rec = do(mux, http.MethodGet, "/apis/adminauth/v1/admin/accounts", tokenWithRole("superadmin"), "")
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestAccounts_Create(t *testing.T) {
	mux, _, ucAccounts := newTestHandlers(t)

	ucAccounts.EXPECT().Create(mock.Anything, mock.Anything).
		Return(entity.AdminAccount{ID: uuid.New(), Name: "Ops", Email: "ops@b.co", Role: entity.RoleAdmin, IsActive: true}, nil)

	body := `{"name":"Ops","email":"ops@b.co","password":"long-enough","role":"admin"}`
	rec := do(mux, http.MethodPost, "/apis/adminauth/v1/admin/accounts", tokenWithRole("superadmin"), body)
	require.Equal(t, http.StatusCreated, rec.Code)
}

func TestAccounts_UpdateLastSuperadminConflict(t *testing.T) {
	mux, _, ucAccounts := newTestHandlers(t)
	id := uuid.New()

	ucAccounts.EXPECT().Update(mock.Anything, id, mock.Anything).
		Return(entity.AdminAccount{}, biz.ErrLastSuperadmin)

	rec := do(mux, http.MethodPut, "/apis/adminauth/v1/admin/accounts/"+id.String(),
		tokenWithRole("superadmin"), `{"is_active":false}`)
	require.Equal(t, http.StatusConflict, rec.Code)
}
