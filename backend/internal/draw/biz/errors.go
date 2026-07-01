package biz

import "errors"

var (
	// ErrResourceNotFound is returned when a draw does not exist.
	ErrResourceNotFound = errors.New("draw not found")
	// ErrResourceInvalid is returned for malformed input.
	ErrResourceInvalid = errors.New("invalid draw resource")
	// ErrAlreadyDrawn is returned when running a draw that is not pending.
	ErrAlreadyDrawn = errors.New("draw already run")
	// ErrNoTickets is returned when a competition has no tickets to draw from.
	ErrNoTickets = errors.New("competition has no tickets")
)
