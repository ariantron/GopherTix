package errors

import (
	"github.com/gofiber/fiber/v2"
)

type ErrIncorrectParameter struct {
	Message string
}

func (e *ErrIncorrectParameter) Error() string {
	return e.Message
}

func NewIncorrectParameter(message string) *ErrIncorrectParameter {
	return &ErrIncorrectParameter{Message: message}
}

func (e *ErrIncorrectParameter) HandleError(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": e.Error(),
	})
}
