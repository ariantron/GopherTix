package errors

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type NotFoundError struct {
	Object string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.Object)
}

func NewNotFoundError(object string) *NotFoundError {
	return &NotFoundError{Object: object}
}

func (e *NotFoundError) HandleError(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": e.Error(),
	})
}
