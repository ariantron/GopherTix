package errors

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ErrValidation struct {
	Errors map[string]string
}

func (e *ErrValidation) Error() string {
	return "validation failed"
}

func NewValidationError(err error) *ErrValidation {
	var validationErrors validator.ValidationErrors
	errors.As(err, &validationErrors)
	errorsMap := make(map[string]string)
	for _, fieldErr := range validationErrors {
		errorsMap[fieldErr.Field()] = fmt.Sprintf("%s failed on the '%s' tag", fieldErr.Field(), fieldErr.Tag())
	}
	return &ErrValidation{Errors: errorsMap}
}

func (e *ErrValidation) HandleError(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error":   e.Error(),
		"details": e.Errors,
	})
}
