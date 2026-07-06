package handler_test

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"application/internal/competition/biz"
	"application/internal/competition/entity"
	comphandler "application/internal/competition/handler"
	"application/internal/competition/mocks"
	"application/pkg/audit"
	"application/pkg/middlewares"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newCategoryHandler(t *testing.T) (*http.ServeMux, *mocks.MockUsecaseCategory, *mocks.MockUsecaseContent) {
	t.Helper()

	ucCat := mocks.NewMockUsecaseCategory(t)
	ucContent := mocks.NewMockUsecaseContent(t)
	mux := http.NewServeMux()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	auth := middlewares.NewJWTAuth(middlewares.JWTSecret(testSecret))
	recorder := audit.NewRecorder(logger, nil)

	ch := comphandler.NewCategoryHandler(logger, mux, ucCat, auth, recorder)
	require.NoError(t, ch.RegisterHandler(context.Background()))

	sh := comphandler.NewContentHandler(logger, mux, ucContent, auth, recorder)
	require.NoError(t, sh.RegisterHandler(context.Background()))

	return mux, ucCat, ucContent
}

func TestCategory_PublicListNoToken(t *testing.T) {
	mux, ucCat, _ := newCategoryHandler(t)
	ucCat.EXPECT().List(mock.Anything).Return([]entity.Category{{ID: uuid.New(), Name: "Cars", Slug: "cars"}}, nil)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/apis/competition/v1/categories", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "Cars")
}

func TestCategory_AdminCRUDGuarded(t *testing.T) {
	mux, _, _ := newCategoryHandler(t)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/apis/competition/v1/admin/categories", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestCategory_DeleteInUse409(t *testing.T) {
	mux, ucCat, _ := newCategoryHandler(t)
	id := uuid.New()
	ucCat.EXPECT().Delete(mock.Anything, id, (*uuid.UUID)(nil)).Return(biz.ErrCategoryInUse)

	rec := doJSON(mux, http.MethodDelete, "/apis/competition/v1/admin/categories/"+id.String(), "")
	require.Equal(t, http.StatusConflict, rec.Code)
}

func TestCategory_DeleteWithReassign(t *testing.T) {
	mux, ucCat, _ := newCategoryHandler(t)
	id := uuid.New()
	target := uuid.New()
	ucCat.EXPECT().Delete(mock.Anything, id, mock.MatchedBy(func(r *uuid.UUID) bool {
		return r != nil && *r == target
	})).Return(nil)

	rec := doJSON(mux, http.MethodDelete,
		"/apis/competition/v1/admin/categories/"+id.String()+"?reassign_to="+target.String(), "")
	require.Equal(t, http.StatusNoContent, rec.Code)
}

func TestContent_PublicRead(t *testing.T) {
	mux, _, ucContent := newCategoryHandler(t)
	ucContent.EXPECT().GetAll(mock.Anything).Return([]entity.SiteContent{
		{Key: "winners.count", Value: "9,700 winners"},
	}, nil)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/apis/competition/v1/content", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "9,700 winners")
}

func TestContent_UpsertGuardedAndWorks(t *testing.T) {
	mux, _, ucContent := newCategoryHandler(t)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPut,
		"/apis/competition/v1/admin/content/winners.count", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	require.Equal(t, http.StatusUnauthorized, rec.Code)

	ucContent.EXPECT().Upsert(mock.Anything, "winners.count", "10,000 winners").
		Return(entity.SiteContent{Key: "winners.count", Value: "10,000 winners"}, nil)

	rec = doJSON(mux, http.MethodPut, "/apis/competition/v1/admin/content/winners.count",
		`{"value":"10,000 winners"}`)
	require.Equal(t, http.StatusOK, rec.Code)
}
