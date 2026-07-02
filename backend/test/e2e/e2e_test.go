//go:build e2e
// +build e2e

package e2e

import (
	"net/http"
	"strings"
	"testing"

	"application/internal/seeddata"
)

// TestPublicReads verifies the public competition endpoints return the seeded
// inventory through the gateway: the live count matches the canonical seed set
// and known titles are present, then a single competition resolves with media.
func TestPublicReads(t *testing.T) {
	requireStack(t)

	var list competitionList

	t.Run("list live competitions matches seed", func(t *testing.T) {
		resp, body := request(t, http.MethodGet, "/apis/competition/v1/competitions?status=live", "", nil)
		requireStatus(t, resp, body, http.StatusOK)
		decode(t, body, &list)

		if list.Count != len(seeddata.Competitions) {
			t.Fatalf("live competition count: got %d want %d (len seeddata.Competitions)",
				list.Count, len(seeddata.Competitions))
		}
		if len(list.Competitions) != list.Count {
			t.Fatalf("count field %d disagrees with competitions length %d", list.Count, len(list.Competitions))
		}

		// A few known titles from seeddata must be present.
		for _, want := range []string{"£1.2M HOME IN ZONE 1", "AUDI R8 FOR 21P!"} {
			if !containsTitle(list.Competitions, want) {
				t.Errorf("expected seeded title %q in live competitions", want)
			}
		}
	})

	t.Run("get single competition with media", func(t *testing.T) {
		comp := firstWithMedia(t, list.Competitions)

		resp, body := request(t, http.MethodGet, "/apis/competition/v1/competitions/"+comp.ID, "", nil)
		requireStatus(t, resp, body, http.StatusOK)

		var got competition
		decode(t, body, &got)

		if got.ID != comp.ID {
			t.Fatalf("competition id: got %q want %q", got.ID, comp.ID)
		}
		if got.Title == "" || got.Slug == "" || got.Status == "" {
			t.Errorf("expected populated fields, got %+v", got)
		}
		if len(got.Media) == 0 {
			t.Fatalf("expected non-empty media for competition %s", got.ID)
		}
	})
}

// TestMediaResolvesFromMinIO proves a competition's media resolves to real image
// bytes out of MinIO. The media service ALSO issues a time-limited presigned URL
// (asserted non-empty), but that URL is signed for the service's configured
// endpoint — the internal `minio:9000` host — so it only resolves from inside
// the compose network. We therefore fetch the bytes the way the BROWSER does:
// via the public-read bucket URL on the host (http://localhost:9000/<bucket>/<key>),
// which is exactly what the frontend builds from the media record.
func TestMediaResolvesFromMinIO(t *testing.T) {
	requireStack(t)

	comp := fetchCompetitionWithMedia(t)
	mediaID := comp.Media[0].ID

	resp, body := request(t, http.MethodGet, "/apis/media/v1/media/"+mediaID, "", nil)
	requireStatus(t, resp, body, http.StatusOK)

	var m media
	decode(t, body, &m)

	if m.ID != mediaID {
		t.Fatalf("media id: got %q want %q", m.ID, mediaID)
	}
	if m.URL == "" {
		t.Fatalf("expected non-empty presigned url for media %s", mediaID)
	}

	publicURL := mediaBaseURL() + "/" + m.Bucket + "/" + m.ObjectKey
	objResp, objBody := request(t, http.MethodGet, publicURL, "", nil)
	if objResp.StatusCode < 200 || objResp.StatusCode >= 300 {
		t.Fatalf("public object GET %s status: got %d want 2xx; body=%s", publicURL, objResp.StatusCode, string(objBody))
	}

	if ct := objResp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "image/") {
		t.Errorf("public object Content-Type: got %q want image/*", ct)
	}
}

// TestDrawPublicReadReturnsWinners discovers a drawn draw via the admin list,
// then reads it through the PUBLIC endpoint and asserts a real seeded winner.
func TestDrawPublicReadReturnsWinners(t *testing.T) {
	requireStack(t)

	token := mintAdminToken(t)

	resp, body := request(t, http.MethodGet, "/apis/draw/v1/admin/draws", token, nil)
	requireStatus(t, resp, body, http.StatusOK)

	var list drawList
	decode(t, body, &list)

	drawn := firstDrawn(list.Draws)
	if drawn == nil {
		t.Fatalf("expected at least one drawn draw in admin list; got %d draws", len(list.Draws))
	}

	// PUBLIC read (no token) of the drawn draw.
	pubResp, pubBody := request(t, http.MethodGet, "/apis/draw/v1/draws/"+drawn.ID, "", nil)
	requireStatus(t, pubResp, pubBody, http.StatusOK)

	var got draw
	decode(t, pubBody, &got)

	if got.Status != "drawn" {
		t.Errorf("draw status: got %q want %q", got.Status, "drawn")
	}
	if got.WinnerUserID == "" {
		t.Errorf("expected non-empty winner_user_id for drawn draw %s", got.ID)
	}
	if !seededWinnerPrize(got.Prize) {
		t.Errorf("draw prize %q is not one of the seeded winner prizes", got.Prize)
	}
}

// TestAdminAuthGuard proves the two-layer admin guard: no token is rejected with
// 401, and a validly minted token passes the guard (any non-401 status).
func TestAdminAuthGuard(t *testing.T) {
	requireStack(t)

	t.Run("no token on admin draws is 401", func(t *testing.T) {
		resp, body := request(t, http.MethodGet, "/apis/draw/v1/admin/draws", "", nil)
		requireStatus(t, resp, body, http.StatusUnauthorized)
	})

	t.Run("no token on admin competition create is 401", func(t *testing.T) {
		resp, body := request(t, http.MethodPost, "/apis/competition/v1/admin/competitions", "",
			map[string]any{"title": "guard-check"})
		requireStatus(t, resp, body, http.StatusUnauthorized)
	})

	t.Run("valid token passes the guard", func(t *testing.T) {
		token := mintAdminToken(t)

		// A minted token must get PAST auth. We do not care whether the create
		// then succeeds or fails validation — only that it is not a 401.
		resp, body := request(t, http.MethodPost, "/apis/competition/v1/admin/competitions", token,
			map[string]any{"title": "guard-check"})
		if resp.StatusCode == http.StatusUnauthorized {
			t.Fatalf("valid token should pass auth, got 401; body=%s", string(body))
		}
	})
}

// TestGatewayRouting checks the gateway's own responses: unknown services 404
// with a clear message, and the friendly root identifies the gateway.
func TestGatewayRouting(t *testing.T) {
	requireStack(t)

	t.Run("unknown service is 404", func(t *testing.T) {
		resp, body := request(t, http.MethodGet, "/apis/nope/v1/x", "", nil)
		requireStatus(t, resp, body, http.StatusNotFound)
		if !strings.Contains(string(body), "unknown service") {
			t.Errorf("expected body to contain %q, got %s", "unknown service", string(body))
		}
	})

	t.Run("root identifies the gateway", func(t *testing.T) {
		resp, body := request(t, http.MethodGet, "/", "", nil)
		requireStatus(t, resp, body, http.StatusOK)
		if !strings.Contains(string(body), "botb-gateway") {
			t.Errorf("expected root body to contain %q, got %s", "botb-gateway", string(body))
		}
	})
}

// TestFullPublicFlow drives the real transactional path against Postgres:
// register a user, purchase a ticket for a live competition, then read the user
// back and assert the balance moved. THIS covers the write path the unit tests
// could only exercise against the in-memory ramsql fake.
func TestFullPublicFlow(t *testing.T) {
	requireStack(t)

	// Pick a live competition with a non-zero ticket price so total_spent_pence
	// is a meaningful assertion (several seeded comps are free at 0 pence).
	comp := firstPaidLiveCompetition(t)
	const qty = 1

	// 1. Register a brand-new user with a per-run unique email.
	email := uniqueEmail()
	regResp, regBody := request(t, http.MethodPost, "/apis/user/v1/users", "",
		registerReq{Name: "E2E Buyer", Email: email})
	requireStatus(t, regResp, regBody, http.StatusCreated)

	var created user
	decode(t, regBody, &created)
	if created.ID == "" {
		t.Fatalf("expected a user id from registration; body=%s", string(regBody))
	}

	// 2. Purchase a ticket for the chosen competition.
	buyResp, buyBody := request(t, http.MethodPost, "/apis/user/v1/tickets", "",
		purchaseReq{CompetitionID: comp.ID, UserID: created.ID, Quantity: qty})
	requireStatus(t, buyResp, buyBody, http.StatusCreated)

	var purchase purchaseResp
	decode(t, buyBody, &purchase)
	if purchase.Count != qty {
		t.Errorf("purchase count: got %d want %d", purchase.Count, qty)
	}
	if want := comp.TicketPricePence * qty; purchase.TotalCostPence != want {
		t.Errorf("purchase total_cost_pence: got %d want %d", purchase.TotalCostPence, want)
	}

	// 3. Read the user back via the admin endpoint and assert the balance moved.
	token := mintAdminToken(t)
	getResp, getBody := request(t, http.MethodGet, "/apis/user/v1/admin/users/"+created.ID, token, nil)
	requireStatus(t, getResp, getBody, http.StatusOK)

	var got user
	decode(t, getBody, &got)
	if got.TicketsOwned != int64(qty) {
		t.Errorf("tickets_owned: got %d want %d", got.TicketsOwned, qty)
	}
	if want := comp.TicketPricePence * int64(qty); got.TotalSpentPence != want {
		t.Errorf("total_spent_pence: got %d want %d (price %d * qty %d)",
			got.TotalSpentPence, want, comp.TicketPricePence, qty)
	}
}

// TestIdempotentSeed asserts the seed did not double-insert: the live count
// equals the canonical seed length and every live title is unique. Re-running
// `make seed` (or `docker compose up`) is idempotent because the seeder derives
// deterministic UUIDv5 IDs from natural keys and upserts — it never duplicates —
// so this count stays stable across re-seeds. We assert the invariant rather
// than re-running the seeder (which needs its own DB/MinIO env).
func TestIdempotentSeed(t *testing.T) {
	requireStack(t)

	resp, body := request(t, http.MethodGet, "/apis/competition/v1/competitions?status=live", "", nil)
	requireStatus(t, resp, body, http.StatusOK)

	var list competitionList
	decode(t, body, &list)

	if list.Count != len(seeddata.Competitions) {
		t.Fatalf("live count: got %d want %d — a double-insert would inflate this",
			list.Count, len(seeddata.Competitions))
	}

	seen := make(map[string]struct{}, len(list.Competitions))
	for _, c := range list.Competitions {
		if _, dup := seen[c.Title]; dup {
			t.Errorf("duplicate live competition title %q — seeder inserted twice", c.Title)
		}
		seen[c.Title] = struct{}{}
	}
}

// --- small local helpers over the fetched data ---

func containsTitle(comps []competition, title string) bool {
	for _, c := range comps {
		if c.Title == title {
			return true
		}
	}

	return false
}

func firstWithMedia(t *testing.T, comps []competition) competition {
	t.Helper()

	for _, c := range comps {
		if len(c.Media) > 0 {
			return c
		}
	}

	// The list envelope may not populate media; fall back to a single-resource
	// fetch which always does.
	return fetchCompetitionWithMedia(t)
}

// fetchCompetitionWithMedia fetches the live list, then does single-resource
// GETs (which always populate media) until it finds one with media.
func fetchCompetitionWithMedia(t *testing.T) competition {
	t.Helper()

	resp, body := request(t, http.MethodGet, "/apis/competition/v1/competitions?status=live", "", nil)
	requireStatus(t, resp, body, http.StatusOK)

	var list competitionList
	decode(t, body, &list)

	for _, c := range list.Competitions {
		cResp, cBody := request(t, http.MethodGet, "/apis/competition/v1/competitions/"+c.ID, "", nil)
		requireStatus(t, cResp, cBody, http.StatusOK)

		var full competition
		decode(t, cBody, &full)
		if len(full.Media) > 0 {
			return full
		}
	}

	t.Fatalf("no live competition with media found among %d competitions", len(list.Competitions))

	return competition{}
}

func firstPaidLiveCompetition(t *testing.T) competition {
	t.Helper()

	resp, body := request(t, http.MethodGet, "/apis/competition/v1/competitions?status=live", "", nil)
	requireStatus(t, resp, body, http.StatusOK)

	var list competitionList
	decode(t, body, &list)

	for _, c := range list.Competitions {
		if c.TicketPricePence > 0 {
			return c
		}
	}

	t.Fatalf("no live competition with a non-zero ticket price found")

	return competition{}
}

func firstDrawn(draws []draw) *draw {
	for i := range draws {
		if draws[i].Status == "drawn" && draws[i].WinnerUserID != "" {
			return &draws[i]
		}
	}

	return nil
}

// seededWinnerPrize reports whether prize matches any seeded winner's prize.
func seededWinnerPrize(prize string) bool {
	for _, w := range seeddata.Winners {
		if w.Prize == prize {
			return true
		}
	}

	return false
}
