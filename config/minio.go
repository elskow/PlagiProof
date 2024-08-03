package config

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"os"
)

func InitMinio() (*minio.Client, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")

	client, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	bucketName := "uploads"
	location := "us-east-1"

	exists, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: location})
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}
