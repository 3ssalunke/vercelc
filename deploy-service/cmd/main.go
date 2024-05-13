package main

import (
	"github.com/3ssalunke/vercelc/deploy-service/pkg/listener"
	"github.com/3ssalunke/vercelc/deploy-service/pkg/services"
)

func main() {
	c := services.NewContainer()
	defer func() {
		if err := c.Shutdown(); err != nil {
			panic(err)
		}
	}()

	listener := listener.NewListener(c)
	if err := listener.Start(); err != nil {
		panic(err)
	}
}
