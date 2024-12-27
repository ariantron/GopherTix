package errors

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	apperrors "gopher_tix/packages/common/errors"
)

var (
	ErrGroupNotFound = errors.New("group not found")
	ErrPermission    = errors.New("permission denied")
	ErrValidation    = errors.New("validation error")
)

func NewGroupNotFoundError(id string) error {
	return apperrors.WrapError(fmt.Errorf("%w: group ID %s", ErrGroupNotFound, id), apperrors.WithStatus(fiber.StatusNotFound))
}

func NewPermissionError(operation string) error {
	return apperrors.WrapError(fmt.Errorf("%w: %s", ErrPermission, operation), apperrors.WithStatus(fiber.StatusForbidden))
}

func NewValidationError(msg string) error {
	return apperrors.WrapError(fmt.Errorf("%w: %s", ErrValidation, msg), apperrors.WithStatus(fiber.StatusBadRequest))
}
