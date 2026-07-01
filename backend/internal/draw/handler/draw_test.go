package handler_test

import (
	"application/internal/draw/biz"
	"application/internal/draw/dto"
	"application/internal/draw/entity"
	drawhandler "application/internal/draw/handler"
	"application/internal/draw/mocks"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newTestHandler(t *testing.T) (*http.ServeMux, *mocks.MockUsecaseDraw) {
	t.Helper()

	uc := mocks.NewMockUsecaseDraw(t)
	mux := http.NewServeMux()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	h := drawhandler.NewDraw(logger, mux, uc)
	require.NoError(t, h.RegisterHandler(context.Background()))

	return mux, uc
}

func doJSON(mux *http.ServeMux, method, target, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequestWithContext(context.Background(), method, target, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	return rec
}

func TestList_OK(t *testing.T) {
	mux, uc := newTestHandler(t)
	uc.EXPECT().List(mock.Anything, mock.Anything).
		Return(biz.DrawPage{Draws: []entity.Draw{{ID: uuid.New()}}, Total: 1}, nil)

	rec := doJSON(mux, http.MethodGet, "/apis/draw/v1/admin/draws?limit=10", "")
	require.Equal(t, http.StatusOK, rec.Code)

	var resp dto.DrawListResp

	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	require.Equal(t, 1, resp.Total)
}

func TestGet_NotFound(t *testing.T) {
	mux, uc := newTestHandler(t)
	id := uuid.New()
	uc.EXPECT().Get(mock.Anything, id).Return(entity.Draw{}, biz.ErrResourceNotFound)

	rec := doJSON(mux, http.MethodGet, "/apis/draw/v1/admin/draws/"+id.String(), "")
	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGet_InvalidID(t *testing.T) {
	mux, _ := newTestHandler(t)
	rec := doJSON(mux, http.MethodGet, "/apis/draw/v1/admin/draws/not-a-uuid", "")
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreate_Created(t *testing.T) {
	mux, uc := newTestHandler(t)
	compID := uuid.New()
	uc.EXPECT().Create(mock.Anything, mock.Anything).
		Return(entity.Draw{ID: uuid.New(), CompetitionID: compID, Status: entity.StatusPending}, nil)

	body := `{"competition_id":"` + compID.String() + `","prize":"Audi RS3"}`
	rec := doJSON(mux, http.MethodPost, "/apis/draw/v1/admin/draws", body)
	require.Equal(t, http.StatusCreated, rec.Code)
}

func TestCreate_BadJSON(t *testing.T) {
	mux, _ := newTestHandler(t)
	rec := doJSON(mux, http.MethodPost, "/apis/draw/v1/admin/draws", `{bad`)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreate_InvalidCompetitionID(t *testing.T) {
	mux, _ := newTestHandler(t)
	rec := doJSON(mux, http.MethodPost, "/apis/draw/v1/admin/draws", `{"competition_id":"x","prize":"y"}`)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRun_OK(t *testing.T) {
	mux, uc := newTestHandler(t)
	id := uuid.New()
	winner := uuid.New()
	uc.EXPECT().Run(mock.Anything, id).
		Return(entity.Draw{ID: id, Status: entity.StatusDrawn, WinnerUserID: &winner}, nil)

	rec := doJSON(mux, http.MethodPost, "/apis/draw/v1/admin/draws/"+id.String()+"/run", "")
	require.Equal(t, http.StatusOK, rec.Code)

	var resp dto.DrawResp

	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	require.Equal(t, "drawn", resp.Status)
	require.NotEmpty(t, resp.WinnerUserID)
}

func TestRun_AlreadyDrawn(t *testing.T) {
	mux, uc := newTestHandler(t)
	id := uuid.New()
	uc.EXPECT().Run(mock.Anything, id).Return(entity.Draw{}, biz.ErrAlreadyDrawn)

	rec := doJSON(mux, http.MethodPost, "/apis/draw/v1/admin/draws/"+id.String()+"/run", "")
	require.Equal(t, http.StatusConflict, rec.Code)
}

func TestRun_NoTickets(t *testing.T) {
	mux, uc := newTestHandler(t)
	id := uuid.New()
	uc.EXPECT().Run(mock.Anything, id).Return(entity.Draw{}, biz.ErrNoTickets)

	rec := doJSON(mux, http.MethodPost, "/apis/draw/v1/admin/draws/"+id.String()+"/run", "")
	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestGetPublic_OK(t *testing.T) {
	mux, uc := newTestHandler(t)
	id := uuid.New()
	uc.EXPECT().GetPublic(mock.Anything, id).
		Return(entity.Draw{ID: id, Status: entity.StatusDrawn}, nil)

	rec := doJSON(mux, http.MethodGet, "/apis/draw/v1/draws/"+id.String(), "")
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestGetPublic_HidesPending(t *testing.T) {
	mux, uc := newTestHandler(t)
	id := uuid.New()
	// The use case maps a pending draw to not-found for public callers.
	uc.EXPECT().GetPublic(mock.Anything, id).Return(entity.Draw{}, biz.ErrResourceNotFound)

	rec := doJSON(mux, http.MethodGet, "/apis/draw/v1/draws/"+id.String(), "")
	require.Equal(t, http.StatusNotFound, rec.Code)
}
