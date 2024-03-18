package validation

import "errors"

var (
	ErrInvalidDescription = errors.New("description is too long")
	ErrInvalidRate        = errors.New("rate should be between 0 and 10")
	ErrInvalidName        = errors.New("name lenght should be between 1 and 150")
)
