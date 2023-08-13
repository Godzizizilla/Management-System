package middlewares

import (
	"github.com/Godzizizilla/Management-System/config"
	"github.com/go-redis/redis/v8"
	"log"
)

var RC *redis.Client

func SetupRedis() {
	client := redis.NewClient(&redis.Options{
		Addr: config.C.Redis.Address,
	})
	RC = client
}

func CloseRedis() {
	if RC != nil {
		if err := RC.Close(); err != nil {
			log.Fatalf("Failed to close Redis connection: %v", err)
		}
	}
}
