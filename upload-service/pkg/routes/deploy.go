package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/3ssalunke/vercelc/shared/util"
	"github.com/3ssalunke/vercelc/upload-service/pkg/controller"
	"github.com/labstack/echo/v4"
)

type (
	deploy struct {
		controller.Controller
	}

	deployRequestPayload struct {
		RepoUrl     string `json:"repoUrl" validate:"required"`
		ProjectName string `json:"projectName" validate:"required"`
		Framework   string `json:"framework" validate:"required"`
	}

	deployResponsePayload struct {
		Id    string `json:"msg"`
		Error string `json:"error"`
	}
)

func (c *deploy) Post(ctx echo.Context) error {
	payload := new(deployRequestPayload)
	if err := ctx.Bind(payload); err != nil {
		log.Printf("error parsing request payload: %v", err)
		return ctx.JSON(http.StatusBadRequest, &deployResponsePayload{Error: err.Error()})
	}
	if err := c.Container.Validator.Validate(payload); err != nil {
		log.Printf("error validating request payload: %v", err)
		return ctx.JSON(http.StatusBadRequest, &deployResponsePayload{Error: err.Error()})
	}
	projectId := util.GenerateRandomId()
	destinationFolder, err := util.GetPathForFolder(fmt.Sprintf("output/%s", projectId))
	if err != nil {
		log.Printf("error getting destination folder for cloning git repo: %v", err)
		return ctx.JSON(http.StatusInternalServerError, &deployResponsePayload{Error: err.Error()})
	}
	if err := util.CloneRepo(payload.RepoUrl, destinationFolder); err != nil {
		log.Printf("error cloning git repo: %v", err)
		return ctx.JSON(http.StatusInternalServerError, &deployResponsePayload{Error: err.Error()})
	}
	if err := c.Container.S3Storage.CopyFolder(fmt.Sprintf("output/%s", projectId)); err != nil {
		log.Printf("error copying repo to s3: %v", err)
		return ctx.JSON(http.StatusInternalServerError, &deployResponsePayload{Error: err.Error()})
	}
	_, err = c.Container.RedisConn.Publisher.LPush(context.TODO(), c.Container.Config.Redis.Buildqueue, projectId).Result()
	if err != nil {
		log.Printf("error publishing to build-queue: %v", err)
		return ctx.JSON(http.StatusInternalServerError, &deployResponsePayload{Error: err.Error()})
	}
	_, err = c.Container.RedisConn.Publisher.HSet(context.TODO(), c.Container.Config.Redis.Statustracker, projectId, "deployed").Result()
	if err != nil {
		log.Printf("error publishing to status tracker: %v", err)
		return ctx.JSON(http.StatusInternalServerError, &deployResponsePayload{Error: err.Error()})
	}

	return c.Redirect(ctx, "project", projectId)
}
