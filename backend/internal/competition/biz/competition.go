package biz

import (
	"context"
	"log/slog"
	"regexp"
	"strings"

	"application/internal/competition/entity"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	nonSlugChars = regexp.MustCompile(`[^a-z0-9]+`)
	validSlug    = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
)

// ObjectStorage is the subset of the S3/MinIO port the competition service
// needs to purge a deleted competition's media objects. Satisfied by
// datasource.MinioStorage (the same client the media service uses).
type ObjectStorage interface {
	Remove(ctx context.Context, key string) error
}

type competition struct {
	logger  *slog.Logger
	tracer  trace.Tracer
	repo    Repository
	storage ObjectStorage
}

var _ UsecaseCompetition = (*competition)(nil)

// NewCompetition constructs the competition use case.
func NewCompetition(logger *slog.Logger, repo Repository, storage ObjectStorage) *competition {
	return &competition{
		logger:  logger.With("layer", "Competition"),
		tracer:  otel.Tracer("CompetitionUseCase"),
		repo:    repo,
		storage: storage,
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

	if !validSlug.MatchString(slug) {
		return entity.Competition{}, ErrResourceInvalid
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
		CategoryID:       in.CategoryID,
		Status:           status,
		StartsAt:         in.StartsAt,
		EndsAt:           in.EndsAt,
	}

	return uc.repo.Create(ctx, c)
}

// Update applies a full-field edit. It validates against the CURRENT row:
// status may only move forward (draft->live->closed, or stay put — a closed
// competition never reopens through this endpoint), and tickets_total may not
// drop below tickets_sold. tickets_sold itself is never written here.
func (uc *competition) Update(ctx context.Context, id uuid.UUID, in UpdateInput) (entity.Competition, error) {
	ctx, span := uc.tracer.Start(ctx, "Update")
	defer span.End()

	if in.Title == "" || in.Prize == "" || !in.Status.Valid() {
		return entity.Competition{}, ErrResourceInvalid
	}

	if in.TicketsTotal <= 0 || in.TicketPricePence < 0 {
		return entity.Competition{}, ErrResourceInvalid
	}

	slug := in.Slug
	if slug == "" {
		slug = Slugify(in.Title)
	}

	if !validSlug.MatchString(slug) {
		return entity.Competition{}, ErrResourceInvalid
	}

	current, err := uc.repo.Get(ctx, id)
	if err != nil {
		return entity.Competition{}, err
	}

	if !transitionAllowed(current.Status, in.Status) {
		return entity.Competition{}, ErrInvalidTransition
	}

	if in.TicketsTotal < current.TicketsSold {
		return entity.Competition{}, ErrResourceInvalid
	}

	c := entity.Competition{
		ID:               id,
		Title:            in.Title,
		Slug:             slug,
		Description:      in.Description,
		Prize:            in.Prize,
		TicketPricePence: in.TicketPricePence,
		TicketsTotal:     in.TicketsTotal,
		CategoryID:       in.CategoryID,
		Status:           in.Status,
		StartsAt:         in.StartsAt,
		EndsAt:           in.EndsAt,
	}

	return uc.repo.Update(ctx, c)
}

// transitionAllowed encodes the status lifecycle: forward-only. Reopening a
// closed competition requires deliberate manual intervention, not a PUT.
func transitionAllowed(from, to entity.Status) bool {
	if from == to {
		return true
	}

	switch from {
	case entity.StatusDraft:
		return to == entity.StatusLive
	case entity.StatusLive:
		return to == entity.StatusClosed
	case entity.StatusClosed:
		return false
	default:
		return false
	}
}

// Delete removes a competition ONLY when it has no entrants — zero sold tickets
// AND no draws. Entrants have paid and draws reference the competition, so a
// live/played competition must be CLOSED or its draw VOIDED instead
// (ErrCompetitionHasEntrants -> 409). On a permitted delete the repo also
// removes the competition's media rows in the same transaction and returns
// their object keys, which we then purge from object storage (DB-first, like
// the media service's own delete ordering).
func (uc *competition) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, span := uc.tracer.Start(ctx, "Delete")
	defer span.End()

	objectKeys, err := uc.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Best-effort object purge: the DB rows are already gone, so a failed
	// removal is a logged storage leak, never a dangling record.
	for _, key := range objectKeys {
		if uc.storage == nil {
			break
		}

		if rmErr := uc.storage.Remove(ctx, key); rmErr != nil {
			uc.logger.WarnContext(ctx, "competition media object removal failed (orphaned in storage)",
				"object_key", key, "error", rmErr)
		}
	}

	return nil
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
