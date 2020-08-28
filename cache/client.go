package cache

import (
	"log"
	"os"

	"github.com/go-redis/redis/v8"

	"github.com/theovidal/105chat/utils"
)

var Client *redis.Client

func OpenCache() {
	addr := utils.GenerateAddress("CACHE")
	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv("CACHE_PASSWORD"),
		DB:       0,
	})

	_, err := Client.Ping(Client.Context()).Result()
	if err != nil {
		log.Panicf("‼ Error while connecting to cache: %s", err.Error())
	}
}
