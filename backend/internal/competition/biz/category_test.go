package biz_test

import (
	"context"
	"testing"

	"application/internal/competition/biz"
	"application/internal/competition/entity"
	"application/internal/competition/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCategory_CreateDerivesSlug(t *testing.T) {
	repo := mocks.NewMockRepositoryCategory(t)
	repo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(c entity.Category) bool {
		return c.Name == "Instant Wins" && c.Slug == "instant-wins"
	})).RunAndReturn(func(_ context.Context, c entity.Category) (entity.Category, error) {
		return c, nil
	})

	uc := biz.NewCategory(discardLogger(), repo)

	got, err := uc.Create(context.Background(), biz.CategoryInput{Name: "Instant Wins"})
	require.NoError(t, err)
	require.Equal(t, "instant-wins", got.Slug)
}

func TestCategory_CreateInvalid(t *testing.T) {
	uc := biz.NewCategory(discardLogger(), mocks.NewMockRepositoryCategory(t))

	_, err := uc.Create(context.Background(), biz.CategoryInput{Name: ""})
	require.ErrorIs(t, err, biz.ErrResourceInvalid)

	_, err = uc.Create(context.Background(), biz.CategoryInput{Name: "X", Slug: "Bad Slug!"})
	require.ErrorIs(t, err, biz.ErrResourceInvalid)
}

// Deleting an in-use category without a reassignment target is blocked (409
// at the API); the repo surfaces ErrCategoryInUse from its transaction.
func TestCategory_DeleteInUseBlocked(t *testing.T) {
	id := uuid.New()
	repo := mocks.NewMockRepositoryCategory(t)
	repo.EXPECT().Delete(mock.Anything, id, (*uuid.UUID)(nil)).Return(biz.ErrCategoryInUse)

	uc := biz.NewCategory(discardLogger(), repo)
	require.ErrorIs(t, uc.Delete(context.Background(), id, nil), biz.ErrCategoryInUse)
}

func TestCategory_DeleteWithReassign(t *testing.T) {
	id := uuid.New()
	target := uuid.New()
	repo := mocks.NewMockRepositoryCategory(t)
	repo.EXPECT().Get(mock.Anything, target).Return(entity.Category{ID: target, Name: "Cash"}, nil)
	repo.EXPECT().Delete(mock.Anything, id, &target).Return(nil)

	uc := biz.NewCategory(discardLogger(), repo)
	require.NoError(t, uc.Delete(context.Background(), id, &target))
}

func TestCategory_DeleteReassignToSelfRejected(t *testing.T) {
	id := uuid.New()
	uc := biz.NewCategory(discardLogger(), mocks.NewMockRepositoryCategory(t))

	require.ErrorIs(t, uc.Delete(context.Background(), id, &id), biz.ErrResourceInvalid)
}

func TestCategory_DeleteReassignTargetMustExist(t *testing.T) {
	id := uuid.New()
	target := uuid.New()
	repo := mocks.NewMockRepositoryCategory(t)
	repo.EXPECT().Get(mock.Anything, target).Return(entity.Category{}, biz.ErrCategoryNotFound)

	uc := biz.NewCategory(discardLogger(), repo)
	require.ErrorIs(t, uc.Delete(context.Background(), id, &target), biz.ErrCategoryNotFound)
}

func TestContent_UpsertValidatesKeyAndSize(t *testing.T) {
	repo := mocks.NewMockRepositoryContent(t)
	uc := biz.NewContent(discardLogger(), repo)

	_, err := uc.Upsert(context.Background(), "Bad Key!", "v")
	require.ErrorIs(t, err, biz.ErrResourceInvalid)

	big := make([]byte, 65*1024)
	_, err = uc.Upsert(context.Background(), "hero.slides", string(big))
	require.ErrorIs(t, err, biz.ErrResourceInvalid)

	repo.EXPECT().Upsert(mock.Anything, "hero.slides", "[]").
		Return(entity.SiteContent{Key: "hero.slides", Value: "[]"}, nil)

	got, err := uc.Upsert(context.Background(), "hero.slides", "[]")
	require.NoError(t, err)
	require.Equal(t, "[]", got.Value)
}
