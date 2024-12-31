package errors

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ErrorInterface interface {
	HandleError(ctx *fiber.Ctx) error
	Error() string
}

var (
	handlers = []ErrorInterface{
		&ErrNotFound{},
		&ErrValidation{},
		&ErrPermissionDenied{},
		&ErrInternalServer{},
		&ErrIncorrectParameter{},
		&ErrInvalidCredentials{},
		&ErrRequestFormat{},
	}
)

func HandleError(ctx *fiber.Ctx) error {
	err := ctx.Next()

	for _, handler := range handlers {
		if errors.As(err, &handler) {
			return handler.HandleError(ctx)
		}
	}

	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "Internal server error",
	})
}

func ParseAndValidateRequest(ctx *fiber.Ctx, req interface{}) error {
	if err := ctx.BodyParser(req); err != nil {
		return NewRequestFormatError()
	}
	if err := validator.New().Struct(req); err != nil {
		return NewValidationError(err)
	}
	return nil
}
