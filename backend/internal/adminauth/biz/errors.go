package biz

import "errors"

var (
	// ErrInvalidCredentials is returned for ANY failed login (unknown email,
	// wrong password, disabled account) — deliberately generic so responses
	// don't reveal which part failed.
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrRateLimited is returned when login attempts exceed the window cap.
	ErrRateLimited = errors.New("too many login attempts")
	// ErrInvalidRefresh is returned for an unknown/expired/revoked/reused
	// refresh token.
	ErrInvalidRefresh = errors.New("invalid refresh token")
	// ErrResourceNotFound is returned when an admin account does not exist.
	ErrResourceNotFound = errors.New("admin account not found")
	// ErrResourceInvalid is returned for malformed input.
	ErrResourceInvalid = errors.New("invalid adminauth resource")
	// ErrResourceExists is returned on a duplicate admin email.
	ErrResourceExists = errors.New("admin account already exists")
	// ErrLastSuperadmin guards against disabling or demoting the last active
	// superadmin, which would lock everyone out of account management.
	ErrLastSuperadmin = errors.New("cannot disable or demote the last active superadmin")
)
