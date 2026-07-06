package biz_test

import (
	"context"
	"testing"

	"application/internal/adminauth/biz"
	"application/internal/adminauth/entity"
	"application/internal/adminauth/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func newAccounts(t *testing.T) (biz.UsecaseAccounts, *mocks.MockRepository) {
	t.Helper()

	repo := mocks.NewMockRepository(t)

	return biz.NewAccounts(discardLogger(), repo), repo
}

func TestCreateAccount_Valid(t *testing.T) {
	uc, repo := newAccounts(t)

	repo.EXPECT().CreateAccount(mock.Anything, mock.Anything).
		RunAndReturn(func(_ context.Context, a entity.AdminAccount) (entity.AdminAccount, error) {
			// Password must be stored hashed, never plaintext.
			require.NotEqual(t, "hunter2-longer", a.PasswordHash)
			require.NoError(t, bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte("hunter2-longer")))
			require.Equal(t, "ops@example.com", a.Email)
			require.True(t, a.IsActive)

			return a, nil
		})

	account, err := uc.Create(context.Background(), biz.CreateAccountInput{
		Name:     "Ops",
		Email:    "Ops@Example.com",
		Password: "hunter2-longer",
		Role:     entity.RoleAdmin,
	})
	require.NoError(t, err)
	require.Equal(t, entity.RoleAdmin, account.Role)
}

func TestCreateAccount_Invalid(t *testing.T) {
	uc, _ := newAccounts(t)

	cases := []biz.CreateAccountInput{
		{Name: "", Email: "a@b.co", Password: "long-enough", Role: entity.RoleAdmin},
		{Name: "X", Email: "not-an-email", Password: "long-enough", Role: entity.RoleAdmin},
		{Name: "X", Email: "a@b.co", Password: "short", Role: entity.RoleAdmin},
		{Name: "X", Email: "a@b.co", Password: "long-enough", Role: entity.Role("root")},
	}

	for _, input := range cases {
		_, err := uc.Create(context.Background(), input)
		require.ErrorIs(t, err, biz.ErrResourceInvalid)
	}
}

func TestUpdateAccount_DisableLastSuperadminRefused(t *testing.T) {
	uc, repo := newAccounts(t)
	acc := entity.AdminAccount{
		ID: uuid.New(), Name: "Root", Email: "root@example.com",
		Role: entity.RoleSuperadmin, IsActive: true,
	}
	inactive := false

	repo.EXPECT().GetAccount(mock.Anything, acc.ID).Return(acc, nil)
	repo.EXPECT().CountActiveSuperadmins(mock.Anything).Return(1, nil)

	_, err := uc.Update(context.Background(), acc.ID, biz.UpdateAccountInput{IsActive: &inactive})
	require.ErrorIs(t, err, biz.ErrLastSuperadmin)
}

func TestUpdateAccount_DemoteLastSuperadminRefused(t *testing.T) {
	uc, repo := newAccounts(t)
	acc := entity.AdminAccount{
		ID: uuid.New(), Name: "Root", Email: "root@example.com",
		Role: entity.RoleSuperadmin, IsActive: true,
	}
	demoted := entity.RoleAdmin

	repo.EXPECT().GetAccount(mock.Anything, acc.ID).Return(acc, nil)
	repo.EXPECT().CountActiveSuperadmins(mock.Anything).Return(1, nil)

	_, err := uc.Update(context.Background(), acc.ID, biz.UpdateAccountInput{Role: &demoted})
	require.ErrorIs(t, err, biz.ErrLastSuperadmin)
}

func TestUpdateAccount_DisableSuperadminAllowedWhenAnotherExists(t *testing.T) {
	uc, repo := newAccounts(t)
	acc := entity.AdminAccount{
		ID: uuid.New(), Name: "Root", Email: "root@example.com",
		Role: entity.RoleSuperadmin, IsActive: true,
	}
	inactive := false

	repo.EXPECT().GetAccount(mock.Anything, acc.ID).Return(acc, nil)
	repo.EXPECT().CountActiveSuperadmins(mock.Anything).Return(2, nil)
	repo.EXPECT().UpdateAccount(mock.Anything, mock.Anything).
		RunAndReturn(func(_ context.Context, a entity.AdminAccount) (entity.AdminAccount, error) {
			require.False(t, a.IsActive)

			return a, nil
		})

	_, err := uc.Update(context.Background(), acc.ID, biz.UpdateAccountInput{IsActive: &inactive})
	require.NoError(t, err)
}

func TestUpdateAccount_DisablePlainAdminNeedsNoGuard(t *testing.T) {
	uc, repo := newAccounts(t)
	acc := entity.AdminAccount{
		ID: uuid.New(), Name: "Ops", Email: "ops@example.com",
		Role: entity.RoleAdmin, IsActive: true,
	}
	inactive := false

	repo.EXPECT().GetAccount(mock.Anything, acc.ID).Return(acc, nil)
	repo.EXPECT().UpdateAccount(mock.Anything, mock.Anything).
		RunAndReturn(func(_ context.Context, a entity.AdminAccount) (entity.AdminAccount, error) {
			return a, nil
		})

	_, err := uc.Update(context.Background(), acc.ID, biz.UpdateAccountInput{IsActive: &inactive})
	require.NoError(t, err)
}

func TestEnsureBootstrap_CreatesFirstSuperadmin(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	cfg := &biz.Config{Bootstrap: biz.Bootstrap{Email: "root@example.com", Password: "first-password"}}

	repo.EXPECT().CountActiveSuperadmins(mock.Anything).Return(0, nil)
	repo.EXPECT().CreateAccount(mock.Anything, mock.Anything).
		RunAndReturn(func(_ context.Context, a entity.AdminAccount) (entity.AdminAccount, error) {
			require.Equal(t, entity.RoleSuperadmin, a.Role)
			require.Equal(t, "root@example.com", a.Email)
			require.NoError(t, bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte("first-password")))

			return a, nil
		})

	require.NoError(t, biz.EnsureBootstrap(context.Background(), discardLogger(), cfg, repo))
}

func TestEnsureBootstrap_SkipsWhenSuperadminExists(t *testing.T) {
	repo := mocks.NewMockRepository(t)
	cfg := &biz.Config{Bootstrap: biz.Bootstrap{Email: "root@example.com", Password: "first-password"}}

	repo.EXPECT().CountActiveSuperadmins(mock.Anything).Return(1, nil)

	require.NoError(t, biz.EnsureBootstrap(context.Background(), discardLogger(), cfg, repo))
}

func TestEnsureBootstrap_SkipsWhenUnconfigured(t *testing.T) {
	repo := mocks.NewMockRepository(t)

	require.NoError(t, biz.EnsureBootstrap(context.Background(), discardLogger(), &biz.Config{}, repo))
}
