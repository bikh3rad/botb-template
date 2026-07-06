package dto

import (
	"encoding/json"
	"errors"
	"net/http"

	"application/internal/media/biz"
)

// ErrorResponse is the standard error envelope for the media API.
type ErrorResponse struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type errorInfo struct {
	Message string
	Code    int
}

// errorsMap translates media biz sentinels to HTTP responses. Mirrors the
// template's dto/error.go pattern, scoped to this service's error set.
var errorsMap = map[error]errorInfo{
	biz.ErrResourceNotFound: {Message: "media not found", Code: http.StatusNotFound},
	biz.ErrResourceInvalid:  {Message: "invalid request", Code: http.StatusBadRequest},
	biz.ErrUnsupportedType:  {Message: "unsupported media type", Code: http.StatusUnsupportedMediaType},
	biz.ErrFileTooLarge:     {Message: "file too large", Code: http.StatusRequestEntityTooLarge},
}

// HandleError writes a JSON error response, mapping known sentinels to codes.
func HandleError(err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	for e, info := range errorsMap {
		if errors.Is(err, e) {
			w.WriteHeader(info.Code)
			_ = encoder.Encode(ErrorResponse{Message: info.Message, Details: err.Error()})

			return
		}
	}

	w.WriteHeader(http.StatusInternalServerError)

	_ = encoder.Encode(ErrorResponse{Message: "internal server error", Details: err.Error()})
}
