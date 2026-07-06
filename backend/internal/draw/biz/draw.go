package biz

import (
	"context"
	"log/slog"
	"strings"

	"application/internal/draw/entity"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const (
	defaultPageSize = 20
	maxPageSize     = 100
)

type draw struct {
	logger *slog.Logger
	tracer trace.Tracer
	repo   Repository
}

var _ UsecaseDraw = (*draw)(nil)

// NewDraw constructs the draw use case.
func NewDraw(logger *slog.Logger, repo Repository) *draw {
	return &draw{
		logger: logger.With("layer", "Draw"),
		tracer: otel.Tracer("DrawUseCase"),
		repo:   repo,
	}
}

func (uc *draw) List(ctx context.Context, filter ListFilter) (DrawPage, error) {
	if filter.Limit <= 0 {
		filter.Limit = defaultPageSize
	}

	if filter.Limit > maxPageSize {
		filter.Limit = maxPageSize
	}

	if filter.Offset < 0 {
		filter.Offset = 0
	}

	return uc.repo.List(ctx, filter)
}

func (uc *draw) Get(ctx context.Context, id uuid.UUID) (entity.Draw, error) {
	return uc.repo.Get(ctx, id)
}

// GetPublic hides pending AND void draws from public callers — a voided
// result must not show a winner.
func (uc *draw) GetPublic(ctx context.Context, id uuid.UUID) (entity.Draw, error) {
	d, err := uc.repo.Get(ctx, id)
	if err != nil {
		return entity.Draw{}, err
	}

	if d.Status != entity.StatusDrawn {
		return entity.Draw{}, ErrResourceNotFound
	}

	return d, nil
}

func (uc *draw) Create(ctx context.Context, in CreateInput) (entity.Draw, error) {
	ctx, span := uc.tracer.Start(ctx, "Create")
	defer span.End()

	if in.CompetitionID == uuid.Nil || in.Prize == "" {
		return entity.Draw{}, ErrResourceInvalid
	}

	d := entity.Draw{
		ID:            uuid.New(),
		CompetitionID: in.CompetitionID,
		Prize:         in.Prize,
		Status:        entity.StatusPending,
	}

	return uc.repo.Create(ctx, d)
}

func (uc *draw) Run(ctx context.Context, id uuid.UUID) (entity.Draw, error) {
	ctx, span := uc.tracer.Start(ctx, "Run")
	defer span.End()

	if id == uuid.Nil {
		return entity.Draw{}, ErrResourceInvalid
	}

	return uc.repo.Run(ctx, id)
}

const maxReasonLength = 500

// UpdatePrize edits the prize text; void draws are frozen.
func (uc *draw) UpdatePrize(ctx context.Context, id uuid.UUID, prize string) (entity.Draw, error) {
	ctx, span := uc.tracer.Start(ctx, "UpdatePrize")
	defer span.End()

	if id == uuid.Nil || prize == "" {
		return entity.Draw{}, ErrResourceInvalid
	}

	current, err := uc.repo.Get(ctx, id)
	if err != nil {
		return entity.Draw{}, err
	}

	if current.Status == entity.StatusVoid {
		return entity.Draw{}, ErrInvalidState
	}

	return uc.repo.UpdatePrize(ctx, id, prize)
}

// Void marks a draw void. Reason is REQUIRED (sensitive mutation) and is
// stored on the row and in the audit trail.
func (uc *draw) Void(ctx context.Context, id uuid.UUID, reason string) (entity.Draw, error) {
	ctx, span := uc.tracer.Start(ctx, "Void")
	defer span.End()

	reason = strings.TrimSpace(reason)
	if reason == "" {
		return entity.Draw{}, ErrReasonRequired
	}

	if len(reason) > maxReasonLength {
		return entity.Draw{}, ErrResourceInvalid
	}

	return uc.repo.Void(ctx, id, reason)
}

// Reassign directly changes a drawn draw's winner. Requires a reason (the
// handler writes it to the audit trail); the repo validates ticket ownership.
func (uc *draw) Reassign(ctx context.Context, id uuid.UUID, ticketID uuid.UUID, reason string) (entity.Draw, error) {
	ctx, span := uc.tracer.Start(ctx, "Reassign")
	defer span.End()

	if id == uuid.Nil || ticketID == uuid.Nil {
		return entity.Draw{}, ErrResourceInvalid
	}

	reason = strings.TrimSpace(reason)
	if reason == "" {
		return entity.Draw{}, ErrReasonRequired
	}

	if len(reason) > maxReasonLength {
		return entity.Draw{}, ErrResourceInvalid
	}

	return uc.repo.Reassign(ctx, id, ticketID)
}
