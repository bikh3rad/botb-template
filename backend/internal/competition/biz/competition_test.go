package biz_test

import (
	"application/internal/competition/biz"
	"application/internal/competition/entity"
	"application/internal/competition/mocks"
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestSlugify(t *testing.T) {
	require.Equal(t, "win-an-audi-rs3", biz.Slugify("  Win an Audi RS3!  "))
	require.Equal(t, "1-2m-home-in-zone-1", biz.Slugify("£1.2M Home in Zone 1"))
}

func TestCreate_GeneratesSlug(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	repo.EXPECT().Create(mock.Anything, mock.Anything).
		RunAndReturn(func(_ context.Context, c entity.Competition) (entity.Competition, error) {
			return c, nil
		})

	uc := biz.NewCompetition(discardLogger(), repo)

	got, err := uc.Create(context.Background(), biz.CreateInput{
		Title:        "Win an Audi RS3",
		Prize:        "Audi RS3",
		TicketsTotal: 1000,
	})

	require.NoError(t, err)
	require.Equal(t, "win-an-audi-rs3", got.Slug)
	require.Equal(t, entity.StatusDraft, got.Status)
	require.NotEqual(t, uuid.Nil, got.ID)
}

func TestCreate_InvalidInput(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	uc := biz.NewCompetition(discardLogger(), repo)

	cases := []biz.CreateInput{
		{Title: "", Prize: "x", TicketsTotal: 10}, // no title
		{Title: "x", Prize: "", TicketsTotal: 10}, // no prize
		{Title: "x", Prize: "y", TicketsTotal: 0}, // no tickets
		{Title: "x", Prize: "y", TicketsTotal: 1, Status: "bogus"},
	}

	for _, in := range cases {
		_, err := uc.Create(context.Background(), in)
		require.ErrorIs(t, err, biz.ErrResourceInvalid)
	}
}

func TestList_InvalidStatusFilter(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	uc := biz.NewCompetition(discardLogger(), repo)

	bogus := entity.Status("bogus")
	_, err := uc.List(context.Background(), biz.ListFilter{Status: &bogus})
	require.ErrorIs(t, err, biz.ErrResourceInvalid)
}

func TestList_PassesFilterThrough(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	live := entity.StatusLive
	repo.EXPECT().List(mock.Anything, biz.ListFilter{Status: &live}).
		Return([]entity.Competition{{Title: "A"}}, nil)

	uc := biz.NewCompetition(discardLogger(), repo)

	got, err := uc.List(context.Background(), biz.ListFilter{Status: &live})
	require.NoError(t, err)
	require.Len(t, got, 1)
}

func TestUpdate_InvalidStatus(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	uc := biz.NewCompetition(discardLogger(), repo)

	_, err := uc.Update(context.Background(), uuid.New(), biz.UpdateInput{Title: "x", Status: "nope"})
	require.ErrorIs(t, err, biz.ErrResourceInvalid)
}

func TestUpdate_Success(t *testing.T) {
	id := uuid.New()
	repo := mocks.NewMockRepository(t)
	repo.EXPECT().Update(mock.Anything, mock.MatchedBy(func(c entity.Competition) bool {
		return c.ID == id && c.Title == "New Title"
	})).Return(entity.Competition{ID: id, Title: "New Title", Status: entity.StatusLive}, nil)

	uc := biz.NewCompetition(discardLogger(), repo)

	got, err := uc.Update(context.Background(), id, biz.UpdateInput{Title: "New Title", Status: entity.StatusLive})
	require.NoError(t, err)
	require.Equal(t, "New Title", got.Title)
}

func TestDelete(t *testing.T) {
	id := uuid.New()
	repo := mocks.NewMockRepository(t)
	repo.EXPECT().Delete(mock.Anything, id).Return(nil)

	uc := biz.NewCompetition(discardLogger(), repo)
	require.NoError(t, uc.Delete(context.Background(), id))
}
