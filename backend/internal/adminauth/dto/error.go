package dto

import (
	"encoding/json"
	"errors"
	"net/http"

	"application/internal/adminauth/biz"
)

// ErrorResponse is the standard error envelope for the adminauth API.
type ErrorResponse struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type errorInfo struct {
	Message string
	Code    int
}

// Login failures deliberately map to ONE generic message (no Details) so a
// caller cannot distinguish unknown email / wrong password / disabled account.
var errorsMap = map[error]errorInfo{
	biz.ErrInvalidCredentials: {Message: "invalid credentials", Code: http.StatusUnauthorized},
	biz.ErrInvalidRefresh:     {Message: "invalid refresh token", Code: http.StatusUnauthorized},
	biz.ErrRateLimited:        {Message: "too many attempts, retry later", Code: http.StatusTooManyRequests},
	biz.ErrResourceNotFound:   {Message: "admin account not found", Code: http.StatusNotFound},
	biz.ErrResourceInvalid:    {Message: "invalid request", Code: http.StatusBadRequest},
	biz.ErrResourceExists:     {Message: "email already in use", Code: http.StatusConflict},
	biz.ErrLastSuperadmin:     {Message: "cannot disable or demote the last active superadmin", Code: http.StatusConflict},
}

// credentialErrors never include error details in the response body.
var credentialErrors = []error{biz.ErrInvalidCredentials, biz.ErrInvalidRefresh, biz.ErrRateLimited}

// HandleError writes a JSON error response, mapping known sentinels to codes.
func HandleError(err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	for e, info := range errorsMap {
		if errors.Is(err, e) {
			w.WriteHeader(info.Code)

			resp := ErrorResponse{Message: info.Message, Details: err.Error()}
			for _, ce := range credentialErrors {
				if errors.Is(err, ce) {
					resp.Details = ""
				}
			}

			_ = encoder.Encode(resp)

			return
		}
	}

	w.WriteHeader(http.StatusInternalServerError)
	_ = encoder.Encode(ErrorResponse{Message: "internal server error"})
}
