package services

import (
	"fmt"

	"github.com/3ssalunke/vercelc/shared/config"
	"github.com/3ssalunke/vercelc/shared/services"
)

type Container struct {
	Config    *config.Config
	S3Storage *services.S3Storage
	RedisConn *services.RedisConnections
}

func NewContainer() *Container {
	c := new(Container)
	c.initConfig()
	c.initS3Storage()
	c.initRedis()
	return c
}

func (c *Container) Shutdown() error {
	return nil
}

func (c *Container) initConfig() {
	cfg, err := config.GetConfig()
	if err != nil {
		panic(fmt.Sprintf("failied to load config: %v", err))
	}
	c.Config = &cfg
}

func (c *Container) initS3Storage() {
	storage, err := services.NewS3Storage(c.Config)
	if err != nil {
		panic(fmt.Sprintf("failed to create s3 client: %v", err))
	}
	c.S3Storage = storage
}

func (c *Container) initRedis() {
	conn, err := services.NewRedisConnection(c.Config)
	if err != nil {
		panic(fmt.Sprintf("failed to create redis connections %v", err))
	}
	c.RedisConn = conn
}
