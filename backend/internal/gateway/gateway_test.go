//nolint:testpackage // needs the unexported gatewayConfig (template keeps config types unexported)
package gateway

import (
	"application/pkg/middlewares"
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

const testSecret = "gateway-test-secret"

func validToken(t *testing.T) string {
	t.Helper()

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "admin",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	s, err := tok.SignedString([]byte(testSecret))
	require.NoError(t, err)

	return s
}

// upstream returns a test server that echoes a marker so the caller can assert
// which upstream received the request.
func upstream(t *testing.T, marker string) *httptest.Server {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, marker+":"+r.URL.Path)
	}))
	t.Cleanup(srv.Close)

	return srv
}

func newGateway(t *testing.T, upstreams map[string]string) *http.ServeMux {
	t.Helper()

	mux := http.NewServeMux()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	auth := middlewares.NewJWTAuth(middlewares.JWTSecret(testSecret))

	g, err := NewGateway(logger, mux, &gatewayConfig{Upstreams: upstreams}, auth)
	require.NoError(t, err)
	require.NoError(t, g.RegisterHandler(context.Background()))

	return mux
}

func request(mux *http.ServeMux, method, target, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequestWithContext(context.Background(), method, target, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	return rec
}

func TestDispatch_PublicNoTokenProxied(t *testing.T) {
	comp := upstream(t, "competition")
	mux := newGateway(t, map[string]string{"competition": comp.URL})

	rec := request(mux, http.MethodGet, "/apis/competition/v1/competitions", "")
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "competition:/apis/competition/v1/competitions")
}

func TestDispatch_AdminRequiresToken(t *testing.T) {
	comp := upstream(t, "competition")
	mux := newGateway(t, map[string]string{"competition": comp.URL})

	rec := request(mux, http.MethodPost, "/apis/competition/v1/admin/competitions", "")
	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestDispatch_AdminWithTokenProxied(t *testing.T) {
	comp := upstream(t, "competition")
	mux := newGateway(t, map[string]string{"competition": comp.URL})

	rec := request(mux, http.MethodPost, "/apis/competition/v1/admin/competitions", validToken(t))
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "competition:")
}

func TestDispatch_UnknownServiceIs404(t *testing.T) {
	comp := upstream(t, "competition")
	mux := newGateway(t, map[string]string{"competition": comp.URL})

	rec := request(mux, http.MethodGet, "/apis/unknown/v1/things", "")
	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDispatch_SelectsCorrectUpstream(t *testing.T) {
	comp := upstream(t, "competition")
	usr := upstream(t, "user")
	mux := newGateway(t, map[string]string{"competition": comp.URL, "user": usr.URL})

	rec := request(mux, http.MethodGet, "/apis/user/v1/users", "")
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "user:/apis/user/v1/users")
}
