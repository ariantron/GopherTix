package handlers

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gopher_tix/modules/authentication/models"
	"gopher_tix/modules/authentication/requests"
	"gopher_tix/modules/authentication/services"
	"log"
	"net"
)

type LoginHandler interface {
	Login(c *fiber.Ctx) error
	RegisterRoutes(router fiber.Router)
}

type loginHandler struct {
	loginService services.LoginService
	validator    *validator.Validate
}

func NewLoginHandler(loginService services.LoginService) LoginHandler {
	return &loginHandler{
		loginService: loginService,
		validator:    validator.New(),
	}
}

func (h *loginHandler) RegisterRoutes(router fiber.Router) {
	auth := router.Group("/auth")
	auth.Post("/login", h.Login)
}

func (h *loginHandler) Login(c *fiber.Ctx) error {
	var req requests.LoginRequest
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

	user, token, err := h.loginService.ValidateUserCredentials(c.Context(), req.Email, req.Password)
	if errors.Is(err, services.ErrInvalidCredentials) {
		if user == nil {
			log.Printf("Invalid login attempt from IP: %s", c.IP())
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid credentials",
			})
		}

		loginRecord := &models.Login{
			UserID:  user.ID,
			IP:      net.ParseIP(c.IP()),
			Succeed: false,
		}

		if err := h.loginService.CreateLoginRecord(c.Context(), loginRecord); err != nil {
			log.Printf("Failed to create login record: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create login record",
			})
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	} else if err != nil {
		log.Printf("Internal server error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	loginRecord := &models.Login{
		UserID:  user.ID,
		IP:      net.ParseIP(c.IP()),
		Succeed: true,
	}

	if err := h.loginService.CreateLoginRecord(c.Context(), loginRecord); err != nil {
		log.Printf("Failed to create login record: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create login record",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}
