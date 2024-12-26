package handlers

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gopher_tix/modules/authorization/requests"
	"gopher_tix/modules/authorization/services"
)

type AuthorizeHandler struct {
	service   *services.AuthorizeService
	validator *validator.Validate
}

func NewAuthorizeHandler(service *services.AuthorizeService) *AuthorizeHandler {
	return &AuthorizeHandler{
		service:   service,
		validator: validator.New(),
	}
}

func (h *AuthorizeHandler) RegisterRoutes(router fiber.Router) {
	auth := router.Group("/authorize")
	auth.Post("/assign-role", h.AssignRole)
	auth.Post("/unassign-role", h.UnassignRole)
}

func (h *AuthorizeHandler) AssignRole(c *fiber.Ctx) error {
	var req requests.RoleManagementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validator.Struct(req); err != nil {
		var validationErrors validator.ValidationErrors
		errors.As(err, &validationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": validationErrors.Error(),
		})
	}

	if err := h.service.AssignRole(c.Context(), req.UserID, req.GroupID, req.RoleID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to assign role",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Role assigned successfully",
	})
}

func (h *AuthorizeHandler) UnassignRole(c *fiber.Ctx) error {
	var req requests.RoleManagementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validator.Struct(req); err != nil {
		var validationErrors validator.ValidationErrors
		errors.As(err, &validationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": validationErrors.Error(),
		})
	}

	if err := h.service.UnassignRole(c.Context(), req.UserID, req.GroupID, req.RoleID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to unassign role",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Role unassigned successfully",
	})
}
