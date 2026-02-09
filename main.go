package main

import (
	"golang-boilerplate-example/database"
	"golang-boilerplate-example/module/note"
	"golang-boilerplate-example/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
)

func main() {
	// Start a new Fiber App
	app := fiber.New()

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

	// Note resources
	noteRepo := note.NewRepository(db)
	noteService := note.NewService(noteRepo)
	noteHandler := note.NewHandler(noteService, queue)

	routes.RegisterNoteRoutes(app, noteHandler)

	// Send string back for GET calls to the endpoint '/'
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("API is up and running")
	})

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	// Listen on port 8080
	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}

}
