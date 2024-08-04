package service

import (
	"context"
	"database/sql"
	"github.com/elskow/Code-Plagiarism-Detector/domain/entity"
)

type PostgresQueueService struct {
	DB *sql.DB
}

func (p *PostgresQueueService) InsertFile(ctx context.Context, file entity.File) error {
	_, err := p.DB.ExecContext(ctx, "INSERT INTO queue (file_id, status) VALUES ($1, $2)", file.ID, "queued")
	return err
}

func (p *PostgresQueueService) GetStatus(ctx context.Context, checkID string) (string, error) {
	var status string
	err := p.DB.QueryRowContext(ctx, "SELECT status FROM queue WHERE file_id = $1", checkID).Scan(&status)
	if err != nil {
		return "", err
	}
	return status, nil
}
