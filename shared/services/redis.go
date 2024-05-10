package services

import (
	"fmt"

	"github.com/3ssalunke/vercelc/shared/config"
	"github.com/redis/go-redis/v9"
)

type RedisConnections struct {
	Publisher  *redis.Client
	Subscriber *redis.Client
}

func NewRedisConnection(config *config.Config) (*RedisConnections, error) {
	pub := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Hostname, config.Redis.Port),
		Password: "",
		DB:       0,
	})

	if pub == nil {
		return nil, fmt.Errorf("publisher redis client is nil")
	}

	sub := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Hostname, config.Redis.Port),
		Password: "",
		DB:       0,
	})

	if sub == nil {
		return nil, fmt.Errorf("subscriber redis client is nil")
	}

	return &RedisConnections{
		Publisher:  pub,
		Subscriber: sub,
	}, nil
}
