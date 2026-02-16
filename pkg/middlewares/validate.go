package middlewares

import (
	"gofiber-starterkit/app/shared"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

var validate = validator.New()

func ValidateBody(c fiber.Ctx, dest any) error {
	if err := c.Bind().Body(dest); err != nil {
		return shared.ErrBadRequest("Invalid request body")
	}

	if err := validate.Struct(dest); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return shared.ErrBadRequest(formatValidationError(validationErrors))
		}
		return shared.ErrBadRequest("Validation failed")
	}

	return nil
}

func formatValidationError(errors validator.ValidationErrors) string {
	if len(errors) == 0 {
		return "Validation failed"
	}
	fe := errors[0]
	switch fe.Tag() {
	case "required":
		return fe.Field() + " is required"
	case "email":
		return fe.Field() + " must be a valid email"
	case "min":
		return fe.Field() + " must be at least " + fe.Param() + " characters"
	case "max":
		return fe.Field() + " must be at most " + fe.Param() + " characters"
	case "url":
		return fe.Field() + " must be a valid URL"
	default:
		return fe.Field() + " is invalid"
	}
}
