package middlewares_test

import (
	"application/pkg/middlewares"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

const secret = "unit-test-secret"

func sign(t *testing.T, key string, claims jwt.MapClaims) string {
	t.Helper()

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := tok.SignedString([]byte(key))
	require.NoError(t, err)

	return s
}

func adminClaims(role string) jwt.MapClaims {
	return jwt.MapClaims{
		"sub":   "admin-1",
		"email": "admin@example.com",
		"role":  role,
		"exp":   time.Now().Add(time.Hour).Unix(),
	}
}

// serve runs the given guard around a 200 handler and returns the status.
func serve(guard func(http.Handler) http.Handler, authorization string) int {
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/x", nil)
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}

	rec := httptest.NewRecorder()
	guard(next).ServeHTTP(rec, req)

	return rec.Code
}

func TestJWT_ValidToken(t *testing.T) {
	auth := middlewares.NewJWTAuth(middlewares.JWTSecret(secret))
	token := sign(t, secret, adminClaims(middlewares.RoleAdmin))

	require.Equal(t, http.StatusOK, serve(auth.Middleware, "Bearer "+token))
	require.Equal(t, http.StatusOK, serve(auth.RequireAdmin, "Bearer "+token))
}

func TestJWT_ExpiredToken(t *testing.T) {
	auth := middlewares.NewJWTAuth(middlewares.JWTSecret(secret))
	claims := adminClaims(middlewares.RoleAdmin)
	claims["exp"] = time.Now().Add(-time.Hour).Unix()
	token := sign(t, secret, claims)

	require.Equal(t, http.StatusUnauthorized, serve(auth.RequireAdmin, "Bearer "+token))
}

func TestJWT_MissingToken(t *testing.T) {
	auth := middlewares.NewJWTAuth(middlewares.JWTSecret(secret))

	require.Equal(t, http.StatusUnauthorized, serve(auth.RequireAdmin, ""))
	require.Equal(t, http.StatusUnauthorized, serve(auth.RequireAdmin, "Bearer "))
	require.Equal(t, http.StatusUnauthorized, serve(auth.RequireAdmin, "not-a-bearer"))
}

func TestJWT_WrongSignature(t *testing.T) {
	auth := middlewares.NewJWTAuth(middlewares.JWTSecret(secret))
	token := sign(t, "some-other-secret", adminClaims(middlewares.RoleAdmin))

	require.Equal(t, http.StatusUnauthorized, serve(auth.RequireAdmin, "Bearer "+token))
}

// A validly-signed token WITHOUT an admin role must be rejected with 403 —
// distinct from 401 — on admin route groups. This closes the old hole where
// any signed token was accepted as an admin.
func TestJWT_ValidTokenWrongRole(t *testing.T) {
	auth := middlewares.NewJWTAuth(middlewares.JWTSecret(secret))

	noRole := sign(t, secret, jwt.MapClaims{"sub": "x", "exp": time.Now().Add(time.Hour).Unix()})
	require.Equal(t, http.StatusForbidden, serve(auth.RequireAdmin, "Bearer "+noRole))
	require.Equal(t, http.StatusOK, serve(auth.Middleware, "Bearer "+noRole))

	userRole := sign(t, secret, adminClaims("user"))
	require.Equal(t, http.StatusForbidden, serve(auth.RequireAdmin, "Bearer "+userRole))
}

func TestJWT_SuperadminGuard(t *testing.T) {
	auth := middlewares.NewJWTAuth(middlewares.JWTSecret(secret))

	adminTok := sign(t, secret, adminClaims(middlewares.RoleAdmin))
	require.Equal(t, http.StatusForbidden, serve(auth.RequireSuperadmin, "Bearer "+adminTok))

	superTok := sign(t, secret, adminClaims(middlewares.RoleSuperadmin))
	require.Equal(t, http.StatusOK, serve(auth.RequireSuperadmin, "Bearer "+superTok))
	require.Equal(t, http.StatusOK, serve(auth.RequireAdmin, "Bearer "+superTok))
}

func TestJWT_ClaimsInContext(t *testing.T) {
	auth := middlewares.NewJWTAuth(middlewares.JWTSecret(secret))
	token := sign(t, secret, adminClaims(middlewares.RoleSuperadmin))

	var got middlewares.Claims

	var found bool

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, found = middlewares.ClaimsFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/x", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	auth.RequireAdmin(next).ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.True(t, found)
	require.Equal(t, "admin-1", got.Subject)
	require.Equal(t, "admin@example.com", got.Email)
	require.Equal(t, middlewares.RoleSuperadmin, got.Role)
}
