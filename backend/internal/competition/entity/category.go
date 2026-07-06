package entity

import (
	"time"

	"github.com/google/uuid"
)

// Category is a first-class competition category (own table, FK from
// competitions.category_id).
type Category struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}
