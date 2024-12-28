package errors

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type InternalServerError struct {
	Operation string
}

func (e *InternalServerError) Error() string {
	return fmt.Sprintf("Internal Server Error: %s", e.Operation)
}

func NewInternalServerError(operation string) *InternalServerError {
	return &InternalServerError{Operation: operation}
}

func (e *InternalServerError) HandleError(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": e.Error(),
	})
}
