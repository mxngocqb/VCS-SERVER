package config

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

// RedisConfiginterface redis config interface
type RedisConfiginterface interface {
	Address() string
	Password() string
	Expiration() int
}


// NewRedisConfig create redis instance
func NewRedisConfig(config *Config) *Redis {
	address := config.REDIS.Addr
	password := config.REDIS.Pass
	expiration := config.REDIS.Expire
	redis := &Redis{
		Addr: address,
		Pass: password,
		Expire: expiration,
	}
	return redis
}

// Address get redis address
func (redis *Redis) Address() string {
	return redis.Addr
}

// Password get redis password
func (redis *Redis) Password() string {
	return redis.Pass
}

func (redis *Redis) Expiration() int {
	return redis.Expire
}

// ConnectRedis connects to Redis using provided configuration
func ConnectRedis(config RedisConfiginterface) (*redis.Client, int, error) {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Address(),
		Password: config.Password(),
		DB:       0, // default DB
	})

	res, err := rdb.Ping(ctx).Result()

	if err != nil {
		log.Fatalf("There is an error while connecting to the Redis ", err)
		return nil, 0, err
	} else {
		log.Printf("Redis connected", res)
	}

	return rdb, config.Expiration(), nil
}
