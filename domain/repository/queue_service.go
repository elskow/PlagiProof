package repository

import (
	"context"
	"github.com/elskow/Code-Plagiarism-Detector/domain/entity"
)

type QueueService interface {
	InsertFile(ctx context.Context, file entity.File) error
	GetStatus(ctx context.Context, checkID string) (string, error)
}
