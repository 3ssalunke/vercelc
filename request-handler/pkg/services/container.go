package services

import (
	"fmt"

	"github.com/3ssalunke/vercelc/shared/config"
	"github.com/3ssalunke/vercelc/shared/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type Container struct {
	Config    *config.Config
	Web       *echo.Echo
	S3Storage *services.S3Storage
	Validator *Validator
}

func NewContainer() *Container {
	c := new(Container)
	c.initConfig()
	c.initS3Storage()
	c.initWeb()
	c.initValidator()
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

func (c *Container) initValidator() {
	c.Validator = NewValidator()
}

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
