package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	Cfg "github.com/3ssalunke/vercelc/shared/config"
	"github.com/3ssalunke/vercelc/shared/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	Client *s3.Client
	Bucket string
}

func NewS3Storage(cfg *Cfg.Config) (*S3Storage, error) {
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}

	awsConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithHTTPClient(httpClient), config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.S3.Accesskey, cfg.S3.Secretkey, "")), config.WithRegion(cfg.S3.Region))
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

func (storage *S3Storage) CopyFolder(folderPath string) error {
	localFolderPath, err := util.GetPathForFolder(folderPath)
	if err != nil {
		return err
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	currentDir = strings.Replace(currentDir, "\\", "/", -1)

	// Walk through local folder and upload files to S3
	err = filepath.Walk(localFolderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		updatedPath := strings.Replace(path, "\\", "/", -1)
		filekey := strings.TrimPrefix(updatedPath, fmt.Sprintf("%s/", currentDir))

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Skip hidden directories/files
		if strings.ContainsRune(filepath.Dir(path), '.') {
			return nil
		}

		// // Open the file
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		var partMiBs int64 = 10
		// Create a new uploader with custom options
		uploader := manager.NewUploader(storage.Client, func(u *manager.Uploader) {
			u.PartSize = partMiBs * 1024 * 1024
		})

		// Upload the file with multipart upload
		_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
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

func (storage *S3Storage) CopyBuildFolder(folderPath string) error {
	localFolderPath, err := util.GetPathForFolder(folderPath)
	if err != nil {
		return err
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	currentDir = strings.Replace(currentDir, "\\", "/", -1)

	// Walk through local folder and upload files to S3
	err = filepath.Walk(localFolderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		updatedPath := strings.Replace(path, "\\", "/", -1)
		filekey := strings.TrimPrefix(updatedPath, fmt.Sprintf("%s/", currentDir+"/build"))
		filekey = strings.Replace(filekey, "build", "dist", -1)
		log.Println(filekey)

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Skip hidden directories/files
		if strings.ContainsRune(filepath.Dir(path), '.') {
			return nil
		}

		// // Open the file
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		var partMiBs int64 = 10
		// Create a new uploader with custom options
		uploader := manager.NewUploader(storage.Client, func(u *manager.Uploader) {
			u.PartSize = partMiBs * 1024 * 1024
		})

		// Upload the file with multipart upload
		_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
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

func (storage *S3Storage) DownloadFolder(folderPath string) error {
	pathToDownloadedFolder, err := util.GetPathForFolder("build")
	if err != nil {
		return err
	}

	listObjectInput := &s3.ListObjectsV2Input{
		Bucket: aws.String(storage.Bucket),
		Prefix: aws.String(folderPath),
	}

	paginator := s3.NewListObjectsV2Paginator(storage.Client, listObjectInput)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return err
		}

		for _, obj := range page.Contents {
			downloadInput := &s3.GetObjectInput{
				Bucket: aws.String(storage.Bucket),
				Key:    obj.Key,
			}

			resp, err := storage.Client.GetObject(context.TODO(), downloadInput)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			filePath := filepath.Join(pathToDownloadedFolder, *obj.Key)
			dirPath := filepath.Dir(filePath)
			if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
				return err
			}

			file, err := os.Create(filePath)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = file.ReadFrom(resp.Body)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
