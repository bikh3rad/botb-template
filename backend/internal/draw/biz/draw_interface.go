package biz

import (
	"context"

	"application/internal/draw/entity"

	"github.com/google/uuid"
)

// ListFilter narrows/paginates a draw listing (same shape as the user service).
type ListFilter struct {
	Query  string // matched against prize (case-insensitive), empty = all
	Limit  int
	Offset int
}

// DrawPage is a page of draws plus the total match count.
type DrawPage struct {
	Draws []entity.Draw
	Total int
}

// CreateInput is the biz-level payload for creating a pending draw.
type CreateInput struct {
	CompetitionID uuid.UUID
	Prize         string
}

// UsecaseDraw is consumed by the HTTP handler.
type UsecaseDraw interface {
	List(ctx context.Context, filter ListFilter) (DrawPage, error)
	Get(ctx context.Context, id uuid.UUID) (entity.Draw, error)
	// GetPublic returns a draw only if it has been drawn, hiding pending AND
	// void draws from the public homepage (returns ErrResourceNotFound
	// otherwise) — a voided result must not show a winner.
	GetPublic(ctx context.Context, id uuid.UUID) (entity.Draw, error)
	Create(ctx context.Context, in CreateInput) (entity.Draw, error)
	Run(ctx context.Context, id uuid.UUID) (entity.Draw, error)
	// UpdatePrize edits the prize text only (any status except void).
	UpdatePrize(ctx context.Context, id uuid.UUID, prize string) (entity.Draw, error)
	// Void marks a pending or drawn draw void with a REQUIRED reason. The safe
	// path to change a winner is void + create a new draw + run it.
	Void(ctx context.Context, id uuid.UUID, reason string) (entity.Draw, error)
	// Reassign directly moves a DRAWN draw's winner to another ticket of the
	// same competition. Requires a reason; exists so the mutation is explicit
	// and audited rather than a hand-edited row.
	Reassign(ctx context.Context, id uuid.UUID, ticketID uuid.UUID, reason string) (entity.Draw, error)
}

// Repository persists draws and runs the winner-selection transaction.
// Implemented by internal/draw/repo (pgx).
type Repository interface {
	List(ctx context.Context, filter ListFilter) (DrawPage, error)
	Get(ctx context.Context, id uuid.UUID) (entity.Draw, error)
	Create(ctx context.Context, d entity.Draw) (entity.Draw, error)
	// Run atomically picks a winning ticket for the draw's competition and marks
	// the draw drawn. It must reject a non-pending draw (ErrAlreadyDrawn) and a
	// competition with no tickets (ErrNoTickets).
	Run(ctx context.Context, id uuid.UUID) (entity.Draw, error)
	UpdatePrize(ctx context.Context, id uuid.UUID, prize string) (entity.Draw, error)
	// Void sets status=void + void_reason atomically, only from pending/drawn.
	Void(ctx context.Context, id uuid.UUID, reason string) (entity.Draw, error)
	// Reassign validates (inside a transaction) that the ticket belongs to the
	// draw's competition and swaps winner_ticket_id/winner_user_id.
	Reassign(ctx context.Context, id uuid.UUID, ticketID uuid.UUID) (entity.Draw, error)
}
