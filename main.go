package main

import (
	"golang-boilerplate-example/database"
	"golang-boilerplate-example/module/note"
	"golang-boilerplate-example/routes"
	"log"

	"github.com/gofiber/fiber/v2"
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

	// Note resources
	noteRepo := note.NewRepository(db)
	noteService := note.NewService(noteRepo)
	noteHandler := note.NewHandler(noteService)

	routes.RegisterNoteRoutes(app, noteHandler)

	// Send string back for GET calls to the endpoint '/'
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("API is up and running")
	})

	// Listen on port 3000
	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}

}
