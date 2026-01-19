package tests

import (
	"golang-boilerplate-example/module/note"
	"golang-boilerplate-example/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupTestApp() *fiber.App {
	// sqlite in-memory
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// migrate schema
	if err := db.AutoMigrate(&note.Note{}); err != nil {
		log.Fatal(err)
	}

	// wire dependencies
	repo := note.NewRepository(db)
	service := note.NewService(repo)
	handler := note.NewHandler(service)

	app := fiber.New()
	routes.RegisterNoteRoutes(app, handler)

	return app
}
