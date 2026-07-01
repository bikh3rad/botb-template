package biz

import "errors"

var (
	// ErrResourceNotFound is returned when a competition does not exist.
	ErrResourceNotFound = errors.New("competition not found")
	// ErrResourceInvalid is returned for malformed input.
	ErrResourceInvalid = errors.New("invalid competition resource")
	// ErrResourceExists is returned on a unique-constraint conflict (e.g. slug).
	ErrResourceExists = errors.New("competition already exists")
)
