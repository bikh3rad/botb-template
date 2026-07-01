package entity

import (
	"time"

	"github.com/google/uuid"
)

// Status is the lifecycle state of a competition.
type Status string

const (
	StatusDraft  Status = "draft"
	StatusLive   Status = "live"
	StatusClosed Status = "closed"
)

// Valid reports whether s is a known status.
func (s Status) Valid() bool {
	switch s {
	case StatusDraft, StatusLive, StatusClosed:
		return true
	default:
		return false
	}
}

// MediaRef is a read-only projection of a media object owned by the media
// service. The competition service reads these rows directly from the shared
// `media` table (owner_type='competition', owner_id=competition.id) — see
// repo/competition.go for the rationale of that vs. an HTTP call to media.
type MediaRef struct {
	ID          uuid.UUID `json:"id"`
	Kind        string    `json:"kind"`
	Bucket      string    `json:"bucket"`
	ObjectKey   string    `json:"object_key"`
	ContentType string    `json:"content_type"`
	Position    int       `json:"position"`
}

// Competition is a prize-draw competition. A competition can have zero, one, or
// many associated Media items (image and/or video), resolved into Media.
type Competition struct {
	ID               uuid.UUID  `json:"id"`
	Title            string     `json:"title"`
	Slug             string     `json:"slug"`
	Description      string     `json:"description"`
	Prize            string     `json:"prize"`
	TicketPricePence int64      `json:"ticket_price_pence"`
	TicketsTotal     int64      `json:"tickets_total"`
	TicketsSold      int64      `json:"tickets_sold"`
	Status           Status     `json:"status"`
	StartsAt         time.Time  `json:"starts_at"`
	EndsAt           time.Time  `json:"ends_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	Media            []MediaRef `json:"media"`
}
