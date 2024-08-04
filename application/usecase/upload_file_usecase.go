package usecase

import (
	"context"
	"github.com/elskow/Code-Plagiarism-Detector/domain/entity"
	"github.com/elskow/Code-Plagiarism-Detector/domain/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"io"
)

type UploadFileUseCase struct {
	FileRepository repository.FileRepository
	QueueService   repository.QueueService
}

func (u *UploadFileUseCase) Execute(ctx context.Context, file entity.File, src io.Reader) (uuid.UUID, error) {
	file.ID = uuid.New()
	err := u.FileRepository.SaveFile(ctx, file, src)
	if err != nil {
		return uuid.Nil, err
	}
	err = u.QueueService.InsertFile(ctx, file)
	if err != nil {
		return uuid.Nil, err
	}

	logrus.Infof("File %s has been uploaded", file.ID)
	logrus.Infof("File withname %s has been uploaded", file.Name)
	logrus.Infof("File withsize %d has been uploaded", file.Size)
	logrus.Infof("File withextension %s has been uploaded", file.Extension)
	return file.ID, nil
}

func (u *UploadFileUseCase) GetStatus(ctx context.Context, checkID string) (string, error) {
	return u.QueueService.GetStatus(ctx, checkID)
}
