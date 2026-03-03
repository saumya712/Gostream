package domain

import "errors"

var (
	ErrInvalidLogLevel    = errors.New("invalid log level")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	Erruseralreadyexits   = errors.New("user already exists")
)
