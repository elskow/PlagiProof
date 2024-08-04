package service

import (
	"context"
	"github.com/elskow/Code-Plagiarism-Detector/domain/entity"
	"github.com/minio/minio-go/v7"
	"io"
)

type MinioFileService struct {
	Client     *minio.Client
	BucketName string
}

func (s *MinioFileService) SaveFile(ctx context.Context, file entity.File, src io.Reader) error {
	filename := file.ID.String() + file.Extension
	_, err := s.Client.PutObject(ctx, s.BucketName, filename, src, file.Size, minio.PutObjectOptions{})

	return err
}
