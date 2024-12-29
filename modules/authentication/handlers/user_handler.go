package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	autnmiddlwares "gopher_tix/modules/authentication/middlewares"
	"gopher_tix/modules/authentication/models"
	"gopher_tix/modules/authentication/requests"
	autnservices "gopher_tix/modules/authentication/services"
	autzmiddlwares "gopher_tix/modules/authorization/middlewares"
	autzservices "gopher_tix/modules/authorization/services"
	errs "gopher_tix/packages/common/errors"
	"gopher_tix/packages/utils"
	"strconv"
)

type UserHandler interface {
	GetByID(ctx *fiber.Ctx) error
	List(ctx *fiber.Ctx) error
	Create(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
	Deactivate(ctx *fiber.Ctx) error
	Activate(ctx *fiber.Ctx) error
	RegisterRoutes(router fiber.Router)
}

type userHandler struct {
	userService      autnservices.UserService
	authorizeService autzservices.AuthorizeService
}

func NewUserHandler(userService autnservices.UserService, authorizeService autzservices.AuthorizeService) UserHandler {
	return &userHandler{
		userService:      userService,
		authorizeService: authorizeService,
	}
}

func (h *userHandler) RegisterRoutes(router fiber.Router) {
	routes := router.Group("/users")
	routes.Use(autnmiddlwares.JwtProtected()).
		Use(autzmiddlwares.NewIsAdminMiddleware(h.authorizeService).Handle)
	routes.Get("/:id", h.GetByID)
	routes.Get("/", h.List)
	routes.Post("/", h.Create)
	routes.Put("/:id", h.Update)
	routes.Delete("/:id", h.Delete)
	routes.Delete("/:id/deactivate", h.Deactivate)
	routes.Post("/:id/activate", h.Activate)
}

func (h *userHandler) Create(ctx *fiber.Ctx) error {
	req := new(requests.UserCreateRequest)
	if err := errs.ParseAndValidateRequest(ctx, &req); err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := &models.User{
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := h.userService.Create(ctx.Context(), user); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"user": user})
}

func (h *userHandler) GetByID(ctx *fiber.Ctx) error {
	id, err := utils.GetUUIDParam(ctx, "id")
	if err != nil {
		return err
	}

	user, err := h.userService.GetByID(ctx.Context(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{"user": user})
}

func (h *userHandler) List(ctx *fiber.Ctx) error {
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	search := ctx.Query("search")

	count, err := h.userService.Count(ctx.Context(), &search)
	if err != nil {
		return err
	}

	totalPages, offset := utils.Paginate(count, page, limit)
	users, err := h.userService.List(ctx.Context(), offset, limit, &search)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"users":      users,
		"total":      count,
		"page":       page,
		"limit":      limit,
		"totalPages": totalPages,
	})
}

func (h *userHandler) Update(ctx *fiber.Ctx) error {
	id, err := utils.GetUUIDParam(ctx, "id")
	if err != nil {
		return err
	}

	req := new(requests.UserUpdateRequest)
	if err := errs.ParseAndValidateRequest(ctx, &req); err != nil {
		return err
	}

	user, err := h.userService.GetByID(ctx.Context(), id)
	if err != nil {
		return err
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Password != "" {
		user.Password, err = utils.HashPassword(req.Password)
		if err != nil {
			return err
		}
	}

	if err := h.userService.Update(ctx.Context(), user); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{"user": user})
}

func (h *userHandler) Delete(ctx *fiber.Ctx) error {
	id, err := utils.GetUUIDParam(ctx, "id")
	if err != nil {
		return err
	}

	user, err := h.userService.GetByID(ctx.Context(), id)
	if err != nil {
		return err
	}

	if err := h.userService.Delete(ctx.Context(), user); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "The user has been successfully deleted"})
}

func (h *userHandler) Deactivate(ctx *fiber.Ctx) error {
	return h.toggleUserStatus(ctx, false, "deactivated")
}

func (h *userHandler) Activate(ctx *fiber.Ctx) error {
	return h.toggleUserStatus(ctx, true, "activated")
}

func (h *userHandler) toggleUserStatus(ctx *fiber.Ctx, activate bool, action string) error {
	id, err := utils.GetUUIDParam(ctx, "id")
	if err != nil {
		return err
	}

	user, err := h.userService.GetByID(ctx.Context(), id)
	if err != nil {
		return err
	}

	if activate {
		err = h.userService.Activate(ctx.Context(), user)
	} else {
		err = h.userService.Deactivate(ctx.Context(), user)
	}

	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": fmt.Sprintf("The user has been successfully %s", action)})
}
