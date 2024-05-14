package listener

import (
	"context"
	"fmt"
	"log"

	"github.com/3ssalunke/vercelc/deploy-service/pkg/builder"
	"github.com/3ssalunke/vercelc/deploy-service/pkg/services"
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
	log.Println("deploy service started...")
	for {
		result, err := l.container.RedisConn.Subscriber.BRPop(context.TODO(), 0, l.container.Config.Redis.Buildqueue).Result()
		if err != nil {
			return nil
		}

		projectId := result[1]
		if err := l.container.S3Storage.DownloadFolder(fmt.Sprintf("output/%s", projectId)); err != nil {
			log.Printf("failed to download project folder from s3 for project %s: %v", projectId, err)
			continue
		}

		if err := builder.BuildProject(projectId); err != nil {
			log.Printf("failed to build project for project %s: %v", projectId, err)
			continue
		}

		if err := l.container.S3Storage.CopyBuildFolder(fmt.Sprintf("build/output/%s/build", projectId)); err != nil {
			log.Printf("failed to copy built folder to s3 for project %s: %v", projectId, err)
			continue
		}
	}
}
