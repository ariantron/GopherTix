package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gopher_tix/modules/authorization/models"
	"gopher_tix/modules/authorization/requests"
	"gopher_tix/modules/authorization/services"
	"gopher_tix/packages/utils"
	"strconv"
)

type GroupHandler struct {
	service   *services.GroupService
	validator *validator.Validate
}

func NewGroupHandler(service *services.GroupService) *GroupHandler {
	return &GroupHandler{
		service:   service,
		validator: validator.New(),
	}
}

func (h *GroupHandler) RegisterRoutes(router fiber.Router) {
	groups := router.Group("/groups")
	groups.Get("/", h.List)
	groups.Post("/", h.Create)
	groups.Get("/:id", h.GetByID)
	groups.Put("/:id", h.Update)
	groups.Delete("/:id", h.Delete)
	groups.Get("/:id/members", h.GetMembers)
}

func (h *GroupHandler) List(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	search := c.Query("search")
	count, err := h.service.Count(c.Context(), &search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	totalPages, offset := utils.Paginate(count, page, limit)
	groups, err := h.service.List(c.Context(), offset, limit, &search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"groups":     groups,
		"total":      count,
		"page":       page,
		"limit":      limit,
		"totalPages": totalPages,
	})
}

func (h *GroupHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}
	group, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Group not found",
		})
	}
	return c.JSON(group)
}

func (h *GroupHandler) Create(c *fiber.Ctx) error {
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(group)
}

func (h *GroupHandler) Update(c *fiber.Ctx) error {
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(group)
}

func (h *GroupHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}
	currentUserID := c.Locals("userID").(uuid.UUID)
	if err := h.service.Delete(c.Context(), id, currentUserID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *GroupHandler) GetMembers(c *fiber.Ctx) error {
	groupID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid group ID format",
		})
	}
	currentUserID := c.Locals("userID").(uuid.UUID)
	members, err := h.service.GetMembers(c.Context(), groupID, currentUserID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(members)
}
