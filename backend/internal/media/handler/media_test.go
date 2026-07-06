package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strings"
	"testing"
	"time"

	"application/internal/media/biz"
	"application/internal/media/dto"
	"application/internal/media/entity"
	mediahandler "application/internal/media/handler"
	"application/internal/media/mocks"
	"application/pkg/audit"
	"application/pkg/middlewares"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const testSecret = "test-secret"

// adminToken signs an HS256 bearer token with role=admin (the media admin
// group requires it now).
func adminToken(t *testing.T) string {
	t.Helper()

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  "admin",
		"role": "admin",
		"exp":  time.Now().Add(time.Hour).Unix(),
	})
	s, err := tok.SignedString([]byte(testSecret))
	require.NoError(t, err)

	return s
}

func newTestHandler(t *testing.T) (*http.ServeMux, *mocks.MockUsecaseMedia) {
	t.Helper()

	uc := mocks.NewMockUsecaseMedia(t)
	mux := http.NewServeMux()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	auth := middlewares.NewJWTAuth(middlewares.JWTSecret(testSecret))

	h := mediahandler.NewMedia(logger, mux, uc, auth, audit.NewRecorder(logger, nil))
	require.NoError(t, h.RegisterHandler(context.Background()))

	return mux, uc
}

// multipartUpload builds a multipart body with a file part carrying the given
// content type plus the owner form fields.
func multipartUpload(t *testing.T, contentType, ownerType, ownerID string) (*bytes.Buffer, string) {
	t.Helper()

	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)

	hdr := textproto.MIMEHeader{}
	hdr.Set("Content-Disposition", `form-data; name="file"; filename="upload.bin"`)
	hdr.Set("Content-Type", contentType)

	part, err := w.CreatePart(hdr)
	require.NoError(t, err)

	_, err = part.Write([]byte("file-bytes"))
	require.NoError(t, err)

	require.NoError(t, w.WriteField("owner_type", ownerType))
	require.NoError(t, w.WriteField("owner_id", ownerID))
	require.NoError(t, w.Close())

	return body, w.FormDataContentType()
}

// doRequest issues req against the handler and returns the recorder.
func doRequest(mux *http.ServeMux, req *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	return rec
}

func TestUpload_Created(t *testing.T) {
	mux, uc := newTestHandler(t)
	ownerID := uuid.New()

	stored := entity.Media{ID: uuid.New(), OwnerType: "competition", OwnerID: ownerID, Kind: entity.KindImage}
	uc.EXPECT().Upload(mock.Anything, mock.Anything).Return(stored, nil)

	body, contentType := multipartUpload(t, "image/png", "competition", ownerID.String())
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/apis/media/v1/admin/uploads", body)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Bearer "+adminToken(t))

	rec := doRequest(mux, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	var resp dto.MediaResp

	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	require.Equal(t, "image", resp.Kind)
}

func TestUpload_UnsupportedType(t *testing.T) {
	mux, uc := newTestHandler(t)
	ownerID := uuid.New()

	uc.EXPECT().Upload(mock.Anything, mock.Anything).Return(entity.Media{}, biz.ErrUnsupportedType)

	body, contentType := multipartUpload(t, "application/zip", "competition", ownerID.String())
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/apis/media/v1/admin/uploads", body)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Bearer "+adminToken(t))

	rec := doRequest(mux, req)
	require.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
}

func TestUpload_MissingFile(t *testing.T) {
	mux, _ := newTestHandler(t)

	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	require.NoError(t, w.WriteField("owner_type", "competition"))
	require.NoError(t, w.WriteField("owner_id", uuid.New().String()))
	require.NoError(t, w.Close())

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/apis/media/v1/admin/uploads", body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+adminToken(t))

	rec := doRequest(mux, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpload_InvalidOwnerID(t *testing.T) {
	mux, _ := newTestHandler(t)

	body, contentType := multipartUpload(t, "image/png", "competition", "not-a-uuid")
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/apis/media/v1/admin/uploads", body)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Bearer "+adminToken(t))

	rec := doRequest(mux, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGet_OK(t *testing.T) {
	mux, uc := newTestHandler(t)
	id := uuid.New()

	uc.EXPECT().Get(mock.Anything, id).Return(biz.MediaWithURL{
		Media: entity.Media{ID: id, Kind: entity.KindImage},
		URL:   "https://minio/presigned",
	}, nil)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/apis/media/v1/media/"+id.String(), nil)

	rec := doRequest(mux, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp dto.MediaResp

	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	require.Equal(t, "https://minio/presigned", resp.URL)
}

func TestGet_NotFound(t *testing.T) {
	mux, uc := newTestHandler(t)
	id := uuid.New()

	uc.EXPECT().Get(mock.Anything, id).Return(biz.MediaWithURL{}, biz.ErrResourceNotFound)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/apis/media/v1/media/"+id.String(), nil)

	rec := doRequest(mux, req)
	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGet_InvalidID(t *testing.T) {
	mux, _ := newTestHandler(t)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/apis/media/v1/media/not-a-uuid", nil)

	rec := doRequest(mux, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

// The old unauthenticated POST /apis/media/v1/uploads is GONE — uploads only
// exist under the admin group now.
func TestOldPublicUploadRouteRemoved(t *testing.T) {
	mux, _ := newTestHandler(t)

	body, contentType := multipartUpload(t, "image/png", "competition", uuid.NewString())
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/apis/media/v1/uploads", body)
	req.Header.Set("Content-Type", contentType)

	rec := doRequest(mux, req)
	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestAdminUploadRequiresToken(t *testing.T) {
	mux, _ := newTestHandler(t)

	body, contentType := multipartUpload(t, "image/png", "competition", uuid.NewString())
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/apis/media/v1/admin/uploads", body)
	req.Header.Set("Content-Type", contentType)
	// no Authorization header

	rec := doRequest(mux, req)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAdminDelete_NoContent(t *testing.T) {
	mux, uc := newTestHandler(t)
	id := uuid.New()
	uc.EXPECT().Delete(mock.Anything, id).Return(nil)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodDelete,
		"/apis/media/v1/admin/media/"+id.String(), nil)
	req.Header.Set("Authorization", "Bearer "+adminToken(t))

	rec := doRequest(mux, req)
	require.Equal(t, http.StatusNoContent, rec.Code)
}

func TestAdminUpdate_Position(t *testing.T) {
	mux, uc := newTestHandler(t)
	id := uuid.New()
	uc.EXPECT().Update(mock.Anything, id, mock.MatchedBy(func(in biz.UpdateInput) bool {
		return in.Position != nil && *in.Position == 3
	})).Return(entity.Media{ID: id, Position: 3, Kind: entity.KindImage}, nil)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPut,
		"/apis/media/v1/admin/media/"+id.String(), strings.NewReader(`{"position":3}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken(t))

	rec := doRequest(mux, req)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestPublicReadsStillTokenless(t *testing.T) {
	mux, uc := newTestHandler(t)
	ownerID := uuid.New()
	uc.EXPECT().ListByOwner(mock.Anything, "competition", ownerID).Return([]entity.Media{}, nil)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet,
		"/apis/media/v1/media?owner_type=competition&owner_id="+ownerID.String(), nil)

	rec := doRequest(mux, req)
	require.Equal(t, http.StatusOK, rec.Code)
}
