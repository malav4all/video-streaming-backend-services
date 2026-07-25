package redis

import (
	"github.com/redis/go-redis/v9"

	"github.com/yourorg/go-user-service/internal/config"
)

// NewRedisClient opens a connection to Redis using the resolved config.
func NewRedisClient(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Username: cfg.RedisUsername,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
}