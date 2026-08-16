package repository

import "errors"

var (
	// ErrNotFound is returned when a record does not exist.
	ErrNotFound = errors.New("not found")
	// ErrConflict is returned when a unique constraint would be violated.
	ErrConflict = errors.New("conflict")
	// ErrNoRowsAffected is returned when an update/delete matched nothing.
	ErrNoRowsAffected = errors.New("no rows affected")
)
