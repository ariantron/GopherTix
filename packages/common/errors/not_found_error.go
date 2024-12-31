package errors

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type ErrNotFound struct {
	Object string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("%s not found", e.Object)
}

func NewNotFoundError(object string) *ErrNotFound {
	return &ErrNotFound{Object: object}
}

func (e *ErrNotFound) HandleError(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": e.Error(),
	})
}
