package biz

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"application/internal/media/entity"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Per-kind upload limits and the content-type allow-list. Kept in the use case
// so validation is authoritative and unit-testable independent of the transport.
const (
	maxImageBytes int64 = 10 << 20  // 10 MiB
	maxVideoBytes int64 = 200 << 20 // 200 MiB
)

// allowedTypes maps an accepted content type to its media kind and file
// extension used when composing the object key.
var allowedTypes = map[string]struct {
	kind entity.Kind
	ext  string
}{
	"image/jpeg": {entity.KindImage, ".jpg"},
	"image/png":  {entity.KindImage, ".png"},
	"image/webp": {entity.KindImage, ".webp"},
	"video/mp4":  {entity.KindVideo, ".mp4"},
	"video/webm": {entity.KindVideo, ".webm"},
}

type media struct {
	logger  *slog.Logger
	tracer  trace.Tracer
	repo    Repository
	storage ObjectStorage
}

var _ UsecaseMedia = (*media)(nil)

// NewMedia constructs the media use case.
func NewMedia(logger *slog.Logger, repo Repository, storage ObjectStorage) *media {
	return &media{
		logger:  logger.With("layer", "Media"),
		tracer:  otel.Tracer("MediaUseCase"),
		repo:    repo,
		storage: storage,
	}
}

// Upload validates the input, streams the object into storage, and persists the
// metadata row. Storage write happens before the DB insert so a failed upload
// never leaves an orphaned metadata row.
func (uc *media) Upload(ctx context.Context, in UploadInput) (entity.Media, error) {
	logger := uc.logger.With("method", "Upload")

	ctx, span := uc.tracer.Start(ctx, "Upload")
	defer span.End()

	if in.OwnerType == "" || in.OwnerID == uuid.Nil {
		return entity.Media{}, ErrResourceInvalid
	}

	spec, ok := allowedTypes[in.ContentType]
	if !ok {
		return entity.Media{}, errors.Join(ErrUnsupportedType, fmt.Errorf("content type %q", in.ContentType))
	}

	if err := validateSize(spec.kind, in.Size); err != nil {
		return entity.Media{}, err
	}

	id := uuid.New()
	objectKey := fmt.Sprintf("%s/%s/%s%s", in.OwnerType, in.OwnerID, id, spec.ext)

	if err := uc.storage.Put(ctx, objectKey, in.Reader, in.Size, in.ContentType); err != nil {
		logger.ErrorContext(ctx, "failed to store object", "error", err)

		return entity.Media{}, err
	}

	m := entity.Media{
		ID:          id,
		OwnerType:   in.OwnerType,
		OwnerID:     in.OwnerID,
		Kind:        spec.kind,
		Bucket:      uc.storage.Bucket(),
		ObjectKey:   objectKey,
		ContentType: in.ContentType,
		SizeBytes:   in.Size,
		Position:    in.Position,
	}

	stored, err := uc.repo.Create(ctx, m)
	if err != nil {
		logger.ErrorContext(ctx, "failed to persist media metadata", "error", err)
		// Best-effort cleanup so a DB failure does not leak a blob.
		if rmErr := uc.storage.Remove(ctx, objectKey); rmErr != nil {
			logger.WarnContext(ctx, "failed to clean up orphaned object", "error", rmErr)
		}

		return entity.Media{}, err
	}

	return stored, nil
}

// Get returns a media record together with a presigned read URL.
func (uc *media) Get(ctx context.Context, id uuid.UUID) (MediaWithURL, error) {
	logger := uc.logger.With("method", "Get")

	m, err := uc.repo.Get(ctx, id)
	if err != nil {
		return MediaWithURL{}, err
	}

	url, err := uc.storage.PresignGet(ctx, m.ObjectKey)
	if err != nil {
		logger.ErrorContext(ctx, "failed to presign url", "error", err)

		return MediaWithURL{}, err
	}

	return MediaWithURL{Media: m, URL: url}, nil
}

const (
	defaultPageSize = 50
	maxPageSize     = 200
)

// ListAll returns a paged global media listing (admin media library).
func (uc *media) ListAll(ctx context.Context, limit, offset int) (MediaPage, error) {
	if limit <= 0 {
		limit = defaultPageSize
	}

	if limit > maxPageSize {
		limit = maxPageSize
	}

	if offset < 0 {
		offset = 0
	}

	return uc.repo.ListAll(ctx, limit, offset)
}

// Delete removes the media record AND its MinIO object. Ordering is
// deliberate: the DB row is deleted FIRST, then the object best-effort. A
// failed object removal leaves only an invisible storage leak (logged loudly,
// reconcilable later); the reverse order could leave a dangling DB row
// pointing at a deleted object, which breaks the UI. A missing object is
// treated as success (idempotent delete).
func (uc *media) Delete(ctx context.Context, id uuid.UUID) error {
	logger := uc.logger.With("method", "Delete")

	ctx, span := uc.tracer.Start(ctx, "Delete")
	defer span.End()

	m, err := uc.repo.Get(ctx, id)
	if err != nil {
		return err
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		return err
	}

	if err := uc.storage.Remove(ctx, m.ObjectKey); err != nil {
		logger.ErrorContext(ctx, "DB row deleted but object removal failed — orphaned object in storage",
			"object_key", m.ObjectKey, "error", err)
	}

	return nil
}

// Update reorders (position) and/or reassigns ownership. Owner reassignment
// is a single UPDATE so it is offered alongside reorder.
func (uc *media) Update(ctx context.Context, id uuid.UUID, in UpdateInput) (entity.Media, error) {
	ctx, span := uc.tracer.Start(ctx, "Update")
	defer span.End()

	if in.Position == nil && in.OwnerType == "" && in.OwnerID == nil {
		return entity.Media{}, ErrResourceInvalid
	}

	if in.Position != nil && *in.Position < 0 {
		return entity.Media{}, ErrResourceInvalid
	}

	// Owner reassignment needs both halves of the owner reference.
	if (in.OwnerType != "") != (in.OwnerID != nil) {
		return entity.Media{}, ErrResourceInvalid
	}

	return uc.repo.Update(ctx, id, in)
}

// ListByOwner returns all media for an owner object, ordered by position.
func (uc *media) ListByOwner(
	ctx context.Context,
	ownerType string,
	ownerID uuid.UUID,
) ([]entity.Media, error) {
	if ownerType == "" || ownerID == uuid.Nil {
		return nil, ErrResourceInvalid
	}

	return uc.repo.ListByOwner(ctx, ownerType, ownerID)
}

func validateSize(kind entity.Kind, size int64) error {
	if size <= 0 {
		return ErrResourceInvalid
	}

	limit := maxImageBytes
	if kind == entity.KindVideo {
		limit = maxVideoBytes
	}

	if size > limit {
		return errors.Join(ErrFileTooLarge, fmt.Errorf("size %d exceeds limit %d", size, limit))
	}

	return nil
}
