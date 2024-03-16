package errors

import (
	"errors"
	"net/http"
)

var (
	ErrStatusForbidden = errors.New(http.StatusText(http.StatusForbidden))
	ErrMovieCreation   = errors.New("movie creation error")
	ErrUserCreation    = errors.New("user creation error")
	ErrDatabase        = errors.New("database error")
	ErrAuthentication  = errors.New("authentication error")
)
