package routes

import (
	"golang-boilerplate-example/internal/auth"
	"golang-boilerplate-example/module/note"
	"golang-boilerplate-example/module/user"

	"github.com/gofiber/fiber/v2"
)

func RegisterNoteRoutes(app *fiber.App, handler *note.Handler) {
	// Note routes
	group := app.Group("/api")
	group.Use(auth.JWTMiddleware())
	group.Post("/notes", handler.CreateNote)
	group.Get("/notes", handler.GetAllNotes)
	group.Get("/notes/:id", handler.GetNote)
	group.Put("/notes/:id", handler.UpdateNote)
	group.Delete("/notes/:id", handler.DeleteNote)
	group.Post("/send-email", handler.SendEmail)
}

func RegisterAuthRoutes(app *fiber.App, handler *user.Handler) {
	group := app.Group("/api/auth")
	group.Post("/register", handler.Register)
	group.Post("/login", handler.Login)
}
