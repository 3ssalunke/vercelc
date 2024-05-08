package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/3ssalunke/vercelc/upload-service/pkg/controller"
	"github.com/labstack/echo/v4"
)

type (
	status struct {
		controller.Controller
	}
	statusResponsePayload struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}
)

func (s *status) Get(ctx echo.Context) error {
	projectId := ctx.QueryParam("id")
	if projectId == "" {
		err := fmt.Errorf("id query param is missing from request")
		log.Print(err.Error())
		return ctx.JSON(http.StatusBadRequest, &statusResponsePayload{Error: err.Error()})
	}
	status, err := s.Container.RedisConn.Subscriber.HGet(context.TODO(), s.Container.Config.Redis.Statustracker, projectId).Result()
	if err != nil {
		log.Printf("error retrieving status for id %s: %v", projectId, err)
		return ctx.JSON(http.StatusInternalServerError, &statusResponsePayload{Error: err.Error()})
	}
	return ctx.JSON(200, &statusResponsePayload{Status: status})
}
