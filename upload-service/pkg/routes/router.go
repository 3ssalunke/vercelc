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
	routePageHome    = "home"
	routePageProject = "project"
	routerApiUpload  = "deploy"
	routerApiStatus  = "status"
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
	apiRoutes(g, ctr)
}

func pageRoutes(g *echo.Group, ctr controller.Controller) {
	home := home{Controller: ctr}
	g.GET("/", home.Get).Name = routePageHome

	project := project{Controller: ctr}
	g.GET("/project/:projectId", project.Get).Name = routePageProject
}

func apiRoutes(g *echo.Group, ctr controller.Controller) {
	deploy := deploy{Controller: ctr}
	g.POST("/deploy", deploy.Post).Name = routerApiUpload

	status := status{Controller: ctr}
	g.GET("/status", status.Get).Name = routerApiStatus
}
