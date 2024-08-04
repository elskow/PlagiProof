package repository

import (
	"context"
	"github.com/elskow/Code-Plagiarism-Detector/domain/entity"
	"io"
)

type FileRepository interface {
	SaveFile(ctx context.Context, file entity.File, src io.Reader) error
}
