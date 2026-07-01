package dto

import (
	"application/internal/competition/biz"
	"application/internal/competition/entity"
	"time"
)

const rfc3339 = "2006-01-02T15:04:05Z07:00"

// CompetitionReq is the create/update request body (admin).
type CompetitionReq struct {
	Title            string `json:"title"`
	Slug             string `json:"slug,omitempty"`
	Description      string `json:"description"`
	Prize            string `json:"prize"`
	TicketPricePence int64  `json:"ticket_price_pence"`
	TicketsTotal     int64  `json:"tickets_total"`
	Status           string `json:"status"`
	StartsAt         string `json:"starts_at"`
	EndsAt           string `json:"ends_at"`
}

// ToCreateInput maps the request to a biz.CreateInput, parsing timestamps.
func (r CompetitionReq) ToCreateInput() (biz.CreateInput, error) {
	starts, ends, err := parseWindow(r.StartsAt, r.EndsAt)
	if err != nil {
		return biz.CreateInput{}, err
	}

	return biz.CreateInput{
		Title:            r.Title,
		Slug:             r.Slug,
		Description:      r.Description,
		Prize:            r.Prize,
		TicketPricePence: r.TicketPricePence,
		TicketsTotal:     r.TicketsTotal,
		Status:           entity.Status(r.Status),
		StartsAt:         starts,
		EndsAt:           ends,
	}, nil
}

// ToUpdateInput maps the request to a biz.UpdateInput, parsing timestamps.
func (r CompetitionReq) ToUpdateInput() (biz.UpdateInput, error) {
	starts, ends, err := parseWindow(r.StartsAt, r.EndsAt)
	if err != nil {
		return biz.UpdateInput{}, err
	}

	return biz.UpdateInput{
		Title:            r.Title,
		Description:      r.Description,
		Prize:            r.Prize,
		TicketPricePence: r.TicketPricePence,
		TicketsTotal:     r.TicketsTotal,
		Status:           entity.Status(r.Status),
		StartsAt:         starts,
		EndsAt:           ends,
	}, nil
}

func parseWindow(startsAt, endsAt string) (starts, ends time.Time, err error) {
	if startsAt != "" {
		starts, err = time.Parse(rfc3339, startsAt)
		if err != nil {
			return time.Time{}, time.Time{}, biz.ErrResourceInvalid
		}
	}

	if endsAt != "" {
		ends, err = time.Parse(rfc3339, endsAt)
		if err != nil {
			return time.Time{}, time.Time{}, biz.ErrResourceInvalid
		}
	}

	return starts, ends, nil
}

// MediaRefResp is the API shape of an associated media object.
type MediaRefResp struct {
	ID          string `json:"id"`
	Kind        string `json:"kind"`
	Bucket      string `json:"bucket"`
	ObjectKey   string `json:"object_key"`
	ContentType string `json:"content_type"`
	Position    int    `json:"position"`
}

// CompetitionResp is the API representation of a competition, media included.
type CompetitionResp struct {
	ID               string         `json:"id"`
	Title            string         `json:"title"`
	Slug             string         `json:"slug"`
	Description      string         `json:"description"`
	Prize            string         `json:"prize"`
	TicketPricePence int64          `json:"ticket_price_pence"`
	TicketsTotal     int64          `json:"tickets_total"`
	TicketsSold      int64          `json:"tickets_sold"`
	Status           string         `json:"status"`
	StartsAt         string         `json:"starts_at"`
	EndsAt           string         `json:"ends_at"`
	CreatedAt        string         `json:"created_at"`
	UpdatedAt        string         `json:"updated_at"`
	Media            []MediaRefResp `json:"media"`
}

// ToCompetitionResp maps an entity to its API shape.
func ToCompetitionResp(c entity.Competition) CompetitionResp {
	media := make([]MediaRefResp, 0, len(c.Media))
	for _, m := range c.Media {
		media = append(media, MediaRefResp{
			ID:          m.ID.String(),
			Kind:        m.Kind,
			Bucket:      m.Bucket,
			ObjectKey:   m.ObjectKey,
			ContentType: m.ContentType,
			Position:    m.Position,
		})
	}

	return CompetitionResp{
		ID:               c.ID.String(),
		Title:            c.Title,
		Slug:             c.Slug,
		Description:      c.Description,
		Prize:            c.Prize,
		TicketPricePence: c.TicketPricePence,
		TicketsTotal:     c.TicketsTotal,
		TicketsSold:      c.TicketsSold,
		Status:           string(c.Status),
		StartsAt:         formatTime(c.StartsAt),
		EndsAt:           formatTime(c.EndsAt),
		CreatedAt:        formatTime(c.CreatedAt),
		UpdatedAt:        formatTime(c.UpdatedAt),
		Media:            media,
	}
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format(rfc3339)
}

// CompetitionListResp is the list envelope.
type CompetitionListResp struct {
	Count        int               `json:"count"`
	Competitions []CompetitionResp `json:"competitions"`
}

// ToCompetitionListResp maps a slice of entities to the list envelope.
func ToCompetitionListResp(cs []entity.Competition) CompetitionListResp {
	items := make([]CompetitionResp, 0, len(cs))
	for i := range cs {
		items = append(items, ToCompetitionResp(cs[i]))
	}

	return CompetitionListResp{Count: len(items), Competitions: items}
}
