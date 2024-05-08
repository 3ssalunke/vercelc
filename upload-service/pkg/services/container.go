package services

import (
	"fmt"

	"github.com/3ssalunke/vercelc/shared/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
)

type RedisConnections struct {
	Publisher  *redis.Client
	Subscriber *redis.Client
}

// Container contains all services used by the application and provides an easy way to handle dependency
// injection including within tests
type Container struct {
	// Validator stores a validator
	Validator *Validator

	// Web stores the web framework
	Web *echo.Echo

	// Config stores the application configuration
	Config *config.Config

	// TemplateRenderer stores a service to easily render and cache templates
	TemplateRenderer *TemplateRenderer

	S3Storage *S3Storage

	RedisConn *RedisConnections
}

// NewContainer creates and initializes a new Container
func NewContainer() *Container {
	c := new(Container)
	c.initConfig()
	c.initValidator()
	c.initWeb()
	c.initTemplateRenderer()
	c.initS3Storage()
	c.initRedis()
	return c
}

// Shutdown shuts the Container down and disconnects all connections
func (c *Container) Shutdown() error {
	return nil
}

// initConfig initializes configuration
func (c *Container) initConfig() {
	cfg, err := config.GetConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	c.Config = &cfg
}

// initValidator initializes the validator
func (c *Container) initValidator() {
	c.Validator = NewValidator()
}

// initWeb initializes the web framework
func (c *Container) initWeb() {
	c.Web = echo.New()

	// Configure logging
	switch c.Config.App.Environment {
	case config.EnvProduction:
		c.Web.Logger.SetLevel(log.WARN)
	default:
		c.Web.Logger.SetLevel(log.DEBUG)
	}

	c.Web.Validator = c.Validator
}

// initTemplateRenderer initializes the template renderer
func (c *Container) initTemplateRenderer() {
	c.TemplateRenderer = NewTemplateRenderer(c.Config)
}

func (c *Container) initS3Storage() {
	storage, err := NewS3Storage(c.Config)
	if err != nil {
		panic(fmt.Sprintf("failed to create s3 client: %v", err))
	}
	c.S3Storage = storage
}

func (c *Container) initRedis() {
	pub, err := NewRedisConnection(c.Config)
	if err != nil {
		panic(fmt.Sprintf("failed to create publisher redis client %v", err))
	}
	sub, err := NewRedisConnection(c.Config)
	if err != nil {
		panic(fmt.Sprintf("failed to create subscriber redis client %v", err))
	}
	c.RedisConn = &RedisConnections{
		Publisher:  pub,
		Subscriber: sub,
	}
}
