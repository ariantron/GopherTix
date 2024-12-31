package errors

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type ErrPermissionDenied struct {
	Action string
}

func (e *ErrPermissionDenied) Error() string {
	if e.Action == "" {
		return "Permission denied for this action"
	}
	return fmt.Sprintf("permission denied for action: %s", e.Action)
}

func NewPermissionDeniedError(action string) *ErrPermissionDenied {
	return &ErrPermissionDenied{Action: action}
}

func (e *ErrPermissionDenied) HandleError(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"error": e.Error(),
	})
}
