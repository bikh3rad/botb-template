package middlewares

import (
	"application/app"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/wire"
)

// JWTSecret is the shared HMAC signing secret. It is defined once in config
// (`jwt.secret`, APP_JWT_SECRET) and read the same way by the gateway and every
// service — it is shared infrastructure config, not per-service state.
type JWTSecret []byte

var (
	// ErrMissingToken is returned when no bearer token is present.
	ErrMissingToken = errors.New("missing bearer token")
	// ErrInvalidToken is returned for a malformed/expired/badly-signed token.
	ErrInvalidToken = errors.New("invalid or expired token")
	// ErrNoJWTSecret is returned when jwt.secret is not configured.
	ErrNoJWTSecret = errors.New("jwt.secret is required")
)

// NewJWTSecret loads the shared secret from config. Fails fast if unset so a
// misconfigured deployment cannot silently run with unprotected admin routes.
func NewJWTSecret(c *app.KConfig) (JWTSecret, error) {
	secret := c.String("jwt.secret")
	if secret == "" {
		return nil, ErrNoJWTSecret
	}

	return JWTSecret(secret), nil
}

// JWTAuth validates HS256 bearer tokens.
type JWTAuth struct {
	secret []byte
}

// NewJWTAuth constructs the middleware from the shared secret.
func NewJWTAuth(secret JWTSecret) *JWTAuth {
	return &JWTAuth{secret: []byte(secret)}
}

// JWTProviderSet wires the shared secret + auth middleware. Included by the
// gateway and by each service so all use the one configured secret.
var JWTProviderSet = wire.NewSet(NewJWTSecret, NewJWTAuth)

// Middleware guards a handler, requiring a valid HS256 bearer token. It is a
// middlewares.Middleware so it composes with MultipleMiddleware.
func (a *JWTAuth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := a.authenticate(r); err != nil {
			writeUnauthorized(w)

			return
		}

		next.ServeHTTP(w, r)
	})
}

// authenticate parses and validates the Authorization bearer token. Token
// expiry (exp) is enforced by the parser's default validation.
func (a *JWTAuth) authenticate(r *http.Request) error {
	raw, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !ok || raw == "" {
		return ErrMissingToken
	}

	token, err := jwt.Parse(raw, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}

		return a.secret, nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil || !token.Valid {
		return ErrInvalidToken
	}

	return nil
}

func writeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"message":"unauthorized"}`))
}
