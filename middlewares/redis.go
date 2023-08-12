package cache

import "github.com/go-redis/redis/v8"

var RC *redis.Client

func SetupRedis() {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	RC = client
}
