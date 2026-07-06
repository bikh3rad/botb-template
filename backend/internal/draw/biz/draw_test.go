package biz_test

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"application/internal/draw/biz"
	"application/internal/draw/entity"
	"application/internal/draw/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestCreate_Success(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	repo.EXPECT().Create(mock.Anything, mock.Anything).
		RunAndReturn(func(_ context.Context, d entity.Draw) (entity.Draw, error) {
			return d, nil
		})

	uc := biz.NewDraw(discardLogger(), repo)

	got, err := uc.Create(context.Background(), biz.CreateInput{CompetitionID: uuid.New(), Prize: "Audi RS3"})
	require.NoError(t, err)
	require.Equal(t, entity.StatusPending, got.Status)
	require.NotEqual(t, uuid.Nil, got.ID)
}

func TestCreate_Invalid(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	uc := biz.NewDraw(discardLogger(), repo)

	cases := []biz.CreateInput{
		{CompetitionID: uuid.Nil, Prize: "x"},
		{CompetitionID: uuid.New(), Prize: ""},
	}

	for _, in := range cases {
		_, err := uc.Create(context.Background(), in)
		require.ErrorIs(t, err, biz.ErrResourceInvalid)
	}
}

func TestRun_HappyPath(t *testing.T) {
	id := uuid.New()
	winner := uuid.New()
	repo := mocks.NewMockRepository(t)
	repo.EXPECT().Run(mock.Anything, id).
		Return(entity.Draw{ID: id, Status: entity.StatusDrawn, WinnerUserID: &winner}, nil)

	uc := biz.NewDraw(discardLogger(), repo)

	got, err := uc.Run(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, entity.StatusDrawn, got.Status)
	require.NotNil(t, got.WinnerUserID)
}

func TestRun_AlreadyDrawn(t *testing.T) {
	id := uuid.New()
	repo := mocks.NewMockRepository(t)
	repo.EXPECT().Run(mock.Anything, id).Return(entity.Draw{}, biz.ErrAlreadyDrawn)

	uc := biz.NewDraw(discardLogger(), repo)

	_, err := uc.Run(context.Background(), id)
	require.ErrorIs(t, err, biz.ErrAlreadyDrawn)
}

func TestRun_NoTickets(t *testing.T) {
	id := uuid.New()
	repo := mocks.NewMockRepository(t)
	repo.EXPECT().Run(mock.Anything, id).Return(entity.Draw{}, biz.ErrNoTickets)

	uc := biz.NewDraw(discardLogger(), repo)

	_, err := uc.Run(context.Background(), id)
	require.ErrorIs(t, err, biz.ErrNoTickets)
}

func TestRun_InvalidID(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	uc := biz.NewDraw(discardLogger(), repo)

	_, err := uc.Run(context.Background(), uuid.Nil)
	require.ErrorIs(t, err, biz.ErrResourceInvalid)
}

func TestGetPublic_HidesPending(t *testing.T) {
	id := uuid.New()
	repo := mocks.NewMockRepository(t)
	repo.EXPECT().Get(mock.Anything, id).Return(entity.Draw{ID: id, Status: entity.StatusPending}, nil)

	uc := biz.NewDraw(discardLogger(), repo)

	_, err := uc.GetPublic(context.Background(), id)
	require.ErrorIs(t, err, biz.ErrResourceNotFound)
}

func TestGetPublic_ReturnsDrawn(t *testing.T) {
	id := uuid.New()
	repo := mocks.NewMockRepository(t)
	repo.EXPECT().Get(mock.Anything, id).Return(entity.Draw{ID: id, Status: entity.StatusDrawn}, nil)

	uc := biz.NewDraw(discardLogger(), repo)

	got, err := uc.GetPublic(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, entity.StatusDrawn, got.Status)
}

func TestList_DefaultsAndCaps(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	repo.EXPECT().List(mock.Anything, mock.MatchedBy(func(f biz.ListFilter) bool {
		return f.Limit == 20 && f.Offset == 0
	})).Return(biz.DrawPage{Total: 0}, nil)

	uc := biz.NewDraw(discardLogger(), repo)

	_, err := uc.List(context.Background(), biz.ListFilter{Limit: 0, Offset: -3})
	require.NoError(t, err)
}

func TestVoid_RequiresReason(t *testing.T) {
	uc := biz.NewDraw(discardLogger(), mocks.NewMockRepository(t))

	_, err := uc.Void(context.Background(), uuid.New(), "   ")
	require.ErrorIs(t, err, biz.ErrReasonRequired)
}

func TestVoid_DelegatesWithReason(t *testing.T) {
	id := uuid.New()
	repo := mocks.NewMockRepository(t)
	repo.EXPECT().Void(mock.Anything, id, "winner unreachable").
		Return(entity.Draw{ID: id, Status: entity.StatusVoid, VoidReason: "winner unreachable"}, nil)

	uc := biz.NewDraw(discardLogger(), repo)

	got, err := uc.Void(context.Background(), id, "winner unreachable")
	require.NoError(t, err)
	require.Equal(t, entity.StatusVoid, got.Status)
}

func TestReassign_RequiresReason(t *testing.T) {
	uc := biz.NewDraw(discardLogger(), mocks.NewMockRepository(t))

	_, err := uc.Reassign(context.Background(), uuid.New(), uuid.New(), "")
	require.ErrorIs(t, err, biz.ErrReasonRequired)
}

func TestUpdatePrize_VoidDrawFrozen(t *testing.T) {
	id := uuid.New()
	repo := mocks.NewMockRepository(t)
	repo.EXPECT().Get(mock.Anything, id).Return(entity.Draw{ID: id, Status: entity.StatusVoid}, nil)

	uc := biz.NewDraw(discardLogger(), repo)

	_, err := uc.UpdatePrize(context.Background(), id, "New prize")
	require.ErrorIs(t, err, biz.ErrInvalidState)
}

// Void draws are hidden from the public read alongside pending ones.
func TestGetPublic_HidesVoid(t *testing.T) {
	id := uuid.New()
	repo := mocks.NewMockRepository(t)
	repo.EXPECT().Get(mock.Anything, id).Return(entity.Draw{ID: id, Status: entity.StatusVoid}, nil)

	uc := biz.NewDraw(discardLogger(), repo)

	_, err := uc.GetPublic(context.Background(), id)
	require.ErrorIs(t, err, biz.ErrResourceNotFound)
}
