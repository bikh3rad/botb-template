package entity

import (
	"time"

	"github.com/google/uuid"
)

// User is a registered player.
type User struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	TicketsOwned    int64     `json:"tickets_owned"`
	TotalSpentPence int64     `json:"total_spent_pence"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
}

// Ticket is a single entry a user bought into a competition.
type Ticket struct {
	ID            uuid.UUID `json:"id"`
	CompetitionID uuid.UUID `json:"competition_id"`
	UserID        uuid.UUID `json:"user_id"`
	PurchasedAt   time.Time `json:"purchased_at"`
}
