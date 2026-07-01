package dto

import (
	"application/internal/media/entity"
)

// MediaResp is the API representation of a stored media object.
type MediaResp struct {
	ID              string   `json:"id"`
	OwnerType       string   `json:"owner_type"`
	OwnerID         string   `json:"owner_id"`
	Kind            string   `json:"kind"`
	Bucket          string   `json:"bucket"`
	ObjectKey       string   `json:"object_key"`
	ContentType     string   `json:"content_type"`
	SizeBytes       int64    `json:"size_bytes"`
	Width           int      `json:"width,omitempty"`
	Height          int      `json:"height,omitempty"`
	DurationSeconds *float64 `json:"duration_seconds,omitempty"`
	Position        int      `json:"position"`
	CreatedAt       string   `json:"created_at"`
	// URL is a time-limited presigned link for reading the object. Populated on
	// single-resource reads, empty on create.
	URL string `json:"url,omitempty"`
}

// ToMediaResp maps an entity to its API shape.
func ToMediaResp(m entity.Media) MediaResp {
	return MediaResp{
		ID:              m.ID.String(),
		OwnerType:       m.OwnerType,
		OwnerID:         m.OwnerID.String(),
		Kind:            string(m.Kind),
		Bucket:          m.Bucket,
		ObjectKey:       m.ObjectKey,
		ContentType:     m.ContentType,
		SizeBytes:       m.SizeBytes,
		Width:           m.Width,
		Height:          m.Height,
		DurationSeconds: m.DurationSeconds,
		Position:        m.Position,
		CreatedAt:       m.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
	}
}

// MediaListResp is the list envelope used for owner queries.
type MediaListResp struct {
	Count int         `json:"count"`
	Media []MediaResp `json:"media"`
}

// ToMediaListResp maps a slice of entities to the list envelope.
func ToMediaListResp(ms []entity.Media) MediaListResp {
	items := make([]MediaResp, 0, len(ms))
	for i := range ms {
		items = append(items, ToMediaResp(ms[i]))
	}

	return MediaListResp{Count: len(items), Media: items}
}
