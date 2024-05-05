package routes

import (
	"net/http"

	"github.com/3ssalunke/vercelc/upload-service/pkg/context"
	"github.com/3ssalunke/vercelc/upload-service/pkg/controller"
	"github.com/labstack/echo/v4"
)

type errorHandler struct {
	controller.Controller
}

func (e *errorHandler) Get(err error, ctx echo.Context) {
	if ctx.Response().Committed || context.IsCanceledError(err) {
		return
	}

	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	if code >= 500 {
		ctx.Logger().Error(err)
	} else {
		ctx.Logger().Info(err)
	}
}
