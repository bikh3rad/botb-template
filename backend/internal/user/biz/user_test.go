package biz_test

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"application/internal/user/biz"
	"application/internal/user/entity"
	"application/internal/user/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestRegister_Success(t *testing.T) {
	repo := mocks.NewMockRepositoryUser(t)
	repo.EXPECT().Create(mock.Anything, mock.Anything).
		RunAndReturn(func(_ context.Context, u entity.User) (entity.User, error) {
			return u, nil
		})

	uc := biz.NewUser(discardLogger(), repo)

	got, err := uc.Register(context.Background(), "  Olivia  ", "Olivia@Example.COM")
	require.NoError(t, err)
	require.Equal(t, "Olivia", got.Name)
	require.Equal(t, "olivia@example.com", got.Email) // trimmed + lowercased
	require.NotEqual(t, uuid.Nil, got.ID)
}

func TestRegister_Invalid(t *testing.T) {
	repo := mocks.NewMockRepositoryUser(t)
	uc := biz.NewUser(discardLogger(), repo)

	cases := [][2]string{
		{"", "a@b.com"},    // no name
		{"A", "not-email"}, // no @ / domain
		{"A", "a@b"},       // no TLD dot
		{"A", ""},          // empty email
	}

	for _, c := range cases {
		_, err := uc.Register(context.Background(), c[0], c[1])
		require.ErrorIs(t, err, biz.ErrResourceInvalid)
	}
}

func TestList_DefaultsAndCaps(t *testing.T) {
	repo := mocks.NewMockRepositoryUser(t)
	repo.EXPECT().List(mock.Anything, biz.UserListFilter{Query: "", Limit: 20, Offset: 0}).
		Return(biz.UserPage{Users: []entity.User{{Name: "A"}}, Total: 1}, nil)

	uc := biz.NewUser(discardLogger(), repo)

	page, err := uc.List(context.Background(), biz.UserListFilter{Limit: 0, Offset: -5})
	require.NoError(t, err)
	require.Equal(t, 1, page.Total)
}

func TestList_LimitCapped(t *testing.T) {
	repo := mocks.NewMockRepositoryUser(t)
	repo.EXPECT().List(mock.Anything, mock.MatchedBy(func(f biz.UserListFilter) bool {
		return f.Limit == 100
	})).Return(biz.UserPage{}, nil)

	uc := biz.NewUser(discardLogger(), repo)

	_, err := uc.List(context.Background(), biz.UserListFilter{Limit: 5000})
	require.NoError(t, err)
}

func TestGet_PassesThrough(t *testing.T) {
	id := uuid.New()
	repo := mocks.NewMockRepositoryUser(t)
	repo.EXPECT().Get(mock.Anything, id).Return(entity.User{ID: id}, nil)

	uc := biz.NewUser(discardLogger(), repo)

	got, err := uc.Get(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, got.ID)
}

func TestUpdate_ValidatesEmail(t *testing.T) {
	repo := mocks.NewMockRepositoryUser(t)
	uc := biz.NewUser(discardLogger(), repo)

	_, err := uc.Update(context.Background(), uuid.New(), "Name", "not-an-email")
	require.ErrorIs(t, err, biz.ErrResourceInvalid)

	_, err = uc.Update(context.Background(), uuid.New(), "", "a@b.co")
	require.ErrorIs(t, err, biz.ErrResourceInvalid)
}

func TestUpdate_NormalizesAndDelegates(t *testing.T) {
	id := uuid.New()
	repo := mocks.NewMockRepositoryUser(t)
	repo.EXPECT().Update(mock.Anything, id, "Sam", "sam@example.com").
		Return(entity.User{ID: id, Name: "Sam", Email: "sam@example.com", IsActive: true}, nil)

	uc := biz.NewUser(discardLogger(), repo)

	got, err := uc.Update(context.Background(), id, " Sam ", " SAM@Example.com ")
	require.NoError(t, err)
	require.Equal(t, "sam@example.com", got.Email)
}
