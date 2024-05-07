package services

import (
	"fmt"

	"github.com/3ssalunke/vercelc/shared/config"
	"github.com/go-redis/redis/v8"
)

func NewRedisConnection(config *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Hostname, config.Redis.Port),
		Password: "",
		DB:       0,
	})

	return rdb
}
