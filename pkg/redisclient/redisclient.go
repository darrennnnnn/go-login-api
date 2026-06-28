package redisclient

import (
	"context"
	"fmt"
	"strconv"

	"github.com/darrennnnnn/go-login-api/config"
	"github.com/redis/go-redis/v9"
)

func Connect(cfg *config.Config) *redis.Client {
	db, err := strconv.Atoi(cfg.Redis.DB)
	if err != nil {
		panic(fmt.Errorf("invalid REDIS_DB: %w", err))
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       db,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return client
}
