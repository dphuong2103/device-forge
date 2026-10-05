package secrets

import "errors"

var (
	ErrUnsupported = errors.New("Windows protected secret storage is unavailable on this host")
	ErrNotFound    = errors.New("secret not found")
)
