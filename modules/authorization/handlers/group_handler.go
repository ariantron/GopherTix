package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gopher_tix/modules/authorization/models"
	"gopher_tix/modules/authorization/requests"
	"gopher_tix/modules/authorization/services"
	"gopher_tix/packages/common/errors"
	"gopher_tix/packages/utils"
	"strconv"
)

type GroupHandler interface {
	RegisterRoutes(router fiber.Router)
	List(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	GetMembers(c *fiber.Ctx) error
}

type groupHandler struct {
	service      services.GroupService
	validator    *validator.Validate
	errorHandler errors.ErrorHandler
}

func NewGroupHandler(service services.GroupService) GroupHandler {
	return &groupHandler{
		service:      service,
		validator:    validator.New(),
		errorHandler: errors.NewErrorHandler(),
	}
}

func (h *groupHandler) RegisterRoutes(router fiber.Router) {
	routes := router.Group("/groups")
	routes.Get("/", h.List)
	routes.Post("/", h.Create)
	routes.Get("/:id", h.GetByID)
	routes.Put("/:id", h.Update)
	routes.Delete("/:id", h.Delete)
	routes.Get("/:id/members", h.GetMembers)
}

func (h *groupHandler) List(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	search := c.Query("search")
	count, err := h.service.Count(c.Context(), &search)
	if err != nil {
		return h.handleServiceError(c, err)
	}
	totalPages, offset := utils.Paginate(count, page, limit)
	groups, err := h.service.List(c.Context(), offset, limit, &search)
	if err != nil {
		return h.handleServiceError(c, err)
	}
	return c.JSON(fiber.Map{
		"groups":     groups,
		"total":      count,
		"page":       page,
		"limit":      limit,
		"totalPages": totalPages,
	})
}

func (h *groupHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}
	group, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return h.handleServiceError(c, err)
	}
	return c.JSON(group)
}

func (h *groupHandler) Create(c *fiber.Ctx) error {
	var req requests.GroupCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	currentUserID := c.Locals("userID").(uuid.UUID)
	group := &models.Group{
		Name:          req.Name,
		ParentGroupID: req.ParentGroupID,
	}

	if err := h.service.Create(c.Context(), group, req.OwnerUserID, currentUserID); err != nil {
		return h.handleServiceError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(group)
}

func (h *groupHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	var req requests.GroupUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	currentUserID := c.Locals("userID").(uuid.UUID)
	group := &models.Group{
		Name: req.Name,
	}
	group.ID = id

	if err := h.service.Update(c.Context(), group, req.OwnerUserID, currentUserID); err != nil {
		return h.handleServiceError(c, err)
	}

	return c.JSON(group)
}

func (h *groupHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	currentUserID := c.Locals("userID").(uuid.UUID)
	if err := h.service.Delete(c.Context(), id, currentUserID); err != nil {
		return h.handleServiceError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *groupHandler) GetMembers(c *fiber.Ctx) error {
	groupID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid group ID format",
		})
	}

	currentUserID := c.Locals("userID").(uuid.UUID)
	members, err := h.service.GetMembers(c.Context(), groupID, currentUserID)
	if err != nil {
		return h.handleServiceError(c, err)
	}

	return c.JSON(members)
}

func (h *groupHandler) handleServiceError(c *fiber.Ctx, err error) error {
	errorResponse := h.errorHandler.HandleError(err)
	return c.Status(errorResponse.Status).JSON(errorResponse)
}
