package biz

import (
	"context"
	"time"

	"application/internal/competition/entity"

	"github.com/google/uuid"
)

// ListFilter narrows a competition listing. A nil field means "no filter".
type ListFilter struct {
	Status *entity.Status
}

// CreateInput is the biz-level payload for creating a competition. Slug is
// derived from Title when empty.
type CreateInput struct {
	Title            string
	Slug             string
	Description      string
	Prize            string
	TicketPricePence int64
	TicketsTotal     int64
	CategoryID       *uuid.UUID
	Status           entity.Status
	StartsAt         time.Time
	EndsAt           time.Time
}

// UpdateInput is the biz-level payload for a full competition update. It covers
// EVERY editable field (title, slug, description, prize, price, total,
// category, status, window). tickets_sold is deliberately absent — it is a
// derived value no endpoint may write.
type UpdateInput struct {
	Title            string
	Slug             string
	Description      string
	Prize            string
	TicketPricePence int64
	TicketsTotal     int64
	CategoryID       *uuid.UUID
	Status           entity.Status
	StartsAt         time.Time
	EndsAt           time.Time
}

// UsecaseCompetition is consumed by the HTTP handler.
type UsecaseCompetition interface {
	List(ctx context.Context, filter ListFilter) ([]entity.Competition, error)
	Get(ctx context.Context, id uuid.UUID) (entity.Competition, error)
	Create(ctx context.Context, in CreateInput) (entity.Competition, error)
	Update(ctx context.Context, id uuid.UUID, in UpdateInput) (entity.Competition, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// Repository persists competitions and resolves their associated media from the
// shared database. Implemented by internal/competition/repo (pgx).
type Repository interface {
	List(ctx context.Context, filter ListFilter) ([]entity.Competition, error)
	Get(ctx context.Context, id uuid.UUID) (entity.Competition, error)
	Create(ctx context.Context, c entity.Competition) (entity.Competition, error)
	Update(ctx context.Context, c entity.Competition) (entity.Competition, error)
	// Delete enforces the no-entrants rule (ErrCompetitionHasEntrants when the
	// competition has sold tickets or any draw), removes the competition and its
	// media rows in one transaction, and returns the removed media object keys
	// so the caller can purge them from object storage.
	Delete(ctx context.Context, id uuid.UUID) ([]string, error)
}
