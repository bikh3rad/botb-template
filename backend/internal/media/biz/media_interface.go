package biz

import (
	"context"
	"io"

	"application/internal/media/entity"

	"github.com/google/uuid"
)

// UploadInput is the biz-level payload for creating a media object. The handler
// parses the multipart request and hands the file stream + metadata here.
type UploadInput struct {
	OwnerType   string
	OwnerID     uuid.UUID
	ContentType string
	Size        int64
	Position    int
	Reader      io.Reader
}

// MediaWithURL pairs a stored media record with a freshly minted presigned URL.
type MediaWithURL struct {
	Media entity.Media
	URL   string
}

// UpdateInput is a partial media edit: reorder (position) and/or owner
// reassignment. Nil/empty fields are unchanged.
type UpdateInput struct {
	Position  *int
	OwnerType string
	OwnerID   *uuid.UUID
}

// MediaPage is a page of media rows plus the total count (media-library view).
type MediaPage struct {
	Items []entity.Media
	Total int
}

// UsecaseMedia is consumed by the HTTP handler. Replace has no endpoint by
// design: replacing = upload new + delete old.
type UsecaseMedia interface {
	Upload(ctx context.Context, in UploadInput) (entity.Media, error)
	Get(ctx context.Context, id uuid.UUID) (MediaWithURL, error)
	ListByOwner(ctx context.Context, ownerType string, ownerID uuid.UUID) ([]entity.Media, error)
	// ListAll returns a paged global listing for the admin media library.
	ListAll(ctx context.Context, limit, offset int) (MediaPage, error)
	// Delete removes the DB row and the MinIO object (see biz.Delete for the
	// deliberate ordering).
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, in UpdateInput) (entity.Media, error)
}

// Repository persists media metadata; implemented by internal/media/repo (pgx).
type Repository interface {
	Create(ctx context.Context, m entity.Media) (entity.Media, error)
	Get(ctx context.Context, id uuid.UUID) (entity.Media, error)
	ListByOwner(ctx context.Context, ownerType string, ownerID uuid.UUID) ([]entity.Media, error)
	ListAll(ctx context.Context, limit, offset int) (MediaPage, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, in UpdateInput) (entity.Media, error)
}

// ObjectStorage is the S3-compatible blob store, owned here so MinIO can be
// swapped for AWS S3 without touching the use case. Implemented by
// datasource.MinioStorage.
type ObjectStorage interface {
	Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error
	PresignGet(ctx context.Context, key string) (string, error)
	Remove(ctx context.Context, key string) error
	Bucket() string
}
