package icinga2

import "errors"

var (
	ErrNotFound     = errors.New("service not found")
	ErrForbidden    = errors.New("forbidden")
	ErrUnauthorized = errors.New("unauthorized")
	ErrUnknown      = errors.New("unknown error")
)
