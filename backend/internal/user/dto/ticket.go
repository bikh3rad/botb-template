package dto

import (
	"application/internal/user/entity"
)

// PurchaseReq is the ticket purchase request body.
type PurchaseReq struct {
	CompetitionID string `json:"competition_id"`
	UserID        string `json:"user_id"`
	Quantity      int    `json:"quantity"`
}

// TicketResp is the API representation of a ticket.
type TicketResp struct {
	ID            string `json:"id"`
	CompetitionID string `json:"competition_id"`
	UserID        string `json:"user_id"`
	PurchasedAt   string `json:"purchased_at"`
}

// ToTicketResp maps an entity to its API shape.
func ToTicketResp(t entity.Ticket) TicketResp {
	return TicketResp{
		ID:            t.ID.String(),
		CompetitionID: t.CompetitionID.String(),
		UserID:        t.UserID.String(),
		PurchasedAt:   formatTime(t.PurchasedAt),
	}
}

func toTicketResps(ts []entity.Ticket) []TicketResp {
	items := make([]TicketResp, 0, len(ts))
	for i := range ts {
		items = append(items, ToTicketResp(ts[i]))
	}

	return items
}

// PurchaseResp is returned after a successful purchase.
type PurchaseResp struct {
	Tickets        []TicketResp `json:"tickets"`
	User           UserResp     `json:"user"`
	TotalCostPence int64        `json:"total_cost_pence"`
	Count          int          `json:"count"`
}

// TicketListResp is the list envelope for a user's tickets.
type TicketListResp struct {
	Count   int          `json:"count"`
	Tickets []TicketResp `json:"tickets"`
}

// ToTicketListResp maps a slice of tickets to the list envelope.
func ToTicketListResp(ts []entity.Ticket) TicketListResp {
	items := toTicketResps(ts)

	return TicketListResp{Count: len(items), Tickets: items}
}
