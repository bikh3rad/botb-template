package entity

import (
	"time"

	"github.com/google/uuid"
)

// Kind is the media type discriminator.
type Kind string

const (
	KindImage Kind = "image"
	KindVideo Kind = "video"
)

// Media is a single stored object (image or video) belonging to some owner
// object (e.g. a competition). An owner can have zero, one, or many Media.
type Media struct {
	ID              uuid.UUID `json:"id"`
	OwnerType       string    `json:"owner_type"`
	OwnerID         uuid.UUID `json:"owner_id"`
	Kind            Kind      `json:"kind"`
	Bucket          string    `json:"bucket"`
	ObjectKey       string    `json:"object_key"`
	ContentType     string    `json:"content_type"`
	SizeBytes       int64     `json:"size_bytes"`
	Width           int       `json:"width"`
	Height          int       `json:"height"`
	DurationSeconds *float64  `json:"duration_seconds,omitempty"`
	Position        int       `json:"position"`
	CreatedAt       time.Time `json:"created_at"`
}
