package usecase

import (
	"context"
	"github.com/elskow/Code-Plagiarism-Detector/domain/entity"
	"github.com/elskow/Code-Plagiarism-Detector/domain/repository"
	"github.com/google/uuid"
)

type CheckPlagiarismUseCase struct {
	PlagiarismCheckRepository repository.PlagiarismCheckRepository
}

func (u *CheckPlagiarismUseCase) Execute(ctx context.Context, fileID uuid.UUID) (uuid.UUID, error) {
	check := entity.PlagiarismCheck{
		ID:     uuid.New(),
		FileID: fileID,
		Status: "queued",
	}
	err := u.PlagiarismCheckRepository.QueueCheck(ctx, check)
	if err != nil {
		return uuid.Nil, err
	}
	return check.ID, nil
}
