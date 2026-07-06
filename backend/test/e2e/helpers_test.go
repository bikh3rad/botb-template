//go:build e2e
// +build e2e

// Package e2e holds end-to-end tests that exercise the LIVE docker-compose stack
// through the gateway (the only published port). There are no mocks and no
// ramsql here: every request travels gateway -> service -> real Postgres/MinIO,
// so these tests assert the wiring the unit tests can only stub.
//
// They are gated behind the `e2e` build tag so `go test ./...` (the unit pass)
// never compiles or runs them. Run explicitly with `make e2e` once the stack is
// up (`docker compose up -d --build && make seed`).
package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

const (
	// defaultBaseURL is the gateway — the single published port of the compose
	// stack. All /apis/ paths are forwarded to their upstream service intact.
	defaultBaseURL = "http://localhost:8080"

	// Bootstrap superadmin credentials — mirror ADMIN_BOOTSTRAP_EMAIL/PASSWORD
	// from .env.example (adminauth seeds this account on first boot). The e2e
	// suite logs in FOR REAL against adminauth; there is no hand-minted token.
	defaultAdminEmail    = "admin@botb.local"
	defaultAdminPassword = "change-me-now"

	// A plain (non-superadmin) admin the suite provisions via the superadmin, so
	// it can assert role=admin gets 403 on superadmin-only routes.
	plainAdminEmail    = "e2e-admin@botb.local"
	plainAdminPassword = "e2e-admin-password"
)

// baseURL returns the gateway base URL (env E2E_BASE_URL overrides the default).
func baseURL() string {
	if v := os.Getenv("E2E_BASE_URL"); v != "" {
		return strings.TrimRight(v, "/")
	}

	return defaultBaseURL
}

// mediaBaseURL returns the host-reachable, public-read MinIO base URL used to
// fetch media objects the same way the browser does (env E2E_MEDIA_BASE_URL
// overrides the default). This mirrors the frontend's NEXT_PUBLIC_MEDIA_BASE_URL.
func mediaBaseURL() string {
	if v := os.Getenv("E2E_MEDIA_BASE_URL"); v != "" {
		return strings.TrimRight(v, "/")
	}

	return "http://localhost:9000"
}

// adminCreds returns the bootstrap superadmin credentials (env
// E2E_ADMIN_EMAIL / E2E_ADMIN_PASSWORD override the compose defaults).
func adminCreds() (email, password string) {
	email = defaultAdminEmail
	if v := os.Getenv("E2E_ADMIN_EMAIL"); v != "" {
		email = v
	}

	password = defaultAdminPassword
	if v := os.Getenv("E2E_ADMIN_PASSWORD"); v != "" {
		password = v
	}

	return email, password
}

// stackErr records whether the gateway was reachable at startup. When non-nil,
// requireStack skips each test with a clear message instead of failing so a
// missing stack degrades gracefully rather than panicking.
var stackErr error

// TestMain probes the gateway root once before any test runs. A single reachable
// check keeps the "stack is down" signal in one place.
func TestMain(m *testing.M) {
	stackErr = probeStack()

	os.Exit(m.Run())
}

// probeStack does a short-timeout GET on the gateway root. Any transport error
// (connection refused, timeout, DNS) means the stack is not up.
func probeStack() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL()+"/", nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// requireStack skips the calling test when the gateway was unreachable at
// startup, turning "no stack" into a clear skip rather than a cascade of
// failures.
func requireStack(t *testing.T) {
	t.Helper()

	if stackErr != nil {
		t.Skipf("gateway unreachable at %s (%v); start the stack with `docker compose up -d --build && make seed`",
			baseURL(), stackErr)
	}
}

// tokenResp mirrors dto.TokenResp from adminauth.
type tokenResp struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Admin        struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	} `json:"admin"`
}

// login authenticates against adminauth and returns the full token response.
func login(t *testing.T, email, password string) tokenResp {
	t.Helper()

	resp, body := request(t, http.MethodPost, "/apis/adminauth/v1/login", "", map[string]string{
		"email":    email,
		"password": password,
	})
	requireStatus(t, resp, body, http.StatusOK)

	var tok tokenResp
	decode(t, body, &tok)

	if tok.AccessToken == "" {
		t.Fatalf("login returned empty access token; body=%s", string(body))
	}

	return tok
}

var (
	superTokenOnce sync.Once
	superTokenVal  string
	adminTokenOnce sync.Once
	adminTokenVal  string
)

// superadminToken logs in as the bootstrapped superadmin (cached per run).
func superadminToken(t *testing.T) string {
	t.Helper()

	superTokenOnce.Do(func() {
		email, password := adminCreds()
		superTokenVal = login(t, email, password).AccessToken
	})

	if superTokenVal == "" {
		t.Fatal("superadmin login failed")
	}

	return superTokenVal
}

// adminToken provisions (idempotently, via the superadmin) a plain role=admin
// account and returns a token for it — used to prove role=admin is accepted on
// admin routes but rejected (403) on superadmin-only routes.
func adminToken(t *testing.T) string {
	t.Helper()

	adminTokenOnce.Do(func() {
		super := superadminToken(t)

		// Create the plain admin; a 409 (already exists from a prior run) is fine.
		// request() already drains + closes the body, so we discard the response.
		_, _ = request(t, http.MethodPost, "/apis/adminauth/v1/admin/accounts", super, map[string]string{
			"name":     "E2E Admin",
			"email":    plainAdminEmail,
			"password": plainAdminPassword,
			"role":     "admin",
		})

		adminTokenVal = login(t, plainAdminEmail, plainAdminPassword).AccessToken
	})

	if adminTokenVal == "" {
		t.Fatal("plain-admin login failed")
	}

	return adminTokenVal
}

// request performs an HTTP call and returns the response plus its fully-read
// body. `target` may be a path (joined to the gateway base URL) or an absolute
// URL (used for the presigned MinIO links, which point straight at object
// storage). `token`, when non-empty, is sent as a bearer credential. `body`,
// when non-nil, is JSON-encoded.
func request(t *testing.T, method, target, token string, body any) (*http.Response, []byte) {
	t.Helper()

	url := target
	if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
		url = baseURL() + target
	}

	var reader *bytes.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal request body: %v", err)
		}
		reader = bytes.NewReader(raw)
	} else {
		reader = bytes.NewReader(nil)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, url, reader)
	if err != nil {
		t.Fatalf("build request %s %s: %v", method, url, err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request %s %s: %v", method, url, err)
	}
	defer resp.Body.Close()

	data, err := readAll(resp)
	if err != nil {
		t.Fatalf("read response body %s %s: %v", method, url, err)
	}

	return resp, data
}

// readAll drains the response body. Split out so request stays readable.
func readAll(resp *http.Response) ([]byte, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// requireStatus fails (fatally) unless the response carries the wanted status,
// printing the body to make mismatches diagnosable.
func requireStatus(t *testing.T, resp *http.Response, body []byte, want int) {
	t.Helper()

	if resp.StatusCode != want {
		t.Fatalf("unexpected status: got %d want %d; body=%s", resp.StatusCode, want, string(body))
	}
}

// decode unmarshals a JSON body into out, failing the test on malformed JSON.
func decode(t *testing.T, body []byte, out any) {
	t.Helper()

	if err := json.Unmarshal(body, out); err != nil {
		t.Fatalf("decode JSON: %v; body=%s", err, string(body))
	}
}

// --- API response shapes (subset of the service DTOs we assert on) ---

type mediaRef struct {
	ID          string `json:"id"`
	Kind        string `json:"kind"`
	Bucket      string `json:"bucket"`
	ObjectKey   string `json:"object_key"`
	ContentType string `json:"content_type"`
	Position    int    `json:"position"`
}

type competition struct {
	ID               string     `json:"id"`
	Title            string     `json:"title"`
	Slug             string     `json:"slug"`
	Description      string     `json:"description"`
	Prize            string     `json:"prize"`
	TicketPricePence int64      `json:"ticket_price_pence"`
	TicketsTotal     int64      `json:"tickets_total"`
	TicketsSold      int64      `json:"tickets_sold"`
	CategoryID       string     `json:"category_id"`
	CategoryName     string     `json:"category_name"`
	Status           string     `json:"status"`
	StartsAt         string     `json:"starts_at"`
	EndsAt           string     `json:"ends_at"`
	CreatedAt        string     `json:"created_at"`
	UpdatedAt        string     `json:"updated_at"`
	Media            []mediaRef `json:"media"`
}

type category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type categoryList struct {
	Count      int        `json:"count"`
	Categories []category `json:"categories"`
}

type contentResp struct {
	Items map[string]string `json:"items"`
}

type competitionList struct {
	Count        int           `json:"count"`
	Competitions []competition `json:"competitions"`
}

type media struct {
	ID          string `json:"id"`
	OwnerType   string `json:"owner_type"`
	OwnerID     string `json:"owner_id"`
	Kind        string `json:"kind"`
	Bucket      string `json:"bucket"`
	ObjectKey   string `json:"object_key"`
	ContentType string `json:"content_type"`
	Position    int    `json:"position"`
	URL         string `json:"url"`
}

type draw struct {
	ID             string `json:"id"`
	CompetitionID  string `json:"competition_id"`
	WinnerUserID   string `json:"winner_user_id"`
	WinnerTicketID string `json:"winner_ticket_id"`
	Prize          string `json:"prize"`
	Status         string `json:"status"`
	VoidReason     string `json:"void_reason"`
	DrawnAt        string `json:"drawn_at"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type drawList struct {
	Count int    `json:"count"`
	Draws []draw `json:"draws"`
}

// user mirrors dto.UserResp (register + admin get responses).
type user struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	TicketsOwned    int64  `json:"tickets_owned"`
	TotalSpentPence int64  `json:"total_spent_pence"`
	IsActive        bool   `json:"is_active"`
	CreatedAt       string `json:"created_at"`
}

// registerReq mirrors dto.RegisterReq: {name, email}.
type registerReq struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// purchaseReq mirrors dto.PurchaseReq: {competition_id, user_id, quantity}.
type purchaseReq struct {
	CompetitionID string `json:"competition_id"`
	UserID        string `json:"user_id"`
	Quantity      int    `json:"quantity"`
}

// purchaseResp mirrors dto.PurchaseResp.
type purchaseResp struct {
	User           user  `json:"user"`
	TotalCostPence int64 `json:"total_cost_pence"`
	Count          int   `json:"count"`
	Tickets        []struct {
		ID            string `json:"id"`
		CompetitionID string `json:"competition_id"`
		UserID        string `json:"user_id"`
	} `json:"tickets"`
}

// uniqueEmail returns a per-run unique address so registration never collides
// with a prior run (the seeder derives stable IDs from email, so reuse would
// upsert instead of insert).
func uniqueEmail() string {
	return fmt.Sprintf("e2e+%d@example.test", time.Now().UnixNano())
}
