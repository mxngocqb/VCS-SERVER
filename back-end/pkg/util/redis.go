package util

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // or other appropriate address
		Password: "",               // no password set
		DB:       0,                // use default DB
	})

	_, err := RDB.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	} else{
		log.Printf("Connected to Redis")
	}
}
