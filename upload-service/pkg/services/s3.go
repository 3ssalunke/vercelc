package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	Cfg "github.com/3ssalunke/vercelc/shared/config"
	"github.com/3ssalunke/vercelc/shared/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	Client *s3.Client
	Bucket string
}

func NewS3Storage(cfg *Cfg.Config) (*S3Storage, error) {
	awsConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.S3.Accesskey, cfg.S3.Secretkey, "")), config.WithRegion(cfg.S3.Region))
	if err != nil {
		log.Printf("unable to load AWS SDK config, %v", err)
		return nil, fmt.Errorf("unable to load AWS SDK config, %v", err)
	}

	client := s3.NewFromConfig(awsConfig)

	return &S3Storage{
		Client: client,
		Bucket: cfg.S3.Bucket,
	}, nil
}

func (storage *S3Storage) CopyFolder(projectId string) error {
	localFolderPath, err := util.GetLocalCloneFolder(projectId)
	if err != nil {
		return err
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	// Walk through local folder and upload files to S3
	err = filepath.Walk(localFolderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		filekey := strings.TrimPrefix(path, fmt.Sprintf("%s/", currentDir))
		// Skip directories
		if info.IsDir() {
			return nil
		}

		if strings.ContainsRune(filepath.Dir(path), '.') {
			return nil
		}

		// Open the file
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Calculate the size of the file
		fileInfo, err := file.Stat()
		if err != nil {
			return err
		}
		fileSize := fileInfo.Size()
		log.Println(fileSize)

		if fileSize > 2000 {
			return nil
		}

		// Upload the file to S3
		_, err = storage.Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(storage.Bucket),
			Key:    aws.String(filekey),
			Body:   file,
		})
		if err != nil {
			return err
		}

		fmt.Println("Uploaded", path, "to", storage.Bucket+"/"+filekey)

		return nil
	})

	if err != nil {
		return err
	}

	fmt.Println("All files uploaded successfully")
	return nil
}
