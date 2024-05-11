package cache

import (
	"context"
	"encoding/json"
	"log"
	"time"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
	"github.com/redis/go-redis/v9"
)

var ctx context.Context = context.Background()

type serverCacheImpl struct {
	client  *redis.Client
	expires time.Duration
}

func (serverCache *serverCacheImpl) Delete(key string) error {
	redisKey := "server:" + key
	return serverCache.client.Del(ctx, redisKey).Err()
}

func NewServerRedisCache(client *redis.Client, expiration int) ServerCache {
	expires := time.Duration(expiration) * time.Second

	return &serverCacheImpl{
		client:  client,
		expires: expires,
	}
}

func (serverCache *serverCacheImpl) Set(key string, value *model.Server) {
	json, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	redisKey := "server:" + key
	serverCache.client.Set(ctx, redisKey, json, serverCache.expires)
}

func (serverCache *serverCacheImpl) Get(key string) *model.Server {

	redisKey := "server:" + key
	val, err := serverCache.client.Get(ctx, redisKey).Result()
	if err != nil {
		return nil
	}

	driver := model.Server{}
	err = json.Unmarshal([]byte(val), &driver)
	if err != nil {
		panic(err)
	}
	return &driver
}

func (serverCache *serverCacheImpl) GetMultiRequest(key string) []model.Server {
	// Retrieve data from Redis
	data, err := serverCache.client.Get(ctx, key).Result()
	if err != nil {
		log.Println("Error getting data from Redis:", err)
		return nil
	}

	// Decode JSON
	var drivers []model.Server
	if err := json.Unmarshal([]byte(data), &drivers); err != nil {
		log.Println("Error unmarshalling JSON:", err)
		return nil
	}

	return drivers
}

func (serverCache *serverCacheImpl) SetMultiRequest(key string, value []model.Server)  {
	// Encode slice of Driver objects to JSON
	jsonData, err := json.Marshal(value)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		panic(err)
	}

	// Save to Redis
	err = serverCache.client.Set(ctx, key, jsonData, serverCache.expires).Err()
	if err != nil {
		log.Println("Error setting data in Redis:", err)
		panic(err)
	}
}
