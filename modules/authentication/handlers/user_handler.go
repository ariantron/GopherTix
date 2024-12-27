package handlers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	autnmiddlwares "gopher_tix/modules/authentication/middlewares"
	"gopher_tix/modules/authentication/models"
	"gopher_tix/modules/authentication/requests"
	autnservices "gopher_tix/modules/authentication/services"
	autzmiddlwares "gopher_tix/modules/authorization/middlewares"
	autzservices "gopher_tix/modules/authorization/services"
	"gopher_tix/packages/common/types"
	"gopher_tix/packages/utils"
	"log"
	"strconv"
)

type UserHandler interface {
	RegisterRoutes(router fiber.Router)
	GetByID(ctx *fiber.Ctx) error
	List(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
	Deactivate(ctx *fiber.Ctx) error
	Activate(ctx *fiber.Ctx) error
}

type userHandler struct {
	userService      autnservices.UserService
	authorizeService autzservices.AuthorizeService
	validator        *validator.Validate
}

func NewUserHandler(
	userService autnservices.UserService,
	authorizeService autzservices.AuthorizeService,
) UserHandler {
	return &userHandler{
		userService:      userService,
		authorizeService: authorizeService,
		validator:        validator.New(),
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
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Printf("Failed to hash password")
	}
	user := &models.User{
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := h.userService.CreateUser(ctx.Context(), user); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(user)
}

func parseUUID(id string) (uuid.UUID, error) {
	return uuid.Parse(id)
}

func (h *userHandler) GetByID(ctx *fiber.Ctx) error {
	id, err := parseUUID(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid UUID format",
		})
	}

	user := &models.User{
		SoftDeleteModel: types.SoftDeleteModel{
			BaseModel: types.BaseModel{
				ID: id,
			},
		},
	}

	result, err := h.userService.GetUserByID(ctx.Context(), user)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return ctx.JSON(result)
}

func (h *userHandler) List(ctx *fiber.Ctx) error {
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	search := ctx.Query(`search`)
	count, err := h.userService.CountUsers(ctx.Context(), &search)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	totalPages, offset := utils.Paginate(count, page, limit)
	users, err := h.userService.ListUsers(ctx.Context(), offset, limit, &search)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
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
	id, err := parseUUID(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid UUID format",
		})
	}

	req := new(requests.UserUpdateRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	user := &models.User{
		SoftDeleteModel: types.SoftDeleteModel{
			BaseModel: types.BaseModel{
				ID: id,
			},
		},
		Email: req.Email,
	}

	if req.Password != "" {
		user.Password, err = utils.HashPassword(req.Password)
	}

	if err := h.userService.UpdateUser(ctx.Context(), user); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(user)
}

func (h *userHandler) Delete(ctx *fiber.Ctx) error {
	id, err := parseUUID(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid UUID format",
		})
	}

	user := &models.User{
		SoftDeleteModel: types.SoftDeleteModel{
			BaseModel: types.BaseModel{
				ID: id,
			},
		},
	}

	if err := h.userService.DeleteUser(ctx.Context(), user); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("User was deleted with id %s", user.ID),
	})
}

func (h *userHandler) Deactivate(ctx *fiber.Ctx) error {
	id, err := parseUUID(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid UUID format",
		})
	}

	user := &models.User{
		SoftDeleteModel: types.SoftDeleteModel{
			BaseModel: types.BaseModel{
				ID: id,
			},
		},
	}

	if err := h.userService.SoftDeleteUser(ctx.Context(), user); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("User was deactivated with id %s", user.ID),
	})
}

func (h *userHandler) Activate(ctx *fiber.Ctx) error {
	id, err := parseUUID(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid UUID format",
		})
	}

	user := &models.User{
		SoftDeleteModel: types.SoftDeleteModel{
			BaseModel: types.BaseModel{
				ID: id,
			},
		},
	}

	if err := h.userService.RestoreUser(ctx.Context(), user); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("User was activated with id %s", user.ID),
	})
}
