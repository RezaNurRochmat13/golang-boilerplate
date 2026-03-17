package container

import (
	"golang-boilerplate-example/module/note"
	"golang-boilerplate-example/module/user"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Container struct {
	NoteHandler *note.Handler
	UserHandler *user.Handler
}

func Build(
	db *gorm.DB,
	redisClient *redis.Client,
	queue *asynq.Client,
) *Container {

	// Note
	noteRepo := note.NewRepository(db)
	noteService := note.NewService(noteRepo, redisClient)
	noteHandler := note.NewHandler(noteService, queue)

	// User
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo, redisClient)
	userHandler := user.NewHandler(userService, redisClient)

	return &Container{
		NoteHandler: noteHandler,
		UserHandler: userHandler,
	}
}
