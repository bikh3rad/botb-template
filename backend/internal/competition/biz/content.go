package biz

import (
	"context"
	"log/slog"
	"regexp"

	"application/internal/competition/entity"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// validContentKey bounds site_content keys: short, lowercase, dot/dash/underscore.
var validContentKey = regexp.MustCompile(`^[a-z0-9._-]{1,128}$`)

const maxContentValueBytes = 64 * 1024

// UsecaseContent is the site-copy store use case (public read, admin write).
type UsecaseContent interface {
	GetAll(ctx context.Context) ([]entity.SiteContent, error)
	Upsert(ctx context.Context, key, value string) (entity.SiteContent, error)
}

// RepositoryContent persists site_content rows.
type RepositoryContent interface {
	GetAll(ctx context.Context) ([]entity.SiteContent, error)
	Upsert(ctx context.Context, key, value string) (entity.SiteContent, error)
}

type content struct {
	logger *slog.Logger
	tracer trace.Tracer
	repo   RepositoryContent
}

var _ UsecaseContent = (*content)(nil)

// NewContent constructs the site-content use case.
func NewContent(logger *slog.Logger, repo RepositoryContent) *content {
	return &content{
		logger: logger.With("layer", "SiteContent"),
		tracer: otel.Tracer("SiteContentUseCase"),
		repo:   repo,
	}
}

func (uc *content) GetAll(ctx context.Context) ([]entity.SiteContent, error) {
	return uc.repo.GetAll(ctx)
}

func (uc *content) Upsert(ctx context.Context, key, value string) (entity.SiteContent, error) {
	ctx, span := uc.tracer.Start(ctx, "UpsertContent")
	defer span.End()

	if !validContentKey.MatchString(key) || len(value) > maxContentValueBytes {
		return entity.SiteContent{}, ErrResourceInvalid
	}

	return uc.repo.Upsert(ctx, key, value)
}
