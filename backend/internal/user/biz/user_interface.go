package biz

import (
	"application/internal/user/entity"
	"context"

	"github.com/google/uuid"
)

// UserListFilter narrows/paginates a user listing.
type UserListFilter struct {
	Query  string // matched against name/email (case-insensitive), empty = all
	Limit  int    // page size (defaulted/capped in the use case)
	Offset int
}

// UserPage is a page of users plus the total match count for pagination.
type UserPage struct {
	Users []entity.User
	Total int
}

// UsecaseUser is consumed by the HTTP handler.
type UsecaseUser interface {
	Register(ctx context.Context, name, email string) (entity.User, error)
	List(ctx context.Context, filter UserListFilter) (UserPage, error)
	Get(ctx context.Context, id uuid.UUID) (entity.User, error)
}

// RepositoryUser persists users. Implemented by internal/user/repo (pgx).
type RepositoryUser interface {
	Create(ctx context.Context, u entity.User) (entity.User, error)
	List(ctx context.Context, filter UserListFilter) (UserPage, error)
	Get(ctx context.Context, id uuid.UUID) (entity.User, error)
}
