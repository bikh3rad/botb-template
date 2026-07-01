package dto

import (
	"application/internal/draw/entity"
	"time"
)

const rfc3339 = "2006-01-02T15:04:05Z07:00"

// CreateDrawReq is the admin request body for creating a pending draw.
type CreateDrawReq struct {
	CompetitionID string `json:"competition_id"`
	Prize         string `json:"prize"`
}

// DrawResp is the API representation of a draw.
type DrawResp struct {
	ID             string `json:"id"`
	CompetitionID  string `json:"competition_id"`
	WinnerUserID   string `json:"winner_user_id,omitempty"`
	WinnerTicketID string `json:"winner_ticket_id,omitempty"`
	Prize          string `json:"prize"`
	Status         string `json:"status"`
	DrawnAt        string `json:"drawn_at,omitempty"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// ToDrawResp maps an entity to its API shape.
func ToDrawResp(d entity.Draw) DrawResp {
	resp := DrawResp{
		ID:            d.ID.String(),
		CompetitionID: d.CompetitionID.String(),
		Prize:         d.Prize,
		Status:        string(d.Status),
		CreatedAt:     formatTime(d.CreatedAt),
		UpdatedAt:     formatTime(d.UpdatedAt),
	}

	if d.WinnerUserID != nil {
		resp.WinnerUserID = d.WinnerUserID.String()
	}

	if d.WinnerTicketID != nil {
		resp.WinnerTicketID = d.WinnerTicketID.String()
	}

	if d.DrawnAt != nil {
		resp.DrawnAt = formatTime(*d.DrawnAt)
	}

	return resp
}

// DrawListResp is the paginated list envelope.
type DrawListResp struct {
	Count  int        `json:"count"`
	Total  int        `json:"total"`
	Limit  int        `json:"limit"`
	Offset int        `json:"offset"`
	Draws  []DrawResp `json:"draws"`
}

// ToDrawListResp maps a page of draws to the list envelope.
func ToDrawListResp(draws []entity.Draw, total, limit, offset int) DrawListResp {
	items := make([]DrawResp, 0, len(draws))
	for i := range draws {
		items = append(items, ToDrawResp(draws[i]))
	}

	return DrawListResp{
		Count:  len(items),
		Total:  total,
		Limit:  limit,
		Offset: offset,
		Draws:  items,
	}
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format(rfc3339)
}
