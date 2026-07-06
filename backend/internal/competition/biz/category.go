package biz

import (
	"context"
	"log/slog"

	"application/internal/competition/entity"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// CategoryInput is the create/update payload for a category. Slug is derived
// from Name when empty.
type CategoryInput struct {
	Name string
	Slug string
}

// UsecaseCategory is consumed by the HTTP handler.
type UsecaseCategory interface {
	List(ctx context.Context) ([]entity.Category, error)
	Create(ctx context.Context, in CategoryInput) (entity.Category, error)
	Update(ctx context.Context, id uuid.UUID, in CategoryInput) (entity.Category, error)
	// Delete removes a category. When competitions still reference it:
	// reassignTo == nil -> ErrCategoryInUse; otherwise competitions are moved
	// to reassignTo atomically before the delete. Competitions are never
	// orphaned by a category delete.
	Delete(ctx context.Context, id uuid.UUID, reassignTo *uuid.UUID) error
}

// RepositoryCategory persists categories.
type RepositoryCategory interface {
	List(ctx context.Context) ([]entity.Category, error)
	Get(ctx context.Context, id uuid.UUID) (entity.Category, error)
	Create(ctx context.Context, c entity.Category) (entity.Category, error)
	Update(ctx context.Context, c entity.Category) (entity.Category, error)
	// Delete implements the in-use guard + optional reassignment in one
	// transaction and returns ErrCategoryInUse when blocked.
	Delete(ctx context.Context, id uuid.UUID, reassignTo *uuid.UUID) error
}

type category struct {
	logger *slog.Logger
	tracer trace.Tracer
	repo   RepositoryCategory
}

var _ UsecaseCategory = (*category)(nil)

// NewCategory constructs the category use case.
func NewCategory(logger *slog.Logger, repo RepositoryCategory) *category {
	return &category{
		logger: logger.With("layer", "Category"),
		tracer: otel.Tracer("CategoryUseCase"),
		repo:   repo,
	}
}

func (uc *category) List(ctx context.Context) ([]entity.Category, error) {
	return uc.repo.List(ctx)
}

func (uc *category) Create(ctx context.Context, in CategoryInput) (entity.Category, error) {
	ctx, span := uc.tracer.Start(ctx, "CreateCategory")
	defer span.End()

	c, err := categoryFromInput(uuid.New(), in)
	if err != nil {
		return entity.Category{}, err
	}

	return uc.repo.Create(ctx, c)
}

func (uc *category) Update(ctx context.Context, id uuid.UUID, in CategoryInput) (entity.Category, error) {
	ctx, span := uc.tracer.Start(ctx, "UpdateCategory")
	defer span.End()

	c, err := categoryFromInput(id, in)
	if err != nil {
		return entity.Category{}, err
	}

	return uc.repo.Update(ctx, c)
}

func (uc *category) Delete(ctx context.Context, id uuid.UUID, reassignTo *uuid.UUID) error {
	ctx, span := uc.tracer.Start(ctx, "DeleteCategory")
	defer span.End()

	if reassignTo != nil {
		if *reassignTo == id {
			return ErrResourceInvalid
		}

		if _, err := uc.repo.Get(ctx, *reassignTo); err != nil {
			return err
		}
	}

	return uc.repo.Delete(ctx, id, reassignTo)
}

func categoryFromInput(id uuid.UUID, in CategoryInput) (entity.Category, error) {
	if in.Name == "" {
		return entity.Category{}, ErrResourceInvalid
	}

	slug := in.Slug
	if slug == "" {
		slug = Slugify(in.Name)
	}

	if !validSlug.MatchString(slug) {
		return entity.Category{}, ErrResourceInvalid
	}

	return entity.Category{ID: id, Name: in.Name, Slug: slug}, nil
}
