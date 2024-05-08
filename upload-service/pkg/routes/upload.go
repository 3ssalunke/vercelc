package routes

import (
	"context"
	"log"
	"net/http"

	"github.com/3ssalunke/vercelc/shared/util"
	"github.com/3ssalunke/vercelc/upload-service/pkg/controller"
	"github.com/labstack/echo/v4"
)

type (
	upload struct {
		controller.Controller
	}

	uploadRequestPayload struct {
		RepoUrl string `json:"repo_url" validate:"required"`
	}

	uploadResponsePayload struct {
		Id    string `json:"msg"`
		Error string `json:"error"`
	}
)

func (c *upload) Post(ctx echo.Context) error {
	payload := new(uploadRequestPayload)
	if err := ctx.Bind(payload); err != nil {
		log.Printf("error parsing request payload: %v", err)
		return ctx.JSON(http.StatusBadRequest, &uploadResponsePayload{Error: err.Error()})
	}
	if err := c.Container.Validator.Validate(payload); err != nil {
		log.Printf("error validating request payload: %v", err)
		return ctx.JSON(http.StatusBadRequest, &uploadResponsePayload{Error: err.Error()})
	}
	projectId := util.GenerateRandomId()
	destinationFolder, err := util.GetLocalCloneFolder(projectId)
	if err != nil {
		log.Printf("error getting destination folder for cloning git repo: %v", err)
		return ctx.JSON(http.StatusInternalServerError, &uploadResponsePayload{Error: err.Error()})
	}
	if err := util.CloneRepo(payload.RepoUrl, destinationFolder); err != nil {
		log.Printf("error cloning git repo: %v", err)
		return ctx.JSON(http.StatusInternalServerError, &uploadResponsePayload{Error: err.Error()})
	}
	if err := c.Container.S3Storage.CopyFolder(projectId); err != nil {
		log.Printf("error copying repo to s3: %v", err)
		return ctx.JSON(http.StatusInternalServerError, &uploadResponsePayload{Error: err.Error()})
	}
	_, err = c.Container.RedisConn.Publisher.LPush(context.TODO(), c.Container.Config.Redis.Buildqueue, projectId).Result()
	if err != nil {
		log.Printf("error publishing to build-queue: %v", err)
		return ctx.JSON(http.StatusInternalServerError, &uploadResponsePayload{Error: err.Error()})
	}
	_, err = c.Container.RedisConn.Publisher.HSet(context.TODO(), c.Container.Config.Redis.Statustracker, projectId, "uploaded").Result()
	if err != nil {
		log.Printf("error publishing to status tracker: %v", err)
		return ctx.JSON(http.StatusInternalServerError, &uploadResponsePayload{Error: err.Error()})
	}
	return ctx.JSON(200, &uploadResponsePayload{Id: projectId})
}
