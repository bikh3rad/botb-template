package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"application/app"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/wire"
)

// JWTSecret is the shared HMAC signing secret. It is defined once in config
// (`jwt.secret`, APP_JWT_SECRET) and read the same way by the gateway and every
// service — it is shared infrastructure config, not per-service state.
type JWTSecret []byte

// Admin roles carried in the `role` JWT claim. Tokens are issued by the
// adminauth service; every other service only verifies.
const (
	RoleAdmin      = "admin"
	RoleSuperadmin = "superadmin"
)

var (
	// ErrMissingToken is returned when no bearer token is present.
	ErrMissingToken = errors.New("missing bearer token")
	// ErrInvalidToken is returned for a malformed/expired/badly-signed token.
	ErrInvalidToken = errors.New("invalid or expired token")
	// ErrNoJWTSecret is returned when jwt.secret is not configured.
	ErrNoJWTSecret = errors.New("jwt.secret is required")
)

// Claims are the verified token claims stashed into the request context so
// handlers can attribute actions (e.g. audit entries) to the acting admin.
type Claims struct {
	Subject string
	Email   string
	Role    string
}

// IsAdmin reports whether the claims carry an admin-capable role.
func (c Claims) IsAdmin() bool {
	return c.Role == RoleAdmin || c.Role == RoleSuperadmin
}

type claimsCtxKey struct{}

// ClaimsFromContext returns the verified claims placed by the JWT middleware.
func ClaimsFromContext(ctx context.Context) (Claims, bool) {
	c, ok := ctx.Value(claimsCtxKey{}).(Claims)

	return c, ok
}

// ContextWithClaims stashes claims; exported for tests that call handlers
// directly without the middleware.
func ContextWithClaims(ctx context.Context, c Claims) context.Context {
	return context.WithValue(ctx, claimsCtxKey{}, c)
}

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

// Middleware guards a handler, requiring a valid HS256 bearer token of any
// role. Prefer RequireAdmin/RequireSuperadmin for admin route groups.
func (a *JWTAuth) Middleware(next http.Handler) http.Handler {
	return a.require(next, nil)
}

// RequireAdmin guards a handler, requiring a valid token whose role claim is
// admin or superadmin. Missing/invalid token → 401; valid token with a wrong
// (or absent) role → 403, so callers can tell authentication failures apart
// from authorization failures.
func (a *JWTAuth) RequireAdmin(next http.Handler) http.Handler {
	return a.require(next, func(c Claims) bool { return c.IsAdmin() })
}

// RequireSuperadmin guards a handler, requiring role=superadmin.
func (a *JWTAuth) RequireSuperadmin(next http.Handler) http.Handler {
	return a.require(next, func(c Claims) bool { return c.Role == RoleSuperadmin })
}

func (a *JWTAuth) require(next http.Handler, allowed func(Claims) bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := a.authenticate(r)
		if err != nil {
			writeUnauthorized(w)

			return
		}

		if allowed != nil && !allowed(claims) {
			writeForbidden(w)

			return
		}

		next.ServeHTTP(w, r.WithContext(ContextWithClaims(r.Context(), claims)))
	})
}

// authenticate parses and validates the Authorization bearer token. Token
// expiry (exp) is enforced by the parser's default validation.
func (a *JWTAuth) authenticate(r *http.Request) (Claims, error) {
	raw, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !ok || raw == "" {
		return Claims{}, ErrMissingToken
	}

	token, err := jwt.Parse(raw, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}

		return a.secret, nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil || !token.Valid {
		return Claims{}, ErrInvalidToken
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return Claims{}, ErrInvalidToken
	}

	claims := Claims{}
	if sub, err := mapClaims.GetSubject(); err == nil {
		claims.Subject = sub
	}

	if email, ok := mapClaims["email"].(string); ok {
		claims.Email = email
	}

	if role, ok := mapClaims["role"].(string); ok {
		claims.Role = role
	}

	return claims, nil
}

func writeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"message":"unauthorized"}`))
}

func writeForbidden(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	_, _ = w.Write([]byte(`{"message":"forbidden"}`))
}
