package biz

import (
	"application/internal/user/entity"
	"context"

	"github.com/google/uuid"
)

// PurchaseInput is the biz-level payload for buying tickets.
type PurchaseInput struct {
	CompetitionID uuid.UUID
	UserID        uuid.UUID
	Quantity      int
}

// PurchaseResult is returned after a successful purchase: the created tickets,
// the refreshed user, and the total charged.
type PurchaseResult struct {
	Tickets        []entity.Ticket
	User           entity.User
	TotalCostPence int64
}

// UsecaseTicket is consumed by the HTTP handler.
type UsecaseTicket interface {
	Purchase(ctx context.Context, in PurchaseInput) (PurchaseResult, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]entity.Ticket, error)
}

// RepositoryTicket persists tickets and executes the purchase transaction.
// Implemented by internal/user/repo (pgx).
type RepositoryTicket interface {
	// Purchase atomically reads the competition price, inserts `quantity` ticket
	// rows, and increments the user's tickets_owned + total_spent.
	Purchase(ctx context.Context, in PurchaseInput) (PurchaseResult, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]entity.Ticket, error)
}
