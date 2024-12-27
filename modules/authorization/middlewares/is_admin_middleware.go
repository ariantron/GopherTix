package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gopher_tix/modules/authorization/services"
)

type IsAdmin struct {
	authorizeService services.AuthorizeService
}

func NewIsAdminMiddleware(authorizeService services.AuthorizeService) *IsAdmin {
	return &IsAdmin{
		authorizeService: authorizeService,
	}
}

func (m *IsAdmin) Handle(c *fiber.Ctx) error {
	currentUserID := c.Locals("user_id").(uuid.UUID)

	isAdmin, err := m.authorizeService.IsAdmin(c.Context(), currentUserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "unauthorized: admin access required",
		})
	}

	return c.Next()
}
