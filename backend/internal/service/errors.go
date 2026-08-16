package service

import "errors"

var (
	// ErrInvalidCredentials indicates a failed login.
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrEmailTaken indicates a registration conflict.
	ErrEmailTaken = errors.New("email already registered")
	// ErrUsernameTaken indicates a registration conflict.
	ErrUsernameTaken = errors.New("username already taken")
	// ErrForbidden indicates the current user lacks permission.
	ErrForbidden = errors.New("forbidden")
	// ErrInvalidMember indicates the invite target cannot be found.
	ErrInvalidMember = errors.New("member not found")
)
