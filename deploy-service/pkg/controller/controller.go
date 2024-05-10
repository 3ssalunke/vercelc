package controller

import (
	"github.com/3ssalunke/vercelc/deploy-service/pkg/services"
)

type Controller struct {
	Container *services.Container
}

func NewController(c *services.Container) Controller {
	return Controller{
		Container: c,
	}
}
