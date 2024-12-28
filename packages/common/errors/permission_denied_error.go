package errors

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type PermissionDeniedError struct {
	Action string
}

func (e *PermissionDeniedError) Error() string {
	if e.Action == "" {
		return "Permission denied for this action"
	}
	return fmt.Sprintf("permission denied for action: %s", e.Action)
}

func NewPermissionDeniedError(action string) *PermissionDeniedError {
	return &PermissionDeniedError{Action: action}
}

func (e *PermissionDeniedError) HandleError(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"error": e.Error(),
	})
}
