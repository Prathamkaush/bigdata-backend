package database

import (
	"context"
	"crypto/tls"
	"log"
	"net/url"

	"bigdata-api/internal/config"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client

func ConnectRedis(cfg *config.Config) {
	redisURL, _ := url.Parse(cfg.RedisURL)

	password := ""
	if redisURL.User != nil {
		password, _ = redisURL.User.Password()
	}

	Redis = redis.NewClient(&redis.Options{
		Addr:      redisURL.Host,
		Password:  password,
		DB:        0,
		TLSConfig: &tls.Config{}, // ✔ REQUIRED for rediss://
	})

	_, err := Redis.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("❌ Redis connection failed: %v", err)
	}

	log.Println("🟢 Connected to Upstash Redis")
}
