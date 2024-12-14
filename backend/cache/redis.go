package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

// Initialize Redis client
func InitRedis(host string, port string) error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, "6379"),
		Password: "", 
		DB:       0,  
	})

	// Test connection
	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return nil
}

// Get cached data
func Get(ctx context.Context, key string, dest interface{}) error {
	val, err := redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return fmt.Errorf("key does not exist")
	} else if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

// Set data in cache
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return redisClient.Set(ctx, key, data, expiration).Err()
}

// Delete cached data
func Delete(ctx context.Context, key string) error {
	return redisClient.Del(ctx, key).Err()
}

// Clear cache by pattern
func ClearByPattern(ctx context.Context, pattern string) error {
	iter := redisClient.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		err := redisClient.Del(ctx, iter.Val()).Err()
		if err != nil {
			return err
		}
	}
	return iter.Err()
}
