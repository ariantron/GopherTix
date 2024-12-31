package errors

import (
	"github.com/gofiber/fiber/v2"
)

type ErrRequestFormat struct {
}

func (e *ErrRequestFormat) Error() string {
	return "The request body format is invalid"
}

func NewRequestFormatError() *ErrRequestFormat {
	return &ErrRequestFormat{}
}

func (e *ErrRequestFormat) HandleError(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": e.Error(),
	})
}
