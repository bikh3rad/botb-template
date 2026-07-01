package biz

import (
	"application/internal/competition/entity"
	"context"
	"log/slog"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var nonSlugChars = regexp.MustCompile(`[^a-z0-9]+`)

type competition struct {
	logger *slog.Logger
	tracer trace.Tracer
	repo   Repository
}

var _ UsecaseCompetition = (*competition)(nil)

// NewCompetition constructs the competition use case.
func NewCompetition(logger *slog.Logger, repo Repository) *competition {
	return &competition{
		logger: logger.With("layer", "Competition"),
		tracer: otel.Tracer("CompetitionUseCase"),
		repo:   repo,
	}
}

func (uc *competition) List(ctx context.Context, filter ListFilter) ([]entity.Competition, error) {
	if filter.Status != nil && !filter.Status.Valid() {
		return nil, ErrResourceInvalid
	}

	return uc.repo.List(ctx, filter)
}

func (uc *competition) Get(ctx context.Context, id uuid.UUID) (entity.Competition, error) {
	return uc.repo.Get(ctx, id)
}

func (uc *competition) Create(ctx context.Context, in CreateInput) (entity.Competition, error) {
	ctx, span := uc.tracer.Start(ctx, "Create")
	defer span.End()

	if err := validateCreate(in); err != nil {
		return entity.Competition{}, err
	}

	slug := in.Slug
	if slug == "" {
		slug = Slugify(in.Title)
	}

	status := in.Status
	if status == "" {
		status = entity.StatusDraft
	}

	c := entity.Competition{
		ID:               uuid.New(),
		Title:            in.Title,
		Slug:             slug,
		Description:      in.Description,
		Prize:            in.Prize,
		TicketPricePence: in.TicketPricePence,
		TicketsTotal:     in.TicketsTotal,
		TicketsSold:      0,
		Status:           status,
		StartsAt:         in.StartsAt,
		EndsAt:           in.EndsAt,
	}

	return uc.repo.Create(ctx, c)
}

func (uc *competition) Update(ctx context.Context, id uuid.UUID, in UpdateInput) (entity.Competition, error) {
	ctx, span := uc.tracer.Start(ctx, "Update")
	defer span.End()

	if in.Title == "" || !in.Status.Valid() {
		return entity.Competition{}, ErrResourceInvalid
	}

	c := entity.Competition{
		ID:               id,
		Title:            in.Title,
		Description:      in.Description,
		Prize:            in.Prize,
		TicketPricePence: in.TicketPricePence,
		TicketsTotal:     in.TicketsTotal,
		Status:           in.Status,
		StartsAt:         in.StartsAt,
		EndsAt:           in.EndsAt,
	}

	return uc.repo.Update(ctx, c)
}

func (uc *competition) Delete(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(ctx, id)
}

func validateCreate(in CreateInput) error {
	if in.Title == "" || in.Prize == "" {
		return ErrResourceInvalid
	}

	if in.TicketsTotal <= 0 || in.TicketPricePence < 0 {
		return ErrResourceInvalid
	}

	if in.Status != "" && !in.Status.Valid() {
		return ErrResourceInvalid
	}

	return nil
}

// Slugify turns a title into a URL-safe slug, e.g. "Win an Audi RS3!" -> "win-an-audi-rs3".
func Slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = nonSlugChars.ReplaceAllString(s, "-")

	return strings.Trim(s, "-")
}
