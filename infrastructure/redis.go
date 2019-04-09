package infrastructure

import (
	"github.com/go-redis/redis"
	"os"
)

var Redis *redis.Client

func init() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
}
