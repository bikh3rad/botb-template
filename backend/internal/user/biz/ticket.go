package biz

import (
	"application/internal/user/entity"
	"context"
	"log/slog"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// maxQuantity caps a single purchase to keep the transaction bounded.
const maxQuantity = 1000

type ticket struct {
	logger *slog.Logger
	tracer trace.Tracer
	repo   RepositoryTicket
}

var _ UsecaseTicket = (*ticket)(nil)

// NewTicket constructs the ticket use case.
func NewTicket(logger *slog.Logger, repo RepositoryTicket) *ticket {
	return &ticket{
		logger: logger.With("layer", "Ticket"),
		tracer: otel.Tracer("TicketUseCase"),
		repo:   repo,
	}
}

func (uc *ticket) Purchase(ctx context.Context, in PurchaseInput) (PurchaseResult, error) {
	ctx, span := uc.tracer.Start(ctx, "Purchase")
	defer span.End()

	if in.CompetitionID == uuid.Nil || in.UserID == uuid.Nil {
		return PurchaseResult{}, ErrResourceInvalid
	}

	if in.Quantity <= 0 || in.Quantity > maxQuantity {
		return PurchaseResult{}, ErrResourceInvalid
	}

	return uc.repo.Purchase(ctx, in)
}

func (uc *ticket) ListByUser(ctx context.Context, userID uuid.UUID) ([]entity.Ticket, error) {
	if userID == uuid.Nil {
		return nil, ErrResourceInvalid
	}

	return uc.repo.ListByUser(ctx, userID)
}
