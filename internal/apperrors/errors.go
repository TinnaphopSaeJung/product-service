package apperrors

import "errors"

var (
	ErrValidation      = errors.New("validation error")
	ErrProductNotFound = errors.New("product not found")
)
