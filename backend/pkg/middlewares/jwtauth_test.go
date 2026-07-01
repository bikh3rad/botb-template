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

// serve runs the auth middleware around a 200 handler and returns the status.
func serve(auth *middlewares.JWTAuth, authorization string) int {
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/x", nil)
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}

	rec := httptest.NewRecorder()
	auth.Middleware(next).ServeHTTP(rec, req)

	return rec.Code
}

func TestJWT_ValidToken(t *testing.T) {
	auth := middlewares.NewJWTAuth(middlewares.JWTSecret(secret))
	token := sign(t, secret, jwt.MapClaims{"sub": "admin", "exp": time.Now().Add(time.Hour).Unix()})

	require.Equal(t, http.StatusOK, serve(auth, "Bearer "+token))
}

func TestJWT_ExpiredToken(t *testing.T) {
	auth := middlewares.NewJWTAuth(middlewares.JWTSecret(secret))
	token := sign(t, secret, jwt.MapClaims{"sub": "admin", "exp": time.Now().Add(-time.Hour).Unix()})

	require.Equal(t, http.StatusUnauthorized, serve(auth, "Bearer "+token))
}

func TestJWT_MissingToken(t *testing.T) {
	auth := middlewares.NewJWTAuth(middlewares.JWTSecret(secret))

	require.Equal(t, http.StatusUnauthorized, serve(auth, ""))
	require.Equal(t, http.StatusUnauthorized, serve(auth, "Bearer "))
	require.Equal(t, http.StatusUnauthorized, serve(auth, "not-a-bearer"))
}

func TestJWT_WrongSignature(t *testing.T) {
	auth := middlewares.NewJWTAuth(middlewares.JWTSecret(secret))
	token := sign(t, "some-other-secret", jwt.MapClaims{"sub": "admin", "exp": time.Now().Add(time.Hour).Unix()})

	require.Equal(t, http.StatusUnauthorized, serve(auth, "Bearer "+token))
}
