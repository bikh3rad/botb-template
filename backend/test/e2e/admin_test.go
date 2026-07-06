//go:build e2e
// +build e2e

package e2e

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// A 1x1 PNG (67 bytes) — enough to prove the upload → MinIO → delete cycle.
var onePixelPNG = []byte{
	0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d,
	0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4, 0x89, 0x00, 0x00, 0x00,
	0x0a, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9c, 0x63, 0x00, 0x01, 0x00, 0x00,
	0x05, 0x00, 0x01, 0x0d, 0x0a, 0x2d, 0xb4, 0x00, 0x00, 0x00, 0x00, 0x49,
	0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
}

// TestAdminCompetitionFullEditAndMediaCycle drives a competition through a
// full-field admin edit and a media upload → reorder → delete cycle, verifying
// the deleted object is really gone from MinIO (browser-facing bucket URL).
func TestAdminCompetitionFullEditAndMediaCycle(t *testing.T) {
	requireStack(t)

	token := superadminToken(t)

	// Grab a category to attach.
	catResp, catBody := request(t, http.MethodGet, "/apis/competition/v1/categories", "", nil)
	requireStatus(t, catResp, catBody, http.StatusOK)

	var cats categoryList
	decode(t, catBody, &cats)
	if len(cats.Categories) == 0 {
		t.Fatal("expected seeded categories")
	}

	// Create a fresh draft competition to edit (avoids mutating seeded rows).
	slug := fmt.Sprintf("e2e-edit-%d", time.Now().UnixNano())
	createResp, createBody := request(t, http.MethodPost, "/apis/competition/v1/admin/competitions", token,
		map[string]any{
			"title":              "E2E Edit Target",
			"slug":               slug,
			"description":        "initial",
			"prize":              "A prize",
			"ticket_price_pence": 100,
			"tickets_total":      500,
			"status":             "draft",
		})
	requireStatus(t, createResp, createBody, http.StatusCreated)

	var comp competition
	decode(t, createBody, &comp)

	// Full-field update: change every editable field, move draft -> live, attach
	// a category. tickets_sold must remain 0 (derived, never written).
	newSlug := slug + "-v2"
	updResp, updBody := request(t, http.MethodPut, "/apis/competition/v1/admin/competitions/"+comp.ID, token,
		map[string]any{
			"title":              "E2E Edited Title",
			"slug":               newSlug,
			"description":        "updated description",
			"prize":              "Updated prize",
			"ticket_price_pence": 250,
			"tickets_total":      1000,
			"category_id":        cats.Categories[0].ID,
			"status":             "live",
		})
	requireStatus(t, updResp, updBody, http.StatusOK)

	var updated competition
	decode(t, updBody, &updated)

	if updated.Title != "E2E Edited Title" || updated.Slug != newSlug ||
		updated.Description != "updated description" || updated.Prize != "Updated prize" ||
		updated.TicketPricePence != 250 || updated.TicketsTotal != 1000 ||
		updated.Status != "live" || updated.CategoryID != cats.Categories[0].ID {
		t.Fatalf("full-field update did not round-trip: %+v", updated)
	}
	if updated.TicketsSold != 0 {
		t.Errorf("tickets_sold must stay derived/0, got %d", updated.TicketsSold)
	}

	// Illegal transition: live -> draft must be rejected (422).
	badResp, badBody := request(t, http.MethodPut, "/apis/competition/v1/admin/competitions/"+comp.ID, token,
		map[string]any{
			"title": "x", "slug": newSlug, "prize": "p",
			"ticket_price_pence": 250, "tickets_total": 1000, "status": "draft",
		})
	requireStatus(t, badResp, badBody, http.StatusUnprocessableEntity)

	// --- Media cycle: upload two images, reorder, delete one, verify in MinIO ---
	first := uploadMedia(t, token, "competition", comp.ID, 0)
	second := uploadMedia(t, token, "competition", comp.ID, 1)

	// Reorder: move `second` to position 0.
	reoResp, reoBody := request(t, http.MethodPut, "/apis/media/v1/admin/media/"+second.ID, token,
		map[string]any{"position": 0})
	requireStatus(t, reoResp, reoBody, http.StatusOK)

	// The object is reachable in MinIO before deletion.
	objURL := mediaBaseURL() + "/" + first.Bucket + "/" + first.ObjectKey
	if st := rawStatus(t, objURL); st < 200 || st >= 300 {
		t.Fatalf("uploaded object should be reachable, got %d at %s", st, objURL)
	}

	// Delete the first media; the DB row and the MinIO object both go.
	delResp, delBody := request(t, http.MethodDelete, "/apis/media/v1/admin/media/"+first.ID, token, nil)
	requireStatus(t, delResp, delBody, http.StatusNoContent)

	// The DB record is gone (404 via the public get).
	getResp, getBody := request(t, http.MethodGet, "/apis/media/v1/media/"+first.ID, "", nil)
	requireStatus(t, getResp, getBody, http.StatusNotFound)

	// The MinIO object is gone too — poll briefly for eventual consistency.
	gone := false
	for range 10 {
		if rawStatus(t, objURL) == http.StatusNotFound {
			gone = true

			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	if !gone {
		t.Errorf("expected deleted object to be removed from MinIO at %s", objURL)
	}

	// Clean up: delete the remaining media + the competition.
	request(t, http.MethodDelete, "/apis/media/v1/admin/media/"+second.ID, token, nil)
	request(t, http.MethodDelete, "/apis/competition/v1/admin/competitions/"+comp.ID, token, nil)
}

// TestAdminCategoryLifecycle: create → assign to a competition → blocked delete
// (409) → reassign delete succeeds.
func TestAdminCategoryLifecycle(t *testing.T) {
	requireStack(t)

	token := superadminToken(t)

	// Create a throwaway category.
	name := fmt.Sprintf("E2E Cat %d", time.Now().UnixNano())
	catResp, catBody := request(t, http.MethodPost, "/apis/competition/v1/admin/categories", token,
		map[string]any{"name": name})
	requireStatus(t, catResp, catBody, http.StatusCreated)

	var cat category
	decode(t, catBody, &cat)

	// A competition assigned to it.
	slug := fmt.Sprintf("e2e-cat-comp-%d", time.Now().UnixNano())
	compResp, compBody := request(t, http.MethodPost, "/apis/competition/v1/admin/competitions", token,
		map[string]any{
			"title": "Cat Comp", "slug": slug, "prize": "p",
			"ticket_price_pence": 100, "tickets_total": 10,
			"category_id": cat.ID, "status": "draft",
		})
	requireStatus(t, compResp, compBody, http.StatusCreated)

	var comp competition
	decode(t, compBody, &comp)

	// Delete blocked while in use → 409.
	blkResp, blkBody := request(t, http.MethodDelete, "/apis/competition/v1/admin/categories/"+cat.ID, token, nil)
	requireStatus(t, blkResp, blkBody, http.StatusConflict)

	// Reassign target: any OTHER category.
	listResp, listBody := request(t, http.MethodGet, "/apis/competition/v1/categories", "", nil)
	requireStatus(t, listResp, listBody, http.StatusOK)

	var cats categoryList
	decode(t, listBody, &cats)

	var target string
	for _, c := range cats.Categories {
		if c.ID != cat.ID {
			target = c.ID

			break
		}
	}
	if target == "" {
		t.Fatal("no reassignment target category available")
	}

	// Delete with reassignment succeeds; the competition moves to the target.
	okResp, okBody := request(t, http.MethodDelete,
		"/apis/competition/v1/admin/categories/"+cat.ID+"?reassign_to="+target, token, nil)
	requireStatus(t, okResp, okBody, http.StatusNoContent)

	movedResp, movedBody := request(t, http.MethodGet, "/apis/competition/v1/competitions/"+comp.ID, "", nil)
	requireStatus(t, movedResp, movedBody, http.StatusOK)

	var moved competition
	decode(t, movedBody, &moved)
	if moved.CategoryID != target {
		t.Errorf("competition should have been reassigned to %s, got %s", target, moved.CategoryID)
	}

	request(t, http.MethodDelete, "/apis/competition/v1/admin/competitions/"+comp.ID, token, nil)
}

// TestAdminUserEditAndSuspendBlocksPurchase: edit a user's profile, suspend it,
// and assert a suspended user cannot purchase (403).
func TestAdminUserEditAndSuspendBlocksPurchase(t *testing.T) {
	requireStack(t)

	token := superadminToken(t)

	// Register a fresh user.
	email := uniqueEmail()
	regResp, regBody := request(t, http.MethodPost, "/apis/user/v1/users", "",
		registerReq{Name: "Suspend Me", Email: email})
	requireStatus(t, regResp, regBody, http.StatusCreated)

	var u user
	decode(t, regBody, &u)

	// Admin edits the profile.
	newEmail := uniqueEmail()
	editResp, editBody := request(t, http.MethodPut, "/apis/user/v1/admin/users/"+u.ID, token,
		map[string]any{"name": "Renamed", "email": newEmail})
	requireStatus(t, editResp, editBody, http.StatusOK)

	var edited user
	decode(t, editBody, &edited)
	if edited.Name != "Renamed" || edited.Email != newEmail {
		t.Fatalf("profile edit did not round-trip: %+v", edited)
	}

	// Suspend the user.
	suspResp, suspBody := request(t, http.MethodPost, "/apis/user/v1/admin/users/"+u.ID+"/suspend", token, nil)
	requireStatus(t, suspResp, suspBody, http.StatusOK)

	var suspended user
	decode(t, suspBody, &suspended)
	if suspended.IsActive {
		t.Fatal("expected is_active=false after suspend")
	}

	// A suspended user's purchase is rejected with 403.
	comp := firstPaidLiveCompetition(t)
	buyResp, buyBody := request(t, http.MethodPost, "/apis/user/v1/tickets", "",
		purchaseReq{CompetitionID: comp.ID, UserID: u.ID, Quantity: 1})
	requireStatus(t, buyResp, buyBody, http.StatusForbidden)

	// Reactivate → purchase now works.
	request(t, http.MethodPost, "/apis/user/v1/admin/users/"+u.ID+"/activate", token, nil)
	okResp, okBody := request(t, http.MethodPost, "/apis/user/v1/tickets", "",
		purchaseReq{CompetitionID: comp.ID, UserID: u.ID, Quantity: 1})
	requireStatus(t, okResp, okBody, http.StatusCreated)
}

// TestAdminDrawVoidWritesAudit: void a draw with a reason and assert the draw is
// voided (and hidden publicly). If the Postgres DSN is reachable it also asserts
// the admin_audit_log row exists.
func TestAdminDrawVoidWritesAudit(t *testing.T) {
	requireStack(t)

	token := superadminToken(t)

	// Create a fresh competition + a draw so we never void a seeded winner.
	slug := fmt.Sprintf("e2e-void-%d", time.Now().UnixNano())
	compResp, compBody := request(t, http.MethodPost, "/apis/competition/v1/admin/competitions", token,
		map[string]any{
			"title": "Void Comp", "slug": slug, "prize": "p",
			"ticket_price_pence": 100, "tickets_total": 10, "status": "draft",
		})
	requireStatus(t, compResp, compBody, http.StatusCreated)

	var comp competition
	decode(t, compBody, &comp)

	drawResp, drawBody := request(t, http.MethodPost, "/apis/draw/v1/admin/draws", token,
		map[string]any{"competition_id": comp.ID, "prize": "Void Prize"})
	requireStatus(t, drawResp, drawBody, http.StatusCreated)

	var d draw
	decode(t, drawBody, &d)

	// Void without a reason is rejected.
	noReasonResp, noReasonBody := request(t, http.MethodPost, "/apis/draw/v1/admin/draws/"+d.ID+"/void", token,
		map[string]any{"reason": ""})
	requireStatus(t, noReasonResp, noReasonBody, http.StatusBadRequest)

	// Void WITH a reason succeeds.
	reason := "e2e compliance void"
	voidResp, voidBody := request(t, http.MethodPost, "/apis/draw/v1/admin/draws/"+d.ID+"/void", token,
		map[string]any{"reason": reason})
	requireStatus(t, voidResp, voidBody, http.StatusOK)

	var voided draw
	decode(t, voidBody, &voided)
	if voided.Status != "void" || voided.VoidReason != reason {
		t.Fatalf("expected voided draw with reason, got %+v", voided)
	}

	// A voided draw is hidden from the public read.
	pubResp, pubBody := request(t, http.MethodGet, "/apis/draw/v1/draws/"+d.ID, "", nil)
	requireStatus(t, pubResp, pubBody, http.StatusNotFound)

	// Optional deeper check: the audit row exists (needs DB access).
	assertAuditRow(t, "draw", d.ID, "draw.void")

	request(t, http.MethodDelete, "/apis/competition/v1/admin/competitions/"+comp.ID, token, nil)
}

// TestAdminRoleGuards: no token → 401; a role=admin token → 403 on a
// superadmin-only route; the superadmin → 200.
func TestAdminRoleGuards(t *testing.T) {
	requireStack(t)

	t.Run("no token on admin route is 401", func(t *testing.T) {
		resp, body := request(t, http.MethodGet, "/apis/user/v1/admin/users", "", nil)
		requireStatus(t, resp, body, http.StatusUnauthorized)
	})

	t.Run("no token on superadmin route is 401", func(t *testing.T) {
		resp, body := request(t, http.MethodGet, "/apis/adminauth/v1/admin/accounts", "", nil)
		requireStatus(t, resp, body, http.StatusUnauthorized)
	})

	t.Run("role=admin is 403 on superadmin route", func(t *testing.T) {
		resp, body := request(t, http.MethodGet, "/apis/adminauth/v1/admin/accounts", adminToken(t), nil)
		requireStatus(t, resp, body, http.StatusForbidden)
	})

	t.Run("role=admin passes a normal admin route", func(t *testing.T) {
		resp, body := request(t, http.MethodGet, "/apis/user/v1/admin/users", adminToken(t), nil)
		requireStatus(t, resp, body, http.StatusOK)
	})

	t.Run("superadmin passes the superadmin route", func(t *testing.T) {
		resp, body := request(t, http.MethodGet, "/apis/adminauth/v1/admin/accounts", superadminToken(t), nil)
		requireStatus(t, resp, body, http.StatusOK)
	})
}

// --- multipart upload + MinIO + DB helpers ---

// uploadMedia posts a 1x1 PNG to the admin upload endpoint and returns the
// created media record.
func uploadMedia(t *testing.T, token, ownerType, ownerID string, position int) media {
	t.Helper()

	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)

	part, err := w.CreateFormFile("file", "pixel.png")
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := part.Write(onePixelPNG); err != nil {
		t.Fatalf("write file part: %v", err)
	}

	_ = w.WriteField("owner_type", ownerType)
	_ = w.WriteField("owner_id", ownerID)
	_ = w.WriteField("position", fmt.Sprintf("%d", position))
	_ = w.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		baseURL()+"/apis/media/v1/admin/uploads", body)
	if err != nil {
		t.Fatalf("build upload request: %v", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do upload: %v", err)
	}
	defer resp.Body.Close()

	data, _ := readAll(resp)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("upload status: got %d want 201; body=%s", resp.StatusCode, string(data))
	}

	var m media
	decode(t, data, &m)
	if m.ID == "" || m.ObjectKey == "" {
		t.Fatalf("upload returned incomplete media: %+v", m)
	}

	return m
}

// rawStatus does a GET and returns just the status code (used against MinIO).
func rawStatus(t *testing.T, url string) int {
	t.Helper()

	resp, body := request(t, http.MethodGet, url, "", nil)
	_ = body

	return resp.StatusCode
}

// e2ePostgresDSN returns the DSN for the optional direct-DB audit check.
func e2ePostgresDSN() string {
	if v := os.Getenv("E2E_POSTGRES_DSN"); v != "" {
		return v
	}

	return "postgres://user:password@localhost:5432/db?sslmode=disable"
}

// assertAuditRow checks the admin_audit_log for a matching entry. The DB is not
// exposed through the gateway, so this is a best-effort direct query: if the DB
// is unreachable (e.g. Postgres port not published) it logs and skips rather
// than failing, keeping the API-level assertions authoritative.
func assertAuditRow(t *testing.T, entityType, entityID, action string) {
	t.Helper()

	db, err := sql.Open("pgx", e2ePostgresDSN())
	if err != nil {
		t.Logf("audit DB check skipped (open): %v", err)

		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		t.Logf("audit DB check skipped (ping): %v", err)

		return
	}

	var count int
	err = db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM admin_audit_log WHERE entity_type = $1 AND entity_id = $2 AND action = $3`,
		entityType, entityID, action).Scan(&count)
	if err != nil {
		t.Logf("audit DB check skipped (query): %v", err)

		return
	}

	if count < 1 {
		t.Errorf("expected an admin_audit_log row for %s/%s action=%s, found none", entityType, entityID, action)
	}
}
