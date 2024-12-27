package router

import (
	"github.com/gofiber/fiber/v2"
	"gopher_tix/modules/authentication/handlers"
	autnrepositories "gopher_tix/modules/authentication/repositories"
	autnservices "gopher_tix/modules/authentication/services"
	autzrepositories "gopher_tix/modules/authorization/repositories"
	autzservices "gopher_tix/modules/authorization/services"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, router fiber.Router) {
	userRepo := autnrepositories.NewUserRepository(db)
	loginRepo := autnrepositories.NewLoginRepository(db)
	authorizeRepo := autzrepositories.NewAuthorizeRepository(db)

	userService := autnservices.NewUserService(userRepo)
	loginService := autnservices.NewLoginService(loginRepo, userService)
	authorizeService := autzservices.NewAuthorizeService(authorizeRepo)

	userHandler := handlers.NewUserHandler(userService, authorizeService)
	loginHandler := handlers.NewLoginHandler(loginService)

	userHandler.RegisterRoutes(router)
	loginHandler.RegisterRoutes(router)
}
