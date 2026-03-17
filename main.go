package main

import (
	"golang-boilerplate-example/internal/bootstrap"
	"golang-boilerplate-example/internal/container"
	"golang-boilerplate-example/internal/logger"
	"golang-boilerplate-example/internal/middleware"
	"golang-boilerplate-example/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func main() {
	// Logger initialization
	logger.Init()
	defer logger.Sync()

	// Bootstrapping app
	db := bootstrap.InitDatabase()
	redisClient := bootstrap.InitRedis()
	queue := bootstrap.InitQueue()
	defer queue.Close()

	// Build container dependencies
	initContainer := container.Build(db, redisClient, queue)

	// Start a new Fiber App
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {

			code := fiber.StatusInternalServerError

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			logger.Log.Error("request failed",
				zap.String("method", c.Method()),
				zap.String("path", c.Path()),
				zap.Int("status", code),
				zap.String("error", err.Error()),
			)

			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(middleware.LoggerMiddleware())

	// Routes registration
	routes.RegisterAuthRoutes(app, initContainer.UserHandler)
	routes.RegisterNoteRoutes(app, initContainer.NoteHandler, redisClient)

	// Health check
	routes.RegisterHealthRoutes(app)

	// Listen on port 8081
	if err := app.Listen(":8081"); err != nil {
		log.Fatal(err)
	}
}
