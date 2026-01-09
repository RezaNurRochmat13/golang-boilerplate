package routes

import (
	"golang-boilerplate-example/module/note"

	"github.com/gofiber/fiber/v2"
)

func RegisterNoteRoutes(app *fiber.App, handler *note.Handler) {
	// Note routes
	group := app.Group("/api")
	group.Post("/notes", handler.CreateNote)
	group.Get("/notes", handler.GetAllNotes)
	group.Get("/notes/:id", handler.GetNote)
	group.Put("/notes/:id", handler.UpdateNote)
	group.Delete("/notes/:id", handler.DeleteNote)
}
