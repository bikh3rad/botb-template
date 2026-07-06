package dto

import (
	"application/internal/competition/entity"
)

// AuditItemResp is one audit-log entry in API shape.
type AuditItemResp struct {
	ID         string `json:"id"`
	ActorEmail string `json:"actor_email"`
	Action     string `json:"action"`
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
	Reason     string `json:"reason,omitempty"`
	CreatedAt  string `json:"created_at"`
}

// AuditListResp is the recent-audit envelope.
type AuditListResp struct {
	Count   int             `json:"count"`
	Entries []AuditItemResp `json:"entries"`
}

// ToAuditListResp maps audit entities to the API envelope.
func ToAuditListResp(entries []entity.AuditEntry) AuditListResp {
	items := make([]AuditItemResp, 0, len(entries))
	for _, e := range entries {
		items = append(items, AuditItemResp{
			ID:         e.ID,
			ActorEmail: e.ActorEmail,
			Action:     e.Action,
			EntityType: e.EntityType,
			EntityID:   e.EntityID,
			Reason:     e.Reason,
			CreatedAt:  formatTime(e.CreatedAt),
		})
	}

	return AuditListResp{Count: len(items), Entries: items}
}
