package listener

import (
	"context"
	"fmt"

	"github.com/3ssalunke/vercelc/deploy-service/pkg/services"
	"github.com/labstack/gommon/log"
)

type Listener struct {
	container *services.Container
}

func NewListener(c *services.Container) *Listener {
	return &Listener{
		container: c,
	}
}

func (l *Listener) Start() error {
	for {
		result, err := l.container.RedisConn.Subscriber.BRPop(context.TODO(), 0, l.container.Config.Redis.Buildqueue).Result()
		if err != nil {
			return nil
		}
		projectId := result[1]
		if err := l.container.S3Storage.DownloadFolder(fmt.Sprintf("output/%s", projectId)); err != nil {
			log.Printf("failed to download project folder from s3: %v", err)
			continue
		}
		// return nil
	}
}
