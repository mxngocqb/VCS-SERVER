package cache

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/config"
)

// ConnectRedis connects to Redis using provided configuration
func ConnectRedis(config *config.Config) (*redis.Client, int, error) {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.REDIS.Addr,
		Password: config.REDIS.Pass,
		DB:       0, // default DB
	})

	res, err := rdb.Ping(ctx).Result()

	if err != nil {
		log.Fatalf("There is an error while connecting to the Redis ", err)
		return nil, 0, err
	} else {
		log.Printf("Redis connected", res)
	}

	return rdb, config.REDIS.Expire, nil
}