package entity

import "github.com/google/uuid"

type File struct {
	ID        uuid.UUID
	Name      string
	Size      int64
	Extension string
	Location  string
}
