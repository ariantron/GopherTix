package errors

import (
	"github.com/gofiber/fiber/v2"
)

type RequestFormatError struct {
}

func (e *RequestFormatError) Error() string {
	return "The request body format is invalid"
}

func NewRequestFormatError() *RequestFormatError {
	return &RequestFormatError{}
}

func (e *RequestFormatError) HandleError(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": e.Error(),
	})
}
