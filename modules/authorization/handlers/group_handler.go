package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gopher_tix/modules/authorization/models"
	"gopher_tix/modules/authorization/requests"
	"gopher_tix/modules/authorization/services"
	errs "gopher_tix/packages/common/errors"
	"gopher_tix/packages/utils"
	"strconv"
)

type GroupHandler interface {
	List(ctx *fiber.Ctx) error
	GetByID(ctx *fiber.Ctx) error
	Create(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
	GetMembers(ctx *fiber.Ctx) error
	RegisterRoutes(router fiber.Router)
}

type groupHandler struct {
	service services.GroupService
}

func NewGroupHandler(service services.GroupService) GroupHandler {
	return &groupHandler{
		service: service,
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

func (h *groupHandler) List(ctx *fiber.Ctx) error {
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	search := ctx.Query("search")
	count, err := h.service.Count(ctx.Context(), &search)
	if err != nil {
		return err
	}
	totalPages, offset := utils.Paginate(count, page, limit)
	groups, err := h.service.List(ctx.Context(), offset, limit, &search)
	if err != nil {
		return err
	}
	return ctx.JSON(fiber.Map{
		"groups":     groups,
		"total":      count,
		"page":       page,
		"limit":      limit,
		"totalPages": totalPages,
	})
}

func (h *groupHandler) GetByID(ctx *fiber.Ctx) error {
	id, err := utils.GetUUIDParam(ctx, "id")
	if err != nil {
		return err
	}

	group, err := h.service.GetByID(ctx.Context(), id)
	if err != nil {
		return err
	}
	return ctx.JSON(group)
}

func (h *groupHandler) Create(ctx *fiber.Ctx) error {
	var req requests.GroupCreateRequest
	if err := errs.ParseAndValidateRequest(ctx, &req); err != nil {
		return err
	}

	group := &models.Group{
		Name:          req.Name,
		ParentGroupID: req.ParentGroupID,
	}

	if err := h.service.Create(ctx.Context(), group, req.OwnerUserID, utils.CurrentUserID(ctx)); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(group)
}

func (h *groupHandler) Update(ctx *fiber.Ctx) error {
	id, err := utils.GetUUIDParam(ctx, "id")
	if err != nil {
		return err
	}

	req := new(requests.GroupUpdateRequest)
	if err := errs.ParseAndValidateRequest(ctx, &req); err != nil {
		return err
	}

	group := &models.Group{
		Name: req.Name,
	}
	group.ID = id

	if err := h.service.Update(ctx.Context(), group, req.OwnerUserID, utils.CurrentUserID(ctx)); err != nil {
		return err
	}

	return ctx.JSON(group)
}

func (h *groupHandler) Delete(ctx *fiber.Ctx) error {
	id, err := utils.GetUUIDParam(ctx, "id")
	if err != nil {
		return err
	}

	if err := h.service.Delete(ctx.Context(), id, utils.CurrentUserID(ctx)); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *groupHandler) GetMembers(ctx *fiber.Ctx) error {
	id, err := utils.GetUUIDParam(ctx, "id")
	if err != nil {
		return err
	}

	members, err := h.service.GetMembers(ctx.Context(), id, utils.CurrentUserID(ctx))
	if err != nil {
		return err
	}

	return ctx.JSON(members)
}
