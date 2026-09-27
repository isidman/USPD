package domain

// Sentinel-style error types the API layer maps to HTTP status codes, and
// tests can assert on with errors.As — see internal/api/handlers.go.

type ErrInvalidInput string

func (e ErrInvalidInput) Error() string { return string(e) }

type ErrNotFound string

func (e ErrNotFound) Error() string { return string(e) }

type ErrConflict string

func (e ErrConflict) Error() string { return string(e) }
