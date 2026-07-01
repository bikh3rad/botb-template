package gateway

import (
	"application/app"
	"application/internal/service"
	"application/pkg/middlewares"
	"context"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// gatewayConfig is the koanf sub-tree `gateway`. `upstreams` maps a service name
// (the <servicename> path segment) to that service's base URL. Env overrides
// follow the APP_ convention, e.g. APP_GATEWAY_UPSTREAMS_COMPETITION.
type gatewayConfig struct {
	Upstreams map[string]string `koanf:"upstreams"`
}

// NewGatewayConfig loads the `gateway` sub-tree.
func NewGatewayConfig(_ context.Context, c *app.KConfig) (*gatewayConfig, error) {
	cfg := new(gatewayConfig)
	if err := c.Unmarshal("gateway", cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

type gateway struct {
	logger  *slog.Logger
	mux     *http.ServeMux
	auth    *middlewares.JWTAuth
	proxies map[string]*httputil.ReverseProxy
}

var _ service.Handler = (*gateway)(nil)

// NewGateway builds one reverse proxy per configured upstream. Trace context is
// propagated to upstreams via an otel-instrumented transport.
func NewGateway(
	logger *slog.Logger,
	mux *http.ServeMux,
	cfg *gatewayConfig,
	auth *middlewares.JWTAuth,
) (*gateway, error) {
	proxies := make(map[string]*httputil.ReverseProxy, len(cfg.Upstreams))

	for name, raw := range cfg.Upstreams {
		target, err := url.Parse(raw)
		if err != nil {
			return nil, err
		}

		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.Transport = otelhttp.NewTransport(http.DefaultTransport)
		proxies[name] = proxy
	}

	return &gateway{
		logger:  logger.With("layer", "Gateway"),
		mux:     mux,
		auth:    auth,
		proxies: proxies,
	}, nil
}

// RegisterHandler mounts the single /apis/ dispatch entrypoint.
func (h *gateway) RegisterHandler(_ context.Context) error {
	h.mux.HandleFunc("/apis/", h.dispatch)

	return nil
}

// dispatch routes by the <servicename> path segment. Admin paths
// (/apis/<svc>/v1/admin/...) are JWT-guarded here (first of two layers — each
// service re-checks on its own port); everything else passes through.
func (h *gateway) dispatch(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With("method", "dispatch", "path", r.URL.Path)

	svc, ok := serviceSegment(r.URL.Path)
	if !ok {
		writeJSONError(w, http.StatusBadRequest, "malformed api path")

		return
	}

	proxy, ok := h.proxies[svc]
	if !ok {
		logger.WarnContext(r.Context(), "no upstream for service", "service", svc)
		writeJSONError(w, http.StatusNotFound, "unknown service")

		return
	}

	if isAdminPath(r.URL.Path) {
		h.auth.Middleware(proxy).ServeHTTP(w, r)

		return
	}

	proxy.ServeHTTP(w, r)
}

// serviceSegment extracts the <servicename> from /apis/<servicename>/...
func serviceSegment(path string) (string, bool) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 || parts[0] != "apis" || parts[1] == "" {
		return "", false
	}

	return parts[1], true
}

// isAdminPath reports whether the path is in an admin route group:
// /apis/<svc>/<version>/admin/...
func isAdminPath(path string) bool {
	parts := strings.Split(strings.Trim(path, "/"), "/")

	return len(parts) >= 4 && parts[0] == "apis" && parts[3] == "admin"
}

func writeJSONError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write([]byte(`{"message":"` + message + `"}`))
}
