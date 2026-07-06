package dto

import (
	"application/internal/competition/entity"
)

// ContentUpsertReq is the admin write body for one site-copy value.
type ContentUpsertReq struct {
	Value string `json:"value"`
}

// ContentItemResp is one site-copy row.
type ContentItemResp struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	UpdatedAt string `json:"updated_at"`
}

// ContentResp is the public read shape: a key->value map (what the frontend
// consumes) plus the row list for the admin editor.
type ContentResp struct {
	Items map[string]string `json:"items"`
	Rows  []ContentItemResp `json:"rows"`
}

// ToContentResp maps site-content rows to the API shape.
func ToContentResp(rows []entity.SiteContent) ContentResp {
	resp := ContentResp{
		Items: make(map[string]string, len(rows)),
		Rows:  make([]ContentItemResp, 0, len(rows)),
	}

	for _, row := range rows {
		resp.Items[row.Key] = row.Value
		resp.Rows = append(resp.Rows, ContentItemResp{
			Key:       row.Key,
			Value:     row.Value,
			UpdatedAt: formatTime(row.UpdatedAt),
		})
	}

	return resp
}
