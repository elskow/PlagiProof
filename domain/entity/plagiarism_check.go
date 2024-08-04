package entity

import "github.com/google/uuid"

type PlagiarismCheck struct {
	ID     uuid.UUID
	FileID uuid.UUID
	Status string
	Result string
}
