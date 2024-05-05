package routes

import (
	"net/http"

	"github.com/3ssalunke/vercelc/shared/config"

	"github.com/3ssalunke/vercelc/upload-service/pkg/controller"
	"github.com/3ssalunke/vercelc/upload-service/pkg/services"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
)

const (
	routePageHome           = "home"
	routerPageCoding        = "coding"
	routerApiProject        = "project"
	routerOrchestratorStart = "start"
)

// BuildRouter builds the router
func BuildRouter(c *services.Container) {
	c.Web.Group("").Static(config.StaticPrefix, config.StaticDir)
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

	// Error handler
	err := errorHandler{Controller: ctr}
	c.Web.HTTPErrorHandler = err.Get

	pageRoutes(g, ctr)
}

func pageRoutes(g *echo.Group, ctr controller.Controller) {
	home := home{Controller: ctr}
	g.GET("/", home.Get).Name = routePageHome
}
