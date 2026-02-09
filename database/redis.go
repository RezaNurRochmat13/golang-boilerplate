package database

import "github.com/hibiken/asynq"

func RedisOpt() asynq.RedisClientOpt {
	return asynq.RedisClientOpt{
		Addr: "localhost:6379",
	}
}
