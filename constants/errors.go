package constants

import "errors"

var (
	ErrFileNotFound  = errors.New("file not found")
	ErrFileExtension = errors.New("file extension not allowed")
	ErrFileSave      = errors.New("failed to save file")
)
