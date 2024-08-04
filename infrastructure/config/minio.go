package config

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"os"
)

func InitMinio() (*minio.Client, error) {
	err := godotenv.Load()
	if err != nil {
		logrus.Errorf("Failed to load .env file: %v", err)
		return nil, err
	}

	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		logrus.Errorf("Failed to create MinIO client: %v", err)
		return nil, err
	}

	bucketName := "uploads"
	location := "us-east-1"

	exists, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		logrus.Errorf("Failed to check if bucket exists: %v", err)
		return nil, err
	}

	if !exists {
		err = client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: location})
		if err != nil {
			logrus.Errorf("Failed to create bucket: %v", err)
			return nil, err
		}
	}

	return client, nil
}
