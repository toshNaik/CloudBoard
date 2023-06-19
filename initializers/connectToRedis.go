package initializers

import (
	"os"

	redis "github.com/redis/go-redis/v9"
)
var Redis *redis.Client

func ConnectToRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:	  os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASS"),
		DB:		  0, 
	})
}