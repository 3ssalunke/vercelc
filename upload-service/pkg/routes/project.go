package routes

import (
	"github.com/3ssalunke/vercelc/upload-service/pkg/controller"
	"github.com/3ssalunke/vercelc/upload-service/templates"
	"github.com/labstack/echo/v4"
)

type (
	project struct {
		controller.Controller
	}
)

func (c *project) Get(ctx echo.Context) error {
	page := controller.NewPage(ctx)
	page.Layout = templates.LayoutMain
	page.Name = templates.PageProject
	page.Metatags.Description = "Project Details"

	return c.RenderPage(ctx, page)
}
