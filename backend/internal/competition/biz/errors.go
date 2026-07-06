package biz

import "errors"

var (
	// ErrResourceNotFound is returned when a competition does not exist.
	ErrResourceNotFound = errors.New("competition not found")
	// ErrResourceInvalid is returned for malformed input.
	ErrResourceInvalid = errors.New("invalid competition resource")
	// ErrResourceExists is returned on a unique-constraint conflict (e.g. slug).
	ErrResourceExists = errors.New("competition already exists")
	// ErrInvalidTransition is returned for a disallowed status change
	// (allowed: draft->live, live->closed; closed never reopens).
	ErrInvalidTransition = errors.New("invalid status transition")
	// ErrCategoryInUse is returned when deleting a category still referenced by
	// competitions without a reassignment target.
	ErrCategoryInUse = errors.New("category in use")
	// ErrCategoryNotFound is returned when a referenced category does not exist.
	ErrCategoryNotFound = errors.New("category not found")
	// ErrCompetitionHasEntrants is returned when deleting a competition that has
	// sold tickets or an existing draw — it must be closed or voided instead.
	ErrCompetitionHasEntrants = errors.New("competition has sold tickets or draws; close or void instead")
)
