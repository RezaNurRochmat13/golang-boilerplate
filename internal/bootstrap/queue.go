package bootstrap

import (
	"golang-boilerplate-example/database"

	"github.com/hibiken/asynq"
)

func InitQueue() *asynq.Client {
	return asynq.NewClient(database.RedisOpt())
}
