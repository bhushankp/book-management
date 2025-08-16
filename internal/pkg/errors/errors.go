package errors

import "net/http"

type Code string

const (
	ErrNotFound     Code = "NOT_FOUND"
	ErrInvalidInput Code = "INVALID_INPUT"
	ErrConflict     Code = "CONFLICT"
	ErrUnauthorized Code = "UNAUTHORIZED"
	ErrForbidden    Code = "FORBIDDEN"
	ErrInternal     Code = "INTERNAL"
)

type AppError struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string { return e.Message }

func HTTPStatus(c Code) int {
	switch c {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrInvalidInput:
		return http.StatusBadRequest
	case ErrConflict:
		return http.StatusConflict
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrForbidden:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

func E(code Code, msg string, err error) *AppError {
	return &AppError{Code: code, Message: msg, Err: err}
}
