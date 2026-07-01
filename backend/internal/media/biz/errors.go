package biz

import "errors"

var (
	// ErrResourceNotFound is returned when a media record does not exist.
	ErrResourceNotFound = errors.New("media resource not found")
	// ErrResourceInvalid is returned for malformed input (bad UUID, empty owner…).
	ErrResourceInvalid = errors.New("invalid media resource")
	// ErrUnsupportedType is returned when the content type is not allow-listed.
	ErrUnsupportedType = errors.New("unsupported media type")
	// ErrFileTooLarge is returned when the upload exceeds the per-kind limit.
	ErrFileTooLarge = errors.New("media file too large")
)
