package routes

import (
	"golang-boilerplate-example/internal/auth"
	"golang-boilerplate-example/module/note"
	"golang-boilerplate-example/module/user"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func RegisterNoteRoutes(app *fiber.App, handler *note.Handler, redis *redis.Client) {
	// Note routes
	group := app.Group("/api")
	group.Use(auth.JWTMiddleware(redis))
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

func RegisterHealthRoutes(app *fiber.App) {
	group := app.Group("/api")

	group.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("API is up and running")
	})
}
