package usecase

import (
	"context"
	"github.com/elskow/Code-Plagiarism-Detector/entity"
	"github.com/elskow/Code-Plagiarism-Detector/repository"
	"io"
)

type FileUseCase struct {
	FileRepository repository.FileRepository
}

func (uc *FileUseCase) UploadFile(ctx context.Context, file entity.File, src io.Reader) error {
	return uc.FileRepository.SaveFile(ctx, file, src)
}
