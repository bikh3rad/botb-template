package biz_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"strings"
	"testing"

	"application/internal/media/biz"
	"application/internal/media/entity"
	"application/internal/media/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestUpload_Success(t *testing.T) {
	ctx := context.Background()
	ownerID := uuid.New()

	repo := mocks.NewMockRepository(t)
	storage := mocks.NewMockObjectStorage(t)

	storage.EXPECT().
		Put(mock.Anything, mock.Anything, mock.Anything, int64(1024), "image/png").
		Return(nil)
	storage.EXPECT().Bucket().Return("botb-media")
	repo.EXPECT().
		Create(mock.Anything, mock.Anything).
		RunAndReturn(func(_ context.Context, m entity.Media) (entity.Media, error) {
			return m, nil
		})

	uc := biz.NewMedia(discardLogger(), repo, storage)

	got, err := uc.Upload(ctx, biz.UploadInput{
		OwnerType:   "competition",
		OwnerID:     ownerID,
		ContentType: "image/png",
		Size:        1024,
		Reader:      strings.NewReader("payload"),
	})

	require.NoError(t, err)
	require.Equal(t, entity.KindImage, got.Kind)
	require.Equal(t, "botb-media", got.Bucket)
	require.Equal(t, ownerID, got.OwnerID)
	require.Contains(t, got.ObjectKey, "competition/"+ownerID.String())
	require.True(t, strings.HasSuffix(got.ObjectKey, ".png"))
}

func TestUpload_UnsupportedType(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	storage := mocks.NewMockObjectStorage(t)

	uc := biz.NewMedia(discardLogger(), repo, storage)

	_, err := uc.Upload(context.Background(), biz.UploadInput{
		OwnerType:   "competition",
		OwnerID:     uuid.New(),
		ContentType: "application/zip",
		Size:        1024,
		Reader:      strings.NewReader("x"),
	})

	require.ErrorIs(t, err, biz.ErrUnsupportedType)
}

func TestUpload_FileTooLarge(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	storage := mocks.NewMockObjectStorage(t)

	uc := biz.NewMedia(discardLogger(), repo, storage)

	_, err := uc.Upload(context.Background(), biz.UploadInput{
		OwnerType:   "competition",
		OwnerID:     uuid.New(),
		ContentType: "image/png",
		Size:        11 << 20, // 11 MiB > 10 MiB image limit
		Reader:      strings.NewReader("x"),
	})

	require.ErrorIs(t, err, biz.ErrFileTooLarge)
}

func TestUpload_InvalidOwner(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	storage := mocks.NewMockObjectStorage(t)

	uc := biz.NewMedia(discardLogger(), repo, storage)

	_, err := uc.Upload(context.Background(), biz.UploadInput{
		OwnerType:   "",
		OwnerID:     uuid.Nil,
		ContentType: "image/png",
		Size:        1024,
		Reader:      strings.NewReader("x"),
	})

	require.ErrorIs(t, err, biz.ErrResourceInvalid)
}

func TestUpload_StorageFailureCleansUp(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	storage := mocks.NewMockObjectStorage(t)

	storage.EXPECT().Put(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("boom"))

	uc := biz.NewMedia(discardLogger(), repo, storage)

	_, err := uc.Upload(context.Background(), biz.UploadInput{
		OwnerType:   "competition",
		OwnerID:     uuid.New(),
		ContentType: "video/mp4",
		Size:        2048,
		Reader:      strings.NewReader("x"),
	})

	require.Error(t, err)
}

func TestGet_Success(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	stored := entity.Media{ID: id, OwnerType: "competition", ObjectKey: "competition/x/y.png", Kind: entity.KindImage}

	repo := mocks.NewMockRepository(t)
	storage := mocks.NewMockObjectStorage(t)

	repo.EXPECT().Get(mock.Anything, id).Return(stored, nil)
	storage.EXPECT().PresignGet(mock.Anything, stored.ObjectKey).Return("https://minio/presigned", nil)

	uc := biz.NewMedia(discardLogger(), repo, storage)

	got, err := uc.Get(ctx, id)

	require.NoError(t, err)
	require.Equal(t, id, got.Media.ID)
	require.Equal(t, "https://minio/presigned", got.URL)
}

func TestGet_NotFound(t *testing.T) {
	id := uuid.New()

	repo := mocks.NewMockRepository(t)
	storage := mocks.NewMockObjectStorage(t)

	repo.EXPECT().Get(mock.Anything, id).Return(entity.Media{}, biz.ErrResourceNotFound)

	uc := biz.NewMedia(discardLogger(), repo, storage)

	_, err := uc.Get(context.Background(), id)

	require.ErrorIs(t, err, biz.ErrResourceNotFound)
}

// Delete removes the DB row FIRST, then best-effort removes the object; an
// object-store failure must not fail the delete (invisible leak vs broken UI).
func TestDelete_RemovesRowThenObject(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	storage := mocks.NewMockObjectStorage(t)
	uc := biz.NewMedia(discardLogger(), repo, storage)

	id := uuid.New()
	m := entity.Media{ID: id, ObjectKey: "competition/x/img.png"}

	repo.EXPECT().Get(mock.Anything, id).Return(m, nil)
	repo.EXPECT().Delete(mock.Anything, id).Return(nil)
	storage.EXPECT().Remove(mock.Anything, "competition/x/img.png").Return(nil)

	require.NoError(t, uc.Delete(context.Background(), id))
}

func TestDelete_ObjectFailureStillSucceeds(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	storage := mocks.NewMockObjectStorage(t)
	uc := biz.NewMedia(discardLogger(), repo, storage)

	id := uuid.New()
	repo.EXPECT().Get(mock.Anything, id).Return(entity.Media{ID: id, ObjectKey: "k"}, nil)
	repo.EXPECT().Delete(mock.Anything, id).Return(nil)
	storage.EXPECT().Remove(mock.Anything, "k").Return(errors.New("minio down"))

	require.NoError(t, uc.Delete(context.Background(), id))
}

func TestDelete_DBFailureAborts(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	storage := mocks.NewMockObjectStorage(t)
	uc := biz.NewMedia(discardLogger(), repo, storage)

	id := uuid.New()
	repo.EXPECT().Get(mock.Anything, id).Return(entity.Media{ID: id, ObjectKey: "k"}, nil)
	repo.EXPECT().Delete(mock.Anything, id).Return(errors.New("db down"))
	// storage.Remove must NOT be called — no partial deletes.

	require.Error(t, uc.Delete(context.Background(), id))
}

func TestUpdate_Validation(t *testing.T) {
	uc := biz.NewMedia(discardLogger(), mocks.NewMockRepository(t), mocks.NewMockObjectStorage(t))

	_, err := uc.Update(context.Background(), uuid.New(), biz.UpdateInput{})
	require.ErrorIs(t, err, biz.ErrResourceInvalid)

	neg := -1
	_, err = uc.Update(context.Background(), uuid.New(), biz.UpdateInput{Position: &neg})
	require.ErrorIs(t, err, biz.ErrResourceInvalid)

	// owner_type without owner_id is rejected.
	_, err = uc.Update(context.Background(), uuid.New(), biz.UpdateInput{OwnerType: "competition"})
	require.ErrorIs(t, err, biz.ErrResourceInvalid)
}

func TestListAll_CapsLimit(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	storage := mocks.NewMockObjectStorage(t)
	repo.EXPECT().ListAll(mock.Anything, 200, 0).Return(biz.MediaPage{}, nil)

	uc := biz.NewMedia(discardLogger(), repo, storage)

	_, err := uc.ListAll(context.Background(), 9999, -5)
	require.NoError(t, err)
}
