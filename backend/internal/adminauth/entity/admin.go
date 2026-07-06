// Package entity holds the adminauth domain types. Admin accounts are a fully
// separate identity from site users (schema adminauth.*, no FKs to public.*).
package entity

import (
	"time"

	"github.com/google/uuid"
)

// Role is an admin account role, carried into the JWT `role` claim.
type Role string

// Admin roles. Superadmins additionally manage admin accounts themselves.
const (
	RoleAdmin      Role = "admin"
	RoleSuperadmin Role = "superadmin"
)

// Valid reports whether the role is one of the known admin roles.
func (r Role) Valid() bool {
	return r == RoleAdmin || r == RoleSuperadmin
}

// AdminAccount is an administrator identity. PasswordHash never leaves the
// service (excluded from DTOs).
type AdminAccount struct {
	ID           uuid.UUID
	Name         string
	Email        string
	PasswordHash string
	Role         Role
	IsActive     bool
	CreatedAt    time.Time
	LastLoginAt  *time.Time
}

// RefreshToken is a stored (hashed) refresh token. Rotation: using a token
// sets RotatedAt and issues a replacement; a token that is expired, revoked,
// or already rotated is rejected.
type RefreshToken struct {
	ID        uuid.UUID
	AdminID   uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	RotatedAt *time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}
