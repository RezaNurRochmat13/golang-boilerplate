package main

import (
	"golang-boilerplate-example/database"
	"golang-boilerplate-example/internal/logger"
	"golang-boilerplate-example/internal/middleware"
	"golang-boilerplate-example/module/note"
	"golang-boilerplate-example/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	// Start logger
	logger.Init()
	defer logger.Sync()

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

	// Connect to the database
	db, err := database.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	// Auto migrate
	db.AutoMigrate(&note.Note{})

	// Initiate Asynq Client
	queue := asynq.NewClient(database.RedisOpt())
	defer queue.Close()

	// Initiate redis connection
	redisConnection := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Note resources
	noteRepo := note.NewRepository(db)
	noteService := note.NewService(noteRepo, redisConnection)
	noteHandler := note.NewHandler(noteService, queue)

	routes.RegisterNoteRoutes(app, noteHandler)

	// Send string back for GET calls to the endpoint '/'
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("API is up and running")
	})

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	// Listen on port 8081
	if err := app.Listen(":8081"); err != nil {
		log.Fatal(err)
	}

}
