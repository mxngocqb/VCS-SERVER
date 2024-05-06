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

// ConstructCacheKey generates a unique key for storing/retrieving data from Redis
func ConstructCacheKey(perPage, offset int, status, field, order string) string {
	return fmt.Sprintf("servers:%d:%d:%s:%s:%s", perPage, offset, status, field, order)
}

func InvalidateCache() {
	// Retrieve all keys that match the pattern used in ConstructCacheKey
	ctx := context.Background()
	keys, err := RDB.Keys(ctx, "servers:*").Result()
	if err != nil {
		// Log the error instead of panicking to ensure service resilience
		log.Printf("Error retrieving cache keys for servers: %v", err)
		return
	}

	// Delete all found keys to invalidate the cache
	if len(keys) > 0 {
		_, delErr := RDB.Del(ctx, keys...).Result()
		if delErr != nil {
			// Log the error and continue
			log.Printf("Error deleting cache keys: %v", delErr)
		}
	}
}
