package initializers

import (
	"os"

	redis "github.com/redis/go-redis/v9"
)
var Redis *redis.Client

func ConnectToRedis() {
	options, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		panic(err)
	}
	Redis = redis.NewClient(options)
}