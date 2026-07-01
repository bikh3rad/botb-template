package dto

import (
	"application/internal/user/entity"
	"time"
)

const rfc3339 = "2006-01-02T15:04:05Z07:00"

// RegisterReq is the public registration request body.
type RegisterReq struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserResp is the API representation of a user.
type UserResp struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	TicketsOwned    int64  `json:"tickets_owned"`
	TotalSpentPence int64  `json:"total_spent_pence"`
	CreatedAt       string `json:"created_at"`
}

// ToUserResp maps an entity to its API shape.
func ToUserResp(u entity.User) UserResp {
	return UserResp{
		ID:              u.ID.String(),
		Name:            u.Name,
		Email:           u.Email,
		TicketsOwned:    u.TicketsOwned,
		TotalSpentPence: u.TotalSpentPence,
		CreatedAt:       formatTime(u.CreatedAt),
	}
}

// UserListResp is the paginated list envelope.
type UserListResp struct {
	Count  int        `json:"count"`
	Total  int        `json:"total"`
	Limit  int        `json:"limit"`
	Offset int        `json:"offset"`
	Users  []UserResp `json:"users"`
}

// ToUserListResp maps a page of users to the list envelope.
func ToUserListResp(users []entity.User, total, limit, offset int) UserListResp {
	items := make([]UserResp, 0, len(users))
	for i := range users {
		items = append(items, ToUserResp(users[i]))
	}

	return UserListResp{
		Count:  len(items),
		Total:  total,
		Limit:  limit,
		Offset: offset,
		Users:  items,
	}
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format(rfc3339)
}
