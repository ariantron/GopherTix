package server

import (
	"fmt"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gopher_tix/configs"
	autnrouter "gopher_tix/modules/authentication/router"
	autzrouter "gopher_tix/modules/authorization/router"
	tcktrouter "gopher_tix/modules/ticketing/router"
	errorhandler "gopher_tix/packages/common/errors"
	"gorm.io/gorm"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Serve(db *gorm.DB) {
	app := fiber.New(fiber.Config{
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	})

	registerMiddlewares(app)
	api := app.Group("/api")
	registerRoutes(db, api)

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

func registerRoutes(db *gorm.DB, router fiber.Router) {
	autnrouter.RegisterRoutes(db, router)
	autzrouter.RegisterRoutes(db, router)
	tcktrouter.RegisterRoutes(db, router)
}

func registerMiddlewares(app *fiber.App) {
	app.Use(logger.New())
	app.Use(recover.New())
	registerSwaggerMiddleware(app)
	registerCorsMiddleware(app)
	app.Use(errorhandler.HandleError)
}

func registerCorsMiddleware(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: configs.CorsAllowOrigins,
		AllowHeaders: configs.CorsAllowHeaders,
		AllowMethods: configs.CorsAllowMethods,
	}))
}

func registerSwaggerMiddleware(app *fiber.App) {
	app.Use(swagger.New(swagger.Config{
		BasePath: "/",
		FilePath: "./docs/swagger.json",
		Path:     "swagger",
		Title:    "Swagger API Docs",
	}))
}
