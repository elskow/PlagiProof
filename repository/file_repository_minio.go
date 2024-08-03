package repository

import (
	"context"
	"github.com/elskow/Code-Plagiarism-Detector/entity"
	"github.com/minio/minio-go/v7"
	"io"
)

type MinioFileRepository struct {
	Client     *minio.Client
	BucketName string
}

func (r *MinioFileRepository) SaveFile(ctx context.Context, file entity.File, src io.Reader) error {
	_, err := r.Client.PutObject(ctx, r.BucketName, file.Name, src, file.Size, minio.PutObjectOptions{})
	return err
}
