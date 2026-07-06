package biz

import (
	"time"

	"application/app"
)

// Config is the koanf `adminauth` sub-tree. Env overrides follow the APP_
// convention, e.g. APP_ADMINAUTH_BOOTSTRAP_EMAIL → adminauth.bootstrap.email.
type Config struct {
	AccessTTL  time.Duration `koanf:"access_ttl"`  // access JWT lifetime
	RefreshTTL time.Duration `koanf:"refresh_ttl"` // refresh token lifetime
	Bootstrap  Bootstrap     `koanf:"bootstrap"`
}

// Bootstrap describes the optional first-superadmin seed (see EnsureBootstrap).
// Credentials come from env/config only — never hardcoded.
type Bootstrap struct {
	Name     string `koanf:"name"`
	Email    string `koanf:"email"`
	Password string `koanf:"password"`
}

const (
	defaultAccessTTL  = 15 * time.Minute
	defaultRefreshTTL = 168 * time.Hour // 7 days
)

// NewConfig loads the `adminauth` sub-tree, applying safe defaults.
func NewConfig(c *app.KConfig) (*Config, error) {
	cfg := new(Config)
	if err := c.Unmarshal("adminauth", cfg); err != nil {
		return nil, err
	}

	if cfg.AccessTTL <= 0 {
		cfg.AccessTTL = defaultAccessTTL
	}

	if cfg.RefreshTTL <= 0 {
		cfg.RefreshTTL = defaultRefreshTTL
	}

	return cfg, nil
}
