//go:build e2e
// +build e2e

package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

// TestCompetitionFullLifecycleAndSafeDelete drives a competition end to end:
// create (all fields + one image) -> visible publicly -> edit (prize + category
// + a second image reordered) -> new prize visible publicly -> a paying entrant
// makes DELETE return 409 -> close instead -> a fresh draft deletes cleanly
// (gone from the public list + DB, its media object gone from MinIO), and the
// audit trail recorded the mutations.
func TestCompetitionFullLifecycleAndSafeDelete(t *testing.T) {
	requireStack(t)

	token := superadminToken(t)

	// A category to attach.
	var cats categoryList
	_, catBody := request(t, http.MethodGet, "/apis/competition/v1/categories", "", nil)
	decode(t, catBody, &cats)
	if len(cats.Categories) < 2 {
		t.Fatal("need at least two seeded categories")
	}

	// 1. Create a draft competition with every field + one image.
	slug := fmt.Sprintf("lifecycle-%d", time.Now().UnixNano())
	_, createBody := request(t, http.MethodPost, "/apis/competition/v1/admin/competitions", token,
		map[string]any{
			"title": "Lifecycle Comp", "slug": slug, "description": "desc",
			"prize": "Original Prize", "ticket_price_pence": 100, "tickets_total": 1000,
			"category_id": cats.Categories[0].ID, "status": "draft",
		})

	var comp competition
	decode(t, createBody, &comp)

	first := uploadMedia(t, token, "competition", comp.ID, 0)

	// 2. Publish (draft -> live) and confirm it appears in the PUBLIC live list.
	publishAndAssertLive(t, token, comp.ID, slug, cats.Categories[0].ID)

	// 3. Edit prize + category + add a second image and reorder it to the front.
	second := uploadMedia(t, token, "competition", comp.ID, 1)
	_, _ = request(t, http.MethodPut, "/apis/media/v1/admin/media/"+second.ID, token,
		map[string]any{"position": 0})

	_, editBody := request(t, http.MethodPut, "/apis/competition/v1/admin/competitions/"+comp.ID, token,
		map[string]any{
			"title": "Lifecycle Comp", "slug": slug, "description": "desc",
			"prize": "Updated Prize!", "ticket_price_pence": 100, "tickets_total": 1000,
			"category_id": cats.Categories[1].ID, "status": "live",
		})
	var edited competition
	decode(t, editBody, &edited)
	if edited.Prize != "Updated Prize!" || edited.CategoryID != cats.Categories[1].ID {
		t.Fatalf("edit did not round-trip: %+v", edited)
	}

	// New prize is visible publicly.
	if pub := publicCompetition(t, comp.ID); pub.Prize != "Updated Prize!" {
		t.Errorf("public prize: got %q want %q", pub.Prize, "Updated Prize!")
	}

	// 4. A paying entrant makes DELETE return 409.
	buyerEmail := uniqueEmail()
	_, regBody := request(t, http.MethodPost, "/apis/user/v1/users", "",
		registerReq{Name: "Lifecycle Buyer", Email: buyerEmail})
	var buyer user
	decode(t, regBody, &buyer)

	buyResp, buyBody := request(t, http.MethodPost, "/apis/user/v1/tickets", "",
		purchaseReq{CompetitionID: comp.ID, UserID: buyer.ID, Quantity: 1})
	requireStatus(t, buyResp, buyBody, http.StatusCreated)

	delResp, delBody := request(t, http.MethodDelete, "/apis/competition/v1/admin/competitions/"+comp.ID, token, nil)
	requireStatus(t, delResp, delBody, http.StatusConflict)

	// 5. Close instead (live -> closed) — the safe alternative.
	_, closeBody := request(t, http.MethodPut, "/apis/competition/v1/admin/competitions/"+comp.ID, token,
		map[string]any{
			"title": "Lifecycle Comp", "slug": slug, "description": "desc",
			"prize": "Updated Prize!", "ticket_price_pence": 100, "tickets_total": 1000,
			"category_id": cats.Categories[1].ID, "status": "closed",
		})
	var closed competition
	decode(t, closeBody, &closed)
	if closed.Status != "closed" {
		t.Fatalf("expected closed, got %q", closed.Status)
	}

	// Media cleanup for the closed (undeletable) comp is out of band; delete via
	// the media endpoint so we don't leak the two objects from this test.
	_, _ = request(t, http.MethodDelete, "/apis/media/v1/admin/media/"+first.ID, token, nil)
	_, _ = request(t, http.MethodDelete, "/apis/media/v1/admin/media/"+second.ID, token, nil)

	// 6. A fresh DRAFT with an image deletes cleanly and purges its media object.
	freshSlug := fmt.Sprintf("lifecycle-fresh-%d", time.Now().UnixNano())
	_, freshBody := request(t, http.MethodPost, "/apis/competition/v1/admin/competitions", token,
		map[string]any{
			"title": "Fresh Draft", "slug": freshSlug, "prize": "P",
			"ticket_price_pence": 100, "tickets_total": 10, "status": "draft",
		})
	var fresh competition
	decode(t, freshBody, &fresh)

	freshMedia := uploadMedia(t, token, "competition", fresh.ID, 0)
	objURL := mediaBaseURL() + "/" + freshMedia.Bucket + "/" + freshMedia.ObjectKey
	if st := rawStatus(t, objURL); st < 200 || st >= 300 {
		t.Fatalf("fresh media should be reachable, got %d", st)
	}

	freshDel, freshDelBody := request(t, http.MethodDelete, "/apis/competition/v1/admin/competitions/"+fresh.ID, token, nil)
	requireStatus(t, freshDel, freshDelBody, http.StatusNoContent)

	// Gone from the public single-get and its media object purged from MinIO.
	getResp, _ := request(t, http.MethodGet, "/apis/competition/v1/competitions/"+fresh.ID, "", nil)
	if getResp.StatusCode != http.StatusNotFound {
		t.Errorf("deleted competition public get: got %d want 404", getResp.StatusCode)
	}

	gone := false
	for range 10 {
		if rawStatus(t, objURL) == http.StatusNotFound {
			gone = true

			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	if !gone {
		t.Errorf("expected competition-delete to purge its media object from MinIO at %s", objURL)
	}

	// Audit trail recorded the lifecycle mutations.
	assertAuditRow(t, "competition", comp.ID, "competition.create")
	assertAuditRow(t, "competition", comp.ID, "competition.update")
	assertAuditRow(t, "competition", fresh.ID, "competition.delete")
}

// TestDrawDoubleRunConflict: a second run of an already-drawn draw conflicts.
func TestDrawDoubleRunConflict(t *testing.T) {
	requireStack(t)

	token := superadminToken(t)

	// Fresh competition, made live, with one ticket so the draw has a winner.
	slug := fmt.Sprintf("double-run-%d", time.Now().UnixNano())
	_, cBody := request(t, http.MethodPost, "/apis/competition/v1/admin/competitions", token,
		map[string]any{
			"title": "Double Run", "slug": slug, "prize": "P",
			"ticket_price_pence": 100, "tickets_total": 10, "status": "draft",
		})
	var comp competition
	decode(t, cBody, &comp)

	_, _ = request(t, http.MethodPut, "/apis/competition/v1/admin/competitions/"+comp.ID, token,
		map[string]any{
			"title": "Double Run", "slug": slug, "prize": "P",
			"ticket_price_pence": 100, "tickets_total": 10, "status": "live",
		})

	_, rBody := request(t, http.MethodPost, "/apis/user/v1/users", "",
		registerReq{Name: "DR Buyer", Email: uniqueEmail()})
	var buyer user
	decode(t, rBody, &buyer)
	_, _ = request(t, http.MethodPost, "/apis/user/v1/tickets", "",
		purchaseReq{CompetitionID: comp.ID, UserID: buyer.ID, Quantity: 1})

	_, dBody := request(t, http.MethodPost, "/apis/draw/v1/admin/draws", token,
		map[string]any{"competition_id": comp.ID, "prize": "DR Prize"})
	var d draw
	decode(t, dBody, &d)

	run1, run1Body := request(t, http.MethodPost, "/apis/draw/v1/admin/draws/"+d.ID+"/run", token, nil)
	requireStatus(t, run1, run1Body, http.StatusOK)

	// Second run of an already-drawn draw → 409.
	run2, run2Body := request(t, http.MethodPost, "/apis/draw/v1/admin/draws/"+d.ID+"/run", token, nil)
	requireStatus(t, run2, run2Body, http.StatusConflict)

	// Clean up so this test leaves the seeded live count unchanged: void the
	// draw and close the competition (it has an entrant, so it can't be
	// deleted — closing removes it from the live grid, which is the point).
	_, _ = request(t, http.MethodPost, "/apis/draw/v1/admin/draws/"+d.ID+"/void", token,
		map[string]any{"reason": "double-run test cleanup"})
	_, _ = request(t, http.MethodPut, "/apis/competition/v1/admin/competitions/"+comp.ID, token,
		map[string]any{
			"title": "Double Run", "slug": slug, "prize": "P",
			"ticket_price_pence": 100, "tickets_total": 10, "status": "closed",
		})
}

// TestSiteContentReflectsPublicly: an admin content edit is visible on the
// public content endpoint the site renders from.
func TestSiteContentReflectsPublicly(t *testing.T) {
	requireStack(t)

	token := superadminToken(t)
	marker := fmt.Sprintf("%d winners", time.Now().UnixNano()%100000)

	_, _ = request(t, http.MethodPut, "/apis/competition/v1/admin/content/winners.count", token,
		map[string]any{"value": marker})

	_, body := request(t, http.MethodGet, "/apis/competition/v1/content", "", nil)

	var content contentResp
	decode(t, body, &content)
	if content.Items["winners.count"] != marker {
		t.Errorf("public content: got %q want %q", content.Items["winners.count"], marker)
	}

	// Restore the seeded value.
	_, _ = request(t, http.MethodPut, "/apis/competition/v1/admin/content/winners.count", token,
		map[string]any{"value": "9,700 winners"})
}

// publishAndAssertLive flips a draft to live and asserts it is in the public
// live list.
func publishAndAssertLive(t *testing.T, token, id, slug, categoryID string) {
	t.Helper()

	_, _ = request(t, http.MethodPut, "/apis/competition/v1/admin/competitions/"+id, token,
		map[string]any{
			"title": "Lifecycle Comp", "slug": slug, "description": "desc",
			"prize": "Original Prize", "ticket_price_pence": 100, "tickets_total": 1000,
			"category_id": categoryID, "status": "live",
		})

	_, body := request(t, http.MethodGet, "/apis/competition/v1/competitions?status=live", "", nil)

	var list competitionList
	decode(t, body, &list)
	for _, c := range list.Competitions {
		if c.ID == id {
			return
		}
	}

	t.Fatalf("published competition %s not found in public live list", id)
}

// publicCompetition fetches a competition through the public single-get.
func publicCompetition(t *testing.T, id string) competition {
	t.Helper()

	_, body := request(t, http.MethodGet, "/apis/competition/v1/competitions/"+id, "", nil)

	var c competition
	decode(t, body, &c)

	return c
}
