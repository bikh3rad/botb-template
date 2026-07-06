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
	// ErrReasonRequired is returned when a void/reassign is attempted without
	// a reason — these are sensitive mutations and must be explainable.
	ErrReasonRequired = errors.New("a reason is required")
	// ErrInvalidState is returned when the draw's status does not allow the
	// requested mutation (e.g. voiding a void draw, reassigning a pending one).
	ErrInvalidState = errors.New("draw status does not allow this operation")
	// ErrTicketMismatch is returned when reassigning to a ticket that does not
	// belong to the draw's competition.
	ErrTicketMismatch = errors.New("ticket does not belong to this competition")
)
