package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gopher_tix/configs"
	"gopher_tix/modules/authentication/controllers"
	"gopher_tix/modules/authentication/repositories"
	"gopher_tix/modules/authentication/services"
	"gorm.io/gorm"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func serve(db *gorm.DB) {
	app := fiber.New(fiber.Config{
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	})

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	userRepo := repositories.NewUserRepository(db)
	loginRepo := repositories.NewLoginRepository(db)

	userService := services.NewUserService(userRepo)
	loginService := services.NewLoginService(loginRepo, userService)

	userController := controllers.NewUserController(userService)
	loginController := controllers.NewLoginController(loginService)

	api := app.Group("/api")
	userController.RegisterRoutes(api)
	loginController.RegisterRoutes(api)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server...")

		if err := app.Shutdown(); err != nil {
			log.Fatalf("Server forced to shutdown: %v\n", err)
		}
	}()

	port := configs.AppPort

	log.Printf("Server starting on port %s...\n", port)
	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Server error: %v\n", err)
	}
}
