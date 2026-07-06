package entity

import (
	"time"

	"github.com/google/uuid"
)

// Status is the lifecycle state of a draw.
type Status string

const (
	StatusPending Status = "pending"
	StatusDrawn   Status = "drawn"
	StatusVoid    Status = "void"
)

// Valid reports whether s is a known status.
func (s Status) Valid() bool {
	switch s {
	case StatusPending, StatusDrawn, StatusVoid:
		return true
	default:
		return false
	}
}

// Draw is a prize draw for a competition. Winner fields are nil until the draw
// is run; drawnAt records when it was run.
type Draw struct {
	ID             uuid.UUID  `json:"id"`
	CompetitionID  uuid.UUID  `json:"competition_id"`
	WinnerUserID   *uuid.UUID `json:"winner_user_id,omitempty"`
	WinnerTicketID *uuid.UUID `json:"winner_ticket_id,omitempty"`
	Prize          string     `json:"prize"`
	Status         Status     `json:"status"`
	VoidReason     string     `json:"void_reason,omitempty"`
	DrawnAt        *time.Time `json:"drawn_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// WinnerItem is one row of the PUBLIC winners feed: a drawn draw joined to
// its winner's display name. Exists so the public site never needs an admin
// token (the old frontend minted one to read the admin draw+user lists).
type WinnerItem struct {
	DrawID       uuid.UUID  `json:"draw_id"`
	Prize        string     `json:"prize"`
	DrawnAt      *time.Time `json:"drawn_at,omitempty"`
	WinnerUserID uuid.UUID  `json:"winner_user_id"`
	WinnerName   string     `json:"winner_name"`
}
