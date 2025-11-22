package utils

import (
	"context"
	"encoding/json"
	"time"

	"bigdata-api/internal/database"
)

var (
	CacheHits   uint64
	CacheMisses uint64
)

// Get value from Redis cache
func CacheGet(key string) (string, error) {
	val, err := database.Redis.Get(context.Background(), key).Result()
	if err == nil {
		CacheHits++
	} else {
		CacheMisses++
	}
	return val, err
}

// Set value into Redis cache
func CacheSet(key string, value interface{}, ttl time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return database.Redis.Set(context.Background(), key, b, ttl).Err()
}
