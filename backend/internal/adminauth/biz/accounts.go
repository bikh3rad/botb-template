package biz

import (
	"context"
	"log/slog"
	"strings"

	"application/internal/adminauth/entity"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
)

const minPasswordLength = 8

type accounts struct {
	logger *slog.Logger
	tracer trace.Tracer
	repo   Repository
}

var _ UsecaseAccounts = (*accounts)(nil)

// NewAccounts constructs the superadmin account-management use case.
func NewAccounts(logger *slog.Logger, repo Repository) *accounts {
	return &accounts{
		logger: logger.With("layer", "AdminAccounts"),
		tracer: otel.Tracer("AdminAccountsUseCase"),
		repo:   repo,
	}
}

func (uc *accounts) List(ctx context.Context) ([]entity.AdminAccount, error) {
	return uc.repo.ListAccounts(ctx)
}

func (uc *accounts) Create(ctx context.Context, input CreateAccountInput) (entity.AdminAccount, error) {
	ctx, span := uc.tracer.Start(ctx, "CreateAccount")
	defer span.End()

	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))

	if input.Name == "" || !validEmail(input.Email) ||
		len(input.Password) < minPasswordLength || !input.Role.Valid() {
		return entity.AdminAccount{}, ErrResourceInvalid
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return entity.AdminAccount{}, err
	}

	account := entity.AdminAccount{
		ID:           uuid.New(),
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: string(hash),
		Role:         input.Role,
		IsActive:     true,
	}

	return uc.repo.CreateAccount(ctx, account)
}

// Update applies a partial edit. Disabling or demoting the LAST ACTIVE
// superadmin is refused (ErrLastSuperadmin) — there is no hard delete, so this
// guard is what keeps account management reachable.
func (uc *accounts) Update(ctx context.Context, id uuid.UUID, input UpdateAccountInput) (entity.AdminAccount, error) {
	ctx, span := uc.tracer.Start(ctx, "UpdateAccount")
	defer span.End()

	account, err := uc.repo.GetAccount(ctx, id)
	if err != nil {
		return entity.AdminAccount{}, err
	}

	losesSuperadmin := false

	if input.Name != nil {
		account.Name = strings.TrimSpace(*input.Name)
		if account.Name == "" {
			return entity.AdminAccount{}, ErrResourceInvalid
		}
	}

	if input.Email != nil {
		account.Email = strings.TrimSpace(strings.ToLower(*input.Email))
		if !validEmail(account.Email) {
			return entity.AdminAccount{}, ErrResourceInvalid
		}
	}

	if input.Password != nil {
		if len(*input.Password) < minPasswordLength {
			return entity.AdminAccount{}, ErrResourceInvalid
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		if err != nil {
			return entity.AdminAccount{}, err
		}

		account.PasswordHash = string(hash)
	}

	if input.Role != nil {
		if !input.Role.Valid() {
			return entity.AdminAccount{}, ErrResourceInvalid
		}

		if account.Role == entity.RoleSuperadmin && *input.Role != entity.RoleSuperadmin && account.IsActive {
			losesSuperadmin = true
		}

		account.Role = *input.Role
	}

	if input.IsActive != nil {
		if account.Role == entity.RoleSuperadmin && account.IsActive && !*input.IsActive {
			losesSuperadmin = true
		}

		account.IsActive = *input.IsActive
	}

	if losesSuperadmin {
		count, err := uc.repo.CountActiveSuperadmins(ctx)
		if err != nil {
			return entity.AdminAccount{}, err
		}

		if count <= 1 {
			return entity.AdminAccount{}, ErrLastSuperadmin
		}
	}

	return uc.repo.UpdateAccount(ctx, account)
}

// EnsureBootstrap creates the FIRST superadmin when (a) bootstrap email +
// password are configured (APP_ADMINAUTH_BOOTSTRAP_EMAIL / _PASSWORD /
// optional _NAME) and (b) no superadmin account exists yet. Idempotent, never
// overwrites an existing account, and logs the email but NEVER the password.
// This is the documented seed path — there are no hardcoded credentials.
func EnsureBootstrap(ctx context.Context, logger *slog.Logger, cfg *Config, repo Repository) error {
	logger = logger.With("layer", "AdminAuthBootstrap")

	if cfg.Bootstrap.Email == "" || cfg.Bootstrap.Password == "" {
		logger.InfoContext(ctx, "no bootstrap credentials configured; skipping first-superadmin seed")

		return nil
	}

	count, err := repo.CountActiveSuperadmins(ctx)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	name := cfg.Bootstrap.Name
	if name == "" {
		name = "Superadmin"
	}

	if len(cfg.Bootstrap.Password) < minPasswordLength {
		logger.ErrorContext(ctx, "bootstrap password too short; refusing to seed", "min_length", minPasswordLength)

		return ErrResourceInvalid
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.Bootstrap.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = repo.CreateAccount(ctx, entity.AdminAccount{
		ID:           uuid.New(),
		Name:         name,
		Email:        strings.TrimSpace(strings.ToLower(cfg.Bootstrap.Email)),
		PasswordHash: string(hash),
		Role:         entity.RoleSuperadmin,
		IsActive:     true,
	})
	if err != nil {
		return err
	}

	logger.InfoContext(ctx, "bootstrapped first superadmin", "email", cfg.Bootstrap.Email)

	return nil
}

// validEmail is a deliberately minimal check — mirrors the user service.
func validEmail(email string) bool {
	at := strings.IndexByte(email, '@')

	return at > 0 && at < len(email)-1 && strings.IndexByte(email[at+1:], '.') > 0
}
