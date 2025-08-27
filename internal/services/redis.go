package services

import (
	"context"
	"log"

	"flypro/internal/config"

	"github.com/redis/go-redis/v9"
)

func MustOpenRedis(cfg *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       0,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Printf("warning: redis not reachable: %v", err)
	}
	return rdb
}
