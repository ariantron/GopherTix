package errors

import (
	"github.com/gofiber/fiber/v2"
)

type IncorrectParameterError struct {
	Message string
}

func (e *IncorrectParameterError) Error() string {
	return e.Message
}

func NewIncorrectParameter(message string) *IncorrectParameterError {
	return &IncorrectParameterError{Message: message}
}

func (e *IncorrectParameterError) HandleError(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": e.Error(),
	})
}
