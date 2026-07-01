package dto

import (
	"application/internal/draw/biz"
	"encoding/json"
	"errors"
	"net/http"
)

// ErrorResponse is the standard error envelope for the draw API.
type ErrorResponse struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type errorInfo struct {
	Message string
	Code    int
}

var errorsMap = map[error]errorInfo{
	biz.ErrResourceNotFound: {Message: "draw not found", Code: http.StatusNotFound},
	biz.ErrResourceInvalid:  {Message: "invalid request", Code: http.StatusBadRequest},
	biz.ErrAlreadyDrawn:     {Message: "draw already run", Code: http.StatusConflict},
	biz.ErrNoTickets:        {Message: "competition has no tickets to draw", Code: http.StatusUnprocessableEntity},
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
