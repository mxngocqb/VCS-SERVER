package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
	"github.com/redis/go-redis/v9"
)

var ctx context.Context = context.Background()

type serverCacheImpl struct {
	client  *redis.Client
	expires time.Duration
}



func NewServerRedisCache(client *redis.Client, expiration int) ServerCache {
	expires := time.Duration(expiration) * time.Second

	return &serverCacheImpl{
		client:  client,
		expires: expires,
	}
}

func (serverCache *serverCacheImpl) Delete(key string) error {
	redisKey := "server:" + key
	return serverCache.client.Del(ctx, redisKey).Err()
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
		log.Printf("Error getting data from Redis:", err)
		return nil
	} else{
		log.Printf("Cached data from Redis")
	}

	// Decode JSON
	var drivers []model.Server
	if err := json.Unmarshal([]byte(data), &drivers); err != nil {
		log.Printf("Error unmarshalling JSON:", err)
		return nil
	}

	return drivers
}

func (serverCache *serverCacheImpl) SetMultiRequest(key string, value []model.Server)  {
	// Encode slice of Driver objects to JSON
	jsonData, err := json.Marshal(value)
	if err != nil {
		log.Printf("Error marshalling JSON:", err)
		panic(err)
	}
	// Save to Redis
	err = serverCache.client.Set(ctx, key, jsonData, serverCache.expires).Err()
	if err != nil {
		log.Printf("Error setting data in Redis:", err)
		panic(err)
	} else{
		log.Printf("Data saved in Redis")
	}
}

func (serverCache *serverCacheImpl) GetTotalServer(key string) int64 {
	// Retrieve data from Redis
	numberOfServerStr, err := serverCache.client.Get(ctx, key).Result()
	if err != nil {
		log.Printf("Error getting data from Redis:", err)
		return -1
	} else{
		log.Printf("Cached total server from Redis")
	}

	numberOfServer, _ := strconv.ParseInt(numberOfServerStr, 10, 64)
	return numberOfServer
}

func (serverCache *serverCacheImpl) SetTotalServer(key string, value int64)  {
	// Encode slice of Driver objects to JSON
	jsonData, err := json.Marshal(value)
	if err != nil {
		log.Printf("Error marshalling JSON:", err)
		panic(err)
	}
	// Save to Redis
	err = serverCache.client.Set(ctx, key, jsonData, serverCache.expires).Err()
	if err != nil {
		log.Printf("Error setting data in Redis:", err)
		panic(err)
	} else{
		log.Printf("Total server saved in Redis")
	}
}

func (serverCache *serverCacheImpl) ConstructCacheKey(perPage, offset int, status, field, order string) string {
	return fmt.Sprintf("servers:%d:%d:%s:%s:%s", perPage, offset, status, field, order)
}

func (serverCache *serverCacheImpl) InvalidateCache() {
	// Retrieve all keys from Redis
	ctx := context.Background()
	keys, err := serverCache.client.Keys(ctx, "servers:*").Result()
	if err != nil {
		log.Printf("Error retrieving cache keys for servers: %v", err)
		return
	}
	// Delete all keys
	if len(keys) > 0 {
		_, delErr := serverCache.client.Del(ctx, keys...).Result()
		if delErr != nil {
			log.Printf("Error deleting cache keys: %v", delErr)
		}
	}
}
