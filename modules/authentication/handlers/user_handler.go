package handlers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gopher_tix/modules/authentication/middlewares"
	"gopher_tix/modules/authentication/models"
	"gopher_tix/modules/authentication/requests"
	"gopher_tix/modules/authentication/services"
	"gopher_tix/packages/common/types"
	"gopher_tix/packages/utils"
	"log"
	"strconv"
)

type UserHandler struct {
	userService services.UserService
	validator   *validator.Validate
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validator.New(),
	}
}

func (h *UserHandler) RegisterRoutes(router fiber.Router) {
	users := router.Group("/users")
	users.Use(middlewares.JwtProtected())
	users.Post("/", h.CreateUser)
	users.Get("/:id", h.GetUser)
	users.Get("/", h.ListUsers)
	users.Put("/:id", h.UpdateUser)
	users.Delete("/:id", h.DeleteUser)
	users.Delete("/:id/soft", h.SoftDeleteUser)
	users.Post("/:id/restore", h.RestoreUser)
}

func (h *UserHandler) CreateUser(ctx *fiber.Ctx) error {
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

func (h *UserHandler) GetUser(ctx *fiber.Ctx) error {
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

func (h *UserHandler) ListUsers(ctx *fiber.Ctx) error {
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

func (h *UserHandler) UpdateUser(ctx *fiber.Ctx) error {
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

func (h *UserHandler) DeleteUser(ctx *fiber.Ctx) error {
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

func (h *UserHandler) SoftDeleteUser(ctx *fiber.Ctx) error {
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

func (h *UserHandler) RestoreUser(ctx *fiber.Ctx) error {
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
