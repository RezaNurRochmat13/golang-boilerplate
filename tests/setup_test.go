package tests

import (
	"golang-boilerplate-example/module/note"
	"golang-boilerplate-example/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
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

	// setup queue
	queue := asynq.NewClient(asynq.RedisClientOpt{Addr: ":6379"})
	defer queue.Close()

	// wire dependencies
	repo := note.NewRepository(db)
	service := note.NewService(repo)
	handler := note.NewHandler(service, queue)

	app := fiber.New()
	routes.RegisterNoteRoutes(app, handler)

	return app
}
