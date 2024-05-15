package routes

import (
	"net/http"

	"github.com/3ssalunke/vercelc/request-handler/pkg/controller"
	"github.com/3ssalunke/vercelc/request-handler/pkg/services"
	echomw "github.com/labstack/echo/v4/middleware"
)

const routeServe = "serve"

func BuildRouter(c *services.Container) {
	// Non static file route group
	g := c.Web.Group("")

	g.Use(
		echomw.RemoveTrailingSlashWithConfig(echomw.TrailingSlashConfig{
			RedirectCode: http.StatusMovedPermanently,
		}),
		echomw.Recover(),
		echomw.Secure(),
		echomw.RequestID(),
		echomw.Gzip(),
		echomw.Logger(),
	)

	// Base controller
	ctr := controller.NewController(c)

	serve := serve{Controller: ctr}
	g.GET("/:projectId/*?", serve.Get).Name = routeServe
}
