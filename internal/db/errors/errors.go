package errors

import "errors"

var (
	ErrUserCreation   = errors.New("user creation error")
	ErrDatabase       = errors.New("database error")
	ErrAuthentication = errors.New("authentication error")
)
