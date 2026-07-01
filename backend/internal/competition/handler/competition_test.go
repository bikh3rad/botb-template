package handler_test

import (
	"application/internal/competition/biz"
	"application/internal/competition/dto"
	"application/internal/competition/entity"
	comphandler "application/internal/competition/handler"
	"application/internal/competition/mocks"
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

// validToken signs an admin bearer token with the test secret.
func validToken() string {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "admin",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	s, _ := tok.SignedString([]byte(testSecret))

	return s
}

func newTestHandler(t *testing.T) (*http.ServeMux, *mocks.MockUsecaseCompetition) {
	t.Helper()

	uc := mocks.NewMockUsecaseCompetition(t)
	mux := http.NewServeMux()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	auth := middlewares.NewJWTAuth(middlewares.JWTSecret(testSecret))
	h := comphandler.NewCompetition(logger, mux, uc, auth)
	require.NoError(t, h.RegisterHandler(context.Background()))

	return mux, uc
}

// doJSON issues a request carrying a valid admin token (public routes ignore it).
func doJSON(mux *http.ServeMux, method, target, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequestWithContext(context.Background(), method, target, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+validToken())

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	return rec
}

// TestAdmin_RequiresToken proves the service's own admin group rejects an
// unauthenticated call (defense in depth is actually wired, not just at the gateway).
func TestAdmin_RequiresToken(t *testing.T) {
	mux, _ := newTestHandler(t)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost,
		"/apis/competition/v1/admin/competitions", strings.NewReader(`{"title":"A"}`))
	req.Header.Set("Content-Type", "application/json")
	// no Authorization header

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestList_OK(t *testing.T) {
	mux, uc := newTestHandler(t)
	uc.EXPECT().List(mock.Anything, mock.Anything).
		Return([]entity.Competition{{ID: uuid.New(), Title: "A"}}, nil)

	rec := doJSON(mux, http.MethodGet, "/apis/competition/v1/competitions", "")
	require.Equal(t, http.StatusOK, rec.Code)

	var resp dto.CompetitionListResp

	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	require.Equal(t, 1, resp.Count)
}

func TestList_StatusFilter(t *testing.T) {
	mux, uc := newTestHandler(t)
	live := entity.StatusLive
	uc.EXPECT().List(mock.Anything, biz.ListFilter{Status: &live}).
		Return([]entity.Competition{}, nil)

	rec := doJSON(mux, http.MethodGet, "/apis/competition/v1/competitions?status=live", "")
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestGet_OK_WithMedia(t *testing.T) {
	mux, uc := newTestHandler(t)
	id := uuid.New()
	uc.EXPECT().Get(mock.Anything, id).Return(entity.Competition{
		ID:    id,
		Title: "A",
		Media: []entity.MediaRef{{ID: uuid.New(), Kind: "image", ObjectKey: "k"}},
	}, nil)

	rec := doJSON(mux, http.MethodGet, "/apis/competition/v1/competitions/"+id.String(), "")
	require.Equal(t, http.StatusOK, rec.Code)

	var resp dto.CompetitionResp

	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	require.Len(t, resp.Media, 1)
}

func TestGet_InvalidID(t *testing.T) {
	mux, _ := newTestHandler(t)
	rec := doJSON(mux, http.MethodGet, "/apis/competition/v1/competitions/not-a-uuid", "")
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGet_NotFound(t *testing.T) {
	mux, uc := newTestHandler(t)
	id := uuid.New()
	uc.EXPECT().Get(mock.Anything, id).Return(entity.Competition{}, biz.ErrResourceNotFound)

	rec := doJSON(mux, http.MethodGet, "/apis/competition/v1/competitions/"+id.String(), "")
	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCreate_Created(t *testing.T) {
	mux, uc := newTestHandler(t)
	uc.EXPECT().Create(mock.Anything, mock.Anything).
		Return(entity.Competition{ID: uuid.New(), Title: "A", Status: entity.StatusDraft}, nil)

	body := `{"title":"A","prize":"Car","ticket_price_pence":125,"tickets_total":1000}`
	rec := doJSON(mux, http.MethodPost, "/apis/competition/v1/admin/competitions", body)
	require.Equal(t, http.StatusCreated, rec.Code)
}

func TestCreate_BadJSON(t *testing.T) {
	mux, _ := newTestHandler(t)
	rec := doJSON(mux, http.MethodPost, "/apis/competition/v1/admin/competitions", `{not json`)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreate_Conflict(t *testing.T) {
	mux, uc := newTestHandler(t)
	uc.EXPECT().Create(mock.Anything, mock.Anything).
		Return(entity.Competition{}, biz.ErrResourceExists)

	body := `{"title":"A","prize":"Car","tickets_total":1000}`
	rec := doJSON(mux, http.MethodPost, "/apis/competition/v1/admin/competitions", body)
	require.Equal(t, http.StatusConflict, rec.Code)
}

func TestUpdate_OK(t *testing.T) {
	mux, uc := newTestHandler(t)
	id := uuid.New()
	uc.EXPECT().Update(mock.Anything, id, mock.Anything).
		Return(entity.Competition{ID: id, Title: "B", Status: entity.StatusLive}, nil)

	body := `{"title":"B","prize":"Car","status":"live","tickets_total":1000}`
	rec := doJSON(mux, http.MethodPut, "/apis/competition/v1/admin/competitions/"+id.String(), body)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestDelete_NoContent(t *testing.T) {
	mux, uc := newTestHandler(t)
	id := uuid.New()
	uc.EXPECT().Delete(mock.Anything, id).Return(nil)

	rec := doJSON(mux, http.MethodDelete, "/apis/competition/v1/admin/competitions/"+id.String(), "")
	require.Equal(t, http.StatusNoContent, rec.Code)
}

func TestDelete_NotFound(t *testing.T) {
	mux, uc := newTestHandler(t)
	id := uuid.New()
	uc.EXPECT().Delete(mock.Anything, id).Return(biz.ErrResourceNotFound)

	rec := doJSON(mux, http.MethodDelete, "/apis/competition/v1/admin/competitions/"+id.String(), "")
	require.Equal(t, http.StatusNotFound, rec.Code)
}
