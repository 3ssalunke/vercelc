package services

import (
	"fmt"

	"github.com/3ssalunke/vercelc/shared/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisConnection(config *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Hostname, config.Redis.Port),
		Password: "",
		DB:       0,
	})

	if rdb == nil {
		return nil, fmt.Errorf("Error: Redis client is nil")
	}

	return rdb, nil
}
