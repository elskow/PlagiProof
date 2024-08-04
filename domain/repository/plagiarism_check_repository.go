package repository

import (
	"context"
	"github.com/elskow/Code-Plagiarism-Detector/domain/entity"
	"github.com/google/uuid"
)

type PlagiarismCheckRepository interface {
	QueueCheck(ctx context.Context, check entity.PlagiarismCheck) error
	GetCheckStatus(ctx context.Context, checkID uuid.UUID) (entity.PlagiarismCheck, error)
}
