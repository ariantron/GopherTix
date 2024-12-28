package errors

import (
	"github.com/gofiber/fiber/v2"
)

type InvalidCredentialsError struct {
}

func (e *InvalidCredentialsError) Error() string {
	return "Invalid credentials"
}

func NewInvalidCredentialsError() *InvalidCredentialsError {
	return &InvalidCredentialsError{}
}

func (e *InvalidCredentialsError) HandleError(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": e.Error(),
	})
}
