package handler_test

import (
	"application/internal/user/biz"
	"application/internal/user/dto"
	"application/internal/user/entity"
	userhandler "application/internal/user/handler"
	"application/internal/user/mocks"
	"application/pkg/middlewares"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const testSecret = "test-secret"

func validToken() string {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  "admin",
		"role": "admin",
		"exp":  time.Now().Add(time.Hour).Unix(),
	})
	s, _ := tok.SignedString([]byte(testSecret))

	return s
}

type harness struct {
	mux    *http.ServeMux
	userUC *mocks.MockUsecaseUser
	tktUC  *mocks.MockUsecaseTicket
}

func newHarness(t *testing.T) harness {
	t.Helper()

	userUC := mocks.NewMockUsecaseUser(t)
	tktUC := mocks.NewMockUsecaseTicket(t)
	mux := http.NewServeMux()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	auth := middlewares.NewJWTAuth(middlewares.JWTSecret(testSecret))

	uh := userhandler.NewUser(logger, mux, userUC, auth)
	require.NoError(t, uh.RegisterHandler(context.Background()))

	th := userhandler.NewTicket(logger, mux, tktUC, auth)
	require.NoError(t, th.RegisterHandler(context.Background()))

	return harness{mux: mux, userUC: userUC, tktUC: tktUC}
}

// do issues a request carrying a valid admin token (public routes ignore it).
func (h harness) do(method, target, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequestWithContext(context.Background(), method, target, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+validToken())

	rec := httptest.NewRecorder()
	h.mux.ServeHTTP(rec, req)

	return rec
}

// TestAdmin_RequiresToken proves the user service's own admin group rejects an
// unauthenticated call (defense in depth actually wired).
func TestAdmin_RequiresToken(t *testing.T) {
	h := newHarness(t)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet,
		"/apis/user/v1/admin/users", nil)
	rec := httptest.NewRecorder()
	h.mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestRegister_Created(t *testing.T) {
	h := newHarness(t)
	h.userUC.EXPECT().Register(mock.Anything, "Olivia", "olivia@example.com").
		Return(entity.User{ID: uuid.New(), Name: "Olivia", Email: "olivia@example.com"}, nil)

	rec := h.do(http.MethodPost, "/apis/user/v1/users", `{"name":"Olivia","email":"olivia@example.com"}`)
	require.Equal(t, http.StatusCreated, rec.Code)
}

func TestRegister_BadJSON(t *testing.T) {
	h := newHarness(t)
	rec := h.do(http.MethodPost, "/apis/user/v1/users", `{bad`)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRegister_Conflict(t *testing.T) {
	h := newHarness(t)
	h.userUC.EXPECT().Register(mock.Anything, mock.Anything, mock.Anything).
		Return(entity.User{}, biz.ErrResourceExists)

	rec := h.do(http.MethodPost, "/apis/user/v1/users", `{"name":"A","email":"a@b.com"}`)
	require.Equal(t, http.StatusConflict, rec.Code)
}

func TestListUsers_OK(t *testing.T) {
	h := newHarness(t)
	h.userUC.EXPECT().List(mock.Anything, mock.Anything).
		Return(biz.UserPage{Users: []entity.User{{ID: uuid.New(), Name: "A"}}, Total: 1}, nil)

	rec := h.do(http.MethodGet, "/apis/user/v1/admin/users?limit=10", "")
	require.Equal(t, http.StatusOK, rec.Code)

	var resp dto.UserListResp

	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	require.Equal(t, 1, resp.Total)
}

func TestGetUser_NotFound(t *testing.T) {
	h := newHarness(t)
	id := uuid.New()
	h.userUC.EXPECT().Get(mock.Anything, id).Return(entity.User{}, biz.ErrResourceNotFound)

	rec := h.do(http.MethodGet, "/apis/user/v1/admin/users/"+id.String(), "")
	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetUser_InvalidID(t *testing.T) {
	h := newHarness(t)
	rec := h.do(http.MethodGet, "/apis/user/v1/admin/users/not-a-uuid", "")
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPurchase_Created(t *testing.T) {
	h := newHarness(t)
	compID := uuid.New()
	userID := uuid.New()

	h.tktUC.EXPECT().Purchase(mock.Anything, mock.Anything).Return(biz.PurchaseResult{
		Tickets:        []entity.Ticket{{ID: uuid.New()}},
		User:           entity.User{ID: userID, TicketsOwned: 1, TotalSpentPence: 125},
		TotalCostPence: 125,
	}, nil)

	body := `{"competition_id":"` + compID.String() + `","user_id":"` + userID.String() + `","quantity":1}`
	rec := h.do(http.MethodPost, "/apis/user/v1/tickets", body)
	require.Equal(t, http.StatusCreated, rec.Code)

	var resp dto.PurchaseResp

	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	require.Equal(t, 1, resp.Count)
	require.Equal(t, int64(125), resp.TotalCostPence)
}

func TestPurchase_InvalidUUID(t *testing.T) {
	h := newHarness(t)
	rec := h.do(http.MethodPost, "/apis/user/v1/tickets", `{"competition_id":"x","user_id":"y","quantity":1}`)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPurchase_CompetitionNotFound(t *testing.T) {
	h := newHarness(t)
	compID := uuid.New()
	userID := uuid.New()

	h.tktUC.EXPECT().Purchase(mock.Anything, mock.Anything).
		Return(biz.PurchaseResult{}, biz.ErrCompetitionNotFound)

	body := `{"competition_id":"` + compID.String() + `","user_id":"` + userID.String() + `","quantity":1}`
	rec := h.do(http.MethodPost, "/apis/user/v1/tickets", body)
	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestListTickets_OK(t *testing.T) {
	h := newHarness(t)
	userID := uuid.New()
	h.tktUC.EXPECT().ListByUser(mock.Anything, userID).
		Return([]entity.Ticket{{ID: uuid.New(), UserID: userID}}, nil)

	rec := h.do(http.MethodGet, "/apis/user/v1/admin/users/"+userID.String()+"/tickets", "")
	require.Equal(t, http.StatusOK, rec.Code)
}
