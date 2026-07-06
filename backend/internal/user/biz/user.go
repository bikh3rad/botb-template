package biz

import (
	"context"
	"log/slog"
	"strings"

	"application/internal/user/entity"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const (
	defaultPageSize = 20
	maxPageSize     = 100
)

type user struct {
	logger *slog.Logger
	tracer trace.Tracer
	repo   RepositoryUser
}

var _ UsecaseUser = (*user)(nil)

// NewUser constructs the user use case.
func NewUser(logger *slog.Logger, repo RepositoryUser) *user {
	return &user{
		logger: logger.With("layer", "User"),
		tracer: otel.Tracer("UserUseCase"),
		repo:   repo,
	}
}

func (uc *user) Register(ctx context.Context, name, email string) (entity.User, error) {
	ctx, span := uc.tracer.Start(ctx, "Register")
	defer span.End()

	name = strings.TrimSpace(name)
	email = strings.TrimSpace(strings.ToLower(email))

	if name == "" || !validEmail(email) {
		return entity.User{}, ErrResourceInvalid
	}

	u := entity.User{
		ID:    uuid.New(),
		Name:  name,
		Email: email,
	}

	return uc.repo.Create(ctx, u)
}

func (uc *user) List(ctx context.Context, filter UserListFilter) (UserPage, error) {
	if filter.Limit <= 0 {
		filter.Limit = defaultPageSize
	}

	if filter.Limit > maxPageSize {
		filter.Limit = maxPageSize
	}

	if filter.Offset < 0 {
		filter.Offset = 0
	}

	filter.Query = strings.TrimSpace(filter.Query)

	return uc.repo.List(ctx, filter)
}

func (uc *user) Get(ctx context.Context, id uuid.UUID) (entity.User, error) {
	return uc.repo.Get(ctx, id)
}

// Update edits the admin-editable profile fields (name, email). Email is
// normalized and validated; a duplicate maps to ErrResourceExists (409).
func (uc *user) Update(ctx context.Context, id uuid.UUID, name, email string) (entity.User, error) {
	ctx, span := uc.tracer.Start(ctx, "Update")
	defer span.End()

	name = strings.TrimSpace(name)
	email = strings.TrimSpace(strings.ToLower(email))

	if name == "" || !validEmail(email) {
		return entity.User{}, ErrResourceInvalid
	}

	return uc.repo.Update(ctx, id, name, email)
}

// SetActive suspends (false) or reactivates (true) a user. Suspension is the
// soft alternative to deletion: history stays, purchases are blocked.
func (uc *user) SetActive(ctx context.Context, id uuid.UUID, active bool) (entity.User, error) {
	ctx, span := uc.tracer.Start(ctx, "SetActive")
	defer span.End()

	return uc.repo.SetActive(ctx, id, active)
}

// validEmail is a deliberately minimal check — enough to reject obvious junk
// without pulling in a full RFC 5322 validator.
func validEmail(email string) bool {
	at := strings.IndexByte(email, '@')

	return at > 0 && at < len(email)-1 && strings.IndexByte(email[at+1:], '.') > 0
}
