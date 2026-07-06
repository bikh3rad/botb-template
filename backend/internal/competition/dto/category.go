package dto

import (
	"application/internal/competition/entity"
)

// CategoryReq is the create/update request body (admin).
type CategoryReq struct {
	Name string `json:"name"`
	Slug string `json:"slug,omitempty"`
}

// CategoryResp is the API representation of a category.
type CategoryResp struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	CreatedAt string `json:"created_at"`
}

// ToCategoryResp maps an entity to its API shape.
func ToCategoryResp(c entity.Category) CategoryResp {
	return CategoryResp{
		ID:        c.ID.String(),
		Name:      c.Name,
		Slug:      c.Slug,
		CreatedAt: formatTime(c.CreatedAt),
	}
}

// CategoryListResp is the list envelope.
type CategoryListResp struct {
	Count      int            `json:"count"`
	Categories []CategoryResp `json:"categories"`
}

// ToCategoryListResp maps categories to the list envelope.
func ToCategoryListResp(cs []entity.Category) CategoryListResp {
	items := make([]CategoryResp, 0, len(cs))
	for i := range cs {
		items = append(items, ToCategoryResp(cs[i]))
	}

	return CategoryListResp{Count: len(items), Categories: items}
}
