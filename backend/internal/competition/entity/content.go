package entity

import "time"

// SiteContent is one editable site-copy value (key-value store; values are
// plain strings, JSON-encoded when structured). Public read, admin write.
type SiteContent struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}
