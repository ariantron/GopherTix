package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gopher_tix/modules/authorization/services"
	errs "gopher_tix/packages/common/errors"
)

type IsAdmin struct {
	authorizeService services.AuthorizeService
}

func NewIsAdminMiddleware(authorizeService services.AuthorizeService) *IsAdmin {
	return &IsAdmin{
		authorizeService: authorizeService,
	}
}

func (m *IsAdmin) Handle(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("user_id").(string)
	currentUserID, err := uuid.Parse(userIDStr)
	if err != nil {
		return errs.NewValidationError(err)
	}

	isAdmin, err := m.authorizeService.IsAdmin(ctx.Context(), currentUserID)
	if err != nil {
		return errs.NewInternalServerError("Failed to checking that the user is an admin")
	}

	if !isAdmin {
		return errs.NewPermissionDeniedError("")
	}

	return ctx.Next()
}
