package dto

import (
	"application/internal/competition/biz"
	"encoding/json"
	"errors"
	"net/http"
)

// ErrorResponse is the standard error envelope for the competition API.
type ErrorResponse struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type errorInfo struct {
	Message string
	Code    int
}

var errorsMap = map[error]errorInfo{
	biz.ErrResourceNotFound: {Message: "competition not found", Code: http.StatusNotFound},
	biz.ErrResourceInvalid:  {Message: "invalid request", Code: http.StatusBadRequest},
	biz.ErrResourceExists:   {Message: "competition already exists", Code: http.StatusConflict},
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
