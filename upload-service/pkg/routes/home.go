package routes

import (
	"github.com/3ssalunke/vercelc/upload-service/pkg/controller"
	"github.com/3ssalunke/vercelc/upload-service/templates"
	"github.com/labstack/echo/v4"
)

type (
	home struct {
		controller.Controller
	}
)

func (c *home) Get(ctx echo.Context) error {
	page := controller.NewPage(ctx)
	page.Layout = templates.LayoutMain
	page.Name = templates.PageHome
	page.Metatags.Description = "Welcome to the Vercelc"

	return c.RenderPage(ctx, page)
}
