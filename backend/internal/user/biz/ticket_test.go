package biz_test

import (
	"application/internal/user/biz"
	"application/internal/user/entity"
	"application/internal/user/mocks"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPurchase_Success(t *testing.T) {
	compID := uuid.New()
	userID := uuid.New()

	repo := mocks.NewMockRepositoryTicket(t)
	repo.EXPECT().Purchase(mock.Anything, biz.PurchaseInput{
		CompetitionID: compID, UserID: userID, Quantity: 3,
	}).Return(biz.PurchaseResult{
		Tickets:        []entity.Ticket{{}, {}, {}},
		User:           entity.User{ID: userID, TicketsOwned: 3, TotalSpentPence: 375},
		TotalCostPence: 375,
	}, nil)

	uc := biz.NewTicket(discardLogger(), repo)

	got, err := uc.Purchase(context.Background(), biz.PurchaseInput{
		CompetitionID: compID, UserID: userID, Quantity: 3,
	})
	require.NoError(t, err)
	require.Len(t, got.Tickets, 3)
	require.Equal(t, int64(375), got.TotalCostPence)
}

func TestPurchase_Invalid(t *testing.T) {
	repo := mocks.NewMockRepositoryTicket(t)
	uc := biz.NewTicket(discardLogger(), repo)

	cases := []biz.PurchaseInput{
		{CompetitionID: uuid.Nil, UserID: uuid.New(), Quantity: 1},
		{CompetitionID: uuid.New(), UserID: uuid.Nil, Quantity: 1},
		{CompetitionID: uuid.New(), UserID: uuid.New(), Quantity: 0},
		{CompetitionID: uuid.New(), UserID: uuid.New(), Quantity: 100000},
	}

	for _, in := range cases {
		_, err := uc.Purchase(context.Background(), in)
		require.ErrorIs(t, err, biz.ErrResourceInvalid)
	}
}

func TestListByUser_Invalid(t *testing.T) {
	repo := mocks.NewMockRepositoryTicket(t)
	uc := biz.NewTicket(discardLogger(), repo)

	_, err := uc.ListByUser(context.Background(), uuid.Nil)
	require.ErrorIs(t, err, biz.ErrResourceInvalid)
}

func TestListByUser_Success(t *testing.T) {
	userID := uuid.New()
	repo := mocks.NewMockRepositoryTicket(t)
	repo.EXPECT().ListByUser(mock.Anything, userID).
		Return([]entity.Ticket{{UserID: userID}}, nil)

	uc := biz.NewTicket(discardLogger(), repo)

	got, err := uc.ListByUser(context.Background(), userID)
	require.NoError(t, err)
	require.Len(t, got, 1)
}
