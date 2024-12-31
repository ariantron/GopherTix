package errors

import (
	"github.com/gofiber/fiber/v2"
)

type ErrInvalidCredentials struct {
}

func (e *ErrInvalidCredentials) Error() string {
	return "Invalid credentials"
}

func NewInvalidCredentialsError() *ErrInvalidCredentials {
	return &ErrInvalidCredentials{}
}

func (e *ErrInvalidCredentials) HandleError(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": e.Error(),
	})
}
