package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gopher_tix/modules/authentication/models"
	"gopher_tix/modules/authentication/services"
	"gopher_tix/packages/common/types"
	"strconv"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (ctrl *UserController) RegisterRoutes(router fiber.Router) {
	users := router.Group("/users")
	users.Post("/", ctrl.CreateUser)
	users.Get("/:id", ctrl.GetUser)
	users.Get("/", ctrl.ListUsers)
	users.Put("/:id", ctrl.UpdateUser)
	users.Delete("/:id", ctrl.DeleteUser)
	users.Delete("/:id/soft", ctrl.SoftDeleteUser)
	users.Post("/:id/restore", ctrl.RestoreUser)
}

func (ctrl *UserController) CreateUser(ctx *fiber.Ctx) error {
	user := new(models.User)
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := ctrl.userService.CreateUser(ctx.Context(), user); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(user)
}

func parseUUID(id string) (uuid.UUID, error) {
	return uuid.Parse(id)
}

func (ctrl *UserController) GetUser(ctx *fiber.Ctx) error {
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

	result, err := ctrl.userService.GetUserByID(ctx.Context(), user)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return ctx.JSON(result)
}

func (ctrl *UserController) ListUsers(ctx *fiber.Ctx) error {
	offset, _ := strconv.Atoi(ctx.Query("offset", "0"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))

	users, err := ctrl.userService.ListUsers(ctx.Context(), offset, limit)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	count, err := ctrl.userService.CountUsers(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"users": users,
		"total": count,
	})
}

func (ctrl *UserController) UpdateUser(ctx *fiber.Ctx) error {
	id, err := parseUUID(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid UUID format",
		})
	}

	user := new(models.User)
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user.ID = id
	if err := ctrl.userService.UpdateUser(ctx.Context(), user); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(user)
}

func (ctrl *UserController) DeleteUser(ctx *fiber.Ctx) error {
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

	if err := ctrl.userService.DeleteUser(ctx.Context(), user); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (ctrl *UserController) SoftDeleteUser(ctx *fiber.Ctx) error {
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

	if err := ctrl.userService.SoftDeleteUser(ctx.Context(), user); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (ctrl *UserController) RestoreUser(ctx *fiber.Ctx) error {
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

	if err := ctrl.userService.RestoreUser(ctx.Context(), user); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.SendStatus(fiber.StatusOK)
}
