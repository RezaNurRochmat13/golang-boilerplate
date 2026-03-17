package bootstrap

import "github.com/redis/go-redis/v9"

func InitRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}
