package errors

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type ErrInternalServer struct {
	Operation string
}

func (e *ErrInternalServer) Error() string {
	return fmt.Sprintf("Internal Server Error: %s", e.Operation)
}

func NewInternalServerError(operation string) *ErrInternalServer {
	return &ErrInternalServer{Operation: operation}
}

func (e *ErrInternalServer) HandleError(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": e.Error(),
	})
}
