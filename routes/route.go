package routes

import (
	"golang-boilerplate-example/module/note"

	"github.com/gofiber/fiber/v2"
)

func RegisterNoteRoutes(app *fiber.App, handler *note.Handler) {
	// Note routes
	app.Post("/notes", handler.CreateNote)
	app.Get("/notes", handler.GetAllNotes)
	app.Get("/notes/:id", handler.GetNote)
	app.Put("/notes/:id", handler.UpdateNote)
	app.Delete("/notes/:id", handler.DeleteNote)
}
