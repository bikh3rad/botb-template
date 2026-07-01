package handler_test

import (
	"application/internal/media/biz"
	"application/internal/media/dto"
	"application/internal/media/entity"
	mediahandler "application/internal/media/handler"
	"application/internal/media/mocks"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newTestHandler(t *testing.T) (*http.ServeMux, *mocks.MockUsecaseMedia) {
	t.Helper()

	uc := mocks.NewMockUsecaseMedia(t)
	mux := http.NewServeMux()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	h := mediahandler.NewMedia(logger, mux, uc)
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
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/apis/media/v1/uploads", body)
	req.Header.Set("Content-Type", contentType)

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
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/apis/media/v1/uploads", body)
	req.Header.Set("Content-Type", contentType)

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

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/apis/media/v1/uploads", body)
	req.Header.Set("Content-Type", w.FormDataContentType())

	rec := doRequest(mux, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpload_InvalidOwnerID(t *testing.T) {
	mux, _ := newTestHandler(t)

	body, contentType := multipartUpload(t, "image/png", "competition", "not-a-uuid")
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/apis/media/v1/uploads", body)
	req.Header.Set("Content-Type", contentType)

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
