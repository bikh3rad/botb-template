package biz

import "errors"

var (
	// ErrResourceNotFound is returned when a user does not exist.
	ErrResourceNotFound = errors.New("user not found")
	// ErrResourceInvalid is returned for malformed input.
	ErrResourceInvalid = errors.New("invalid user resource")
	// ErrResourceExists is returned on a duplicate email.
	ErrResourceExists = errors.New("user already exists")
	// ErrCompetitionNotFound is returned when purchasing against an unknown competition.
	ErrCompetitionNotFound = errors.New("competition not found")
	// ErrUserSuspended is returned when a suspended (is_active=false) user
	// attempts a purchase.
	ErrUserSuspended = errors.New("user is suspended")
)
