package routes

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/3ssalunke/vercelc/request-handler/pkg/controller"
	"github.com/labstack/echo/v4"
)

type serve struct {
	controller.Controller
}

func (s *serve) Get(ctx echo.Context) error {
	projectId := ctx.Param("projectId")
	path := ctx.Param("*")

	key := fmt.Sprintf("output/%s/dist/%s", projectId, path)
	content, err := s.Container.S3Storage.DownloadFile(key)
	if err != nil {
		log.Printf("error occured while getting file from s3: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	contentType := "text/plain" // Default content type
	if ext := filepath.Ext(key); ext != "" {
		contentType = mime.TypeByExtension(ext)
	}

	ctx.Response().Header().Set(echo.HeaderContentType, contentType)
	return ctx.Blob(http.StatusOK, "application/octet-stream", content)
}
