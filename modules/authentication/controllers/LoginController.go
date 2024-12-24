package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gopher_tix/modules/authentication/models"
	"gopher_tix/modules/authentication/services"
	"net"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginController interface {
	Login(c *fiber.Ctx) error
}

type loginController struct {
	loginService services.LoginService
	validate     *validator.Validate
}

func NewLoginController(loginService services.LoginService) LoginController {
	return &loginController{
		loginService: loginService,
		validate:     validator.New(),
	}
}

func (ctrl *loginController) RegisterRoutes(router fiber.Router) {
	auth := router.Group("/auth")
	auth.Post("/", ctrl.Login)
}

func (ctrl *loginController) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := ctrl.validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": validationErrors.Error(),
		})
	}

	user, token, err := ctrl.loginService.ValidateUserCredentials(c.Context(), req.Email, req.Password)
	if err == services.ErrInvalidCredentials {
		loginRecord := &models.Login{
			UserID:  uuid.Nil,
			IP:      net.ParseIP(c.IP()),
			Succeed: false,
		}
		_ = ctrl.loginService.CreateLoginRecord(c.Context(), loginRecord)

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	loginRecord := &models.Login{
		UserID:  user.ID,
		IP:      net.ParseIP(c.IP()),
		Succeed: true,
	}
	if err := ctrl.loginService.CreateLoginRecord(c.Context(), loginRecord); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create login record",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":          user.ID,
			"email":       user.Email,
		},
	})
}