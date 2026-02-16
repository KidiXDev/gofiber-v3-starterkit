package shared

import (
	"github.com/gofiber/fiber/v3"
)

type HTTPError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

func NewHTTPError(code int, message string) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: message,
	}
}

var (
	ErrBadRequest          = func(msg string) *HTTPError { return NewHTTPError(fiber.StatusBadRequest, msg) }
	ErrUnauthorized        = func(msg string) *HTTPError { return NewHTTPError(fiber.StatusUnauthorized, msg) }
	ErrForbidden           = func(msg string) *HTTPError { return NewHTTPError(fiber.StatusForbidden, msg) }
	ErrNotFound            = func(msg string) *HTTPError { return NewHTTPError(fiber.StatusNotFound, msg) }
	ErrConflict            = func(msg string) *HTTPError { return NewHTTPError(fiber.StatusConflict, msg) }
	ErrPaymentRequired     = func(msg string) *HTTPError { return NewHTTPError(fiber.StatusPaymentRequired, msg) }
	ErrInternalServerError = func(msg string) *HTTPError { return NewHTTPError(fiber.StatusInternalServerError, msg) }
	ErrUnprocessableEntity = func(msg string) *HTTPError { return NewHTTPError(fiber.StatusUnprocessableEntity, msg) }
)
