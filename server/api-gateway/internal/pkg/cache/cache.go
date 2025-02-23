package cache

import (
	"fmt"

	"github.com/LeonLow97/internal/config"
	"github.com/go-redis/redis/v8"
)

type Cache struct {
	RedisClient *redis.Client
	cfg         config.Config
}

func NewRedisClient(cfg config.Config) Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisServer.Address, cfg.RedisServer.Port),
		Password: cfg.RedisServer.Password,
		DB:       cfg.RedisServer.DatabaseIndex,
	})

	return Cache{
		RedisClient: client,
		cfg:         cfg,
	}
}
