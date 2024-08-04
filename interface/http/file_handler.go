package http

import (
	"github.com/elskow/Code-Plagiarism-Detector/application/usecase"
	"github.com/elskow/Code-Plagiarism-Detector/constants"
	"github.com/elskow/Code-Plagiarism-Detector/domain/entity"
	"github.com/elskow/Code-Plagiarism-Detector/dto"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

type FileHandler struct {
	UploadFileUseCase *usecase.UploadFileUseCase
}

func (h *FileHandler) UploadHandler(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseDTO{
			Error:   true,
			Message: constants.ErrFileNotFound.Error(),
		})
		return
	}
	defer file.Close()

	if !isAllowedExtension(header.Filename, constants.SupportedExtensions) {
		c.JSON(http.StatusBadRequest, dto.ResponseDTO{
			Error:   true,
			Message: constants.ErrFileExtension.Error(),
		})
		return
	}

	fileEntity := entity.File{
		Name:      header.Filename,
		Size:      header.Size,
		Extension: filepath.Ext(header.Filename),
	}

	fileID, err := h.UploadFileUseCase.Execute(c.Request.Context(), fileEntity, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseDTO{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseDTO{
		Error:   false,
		Message: "File uploaded successfully with ID " + fileID.String(),
	})
}

func isAllowedExtension(filename string, allowedExtensions []string) bool {
	ext := filepath.Ext(filename)
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			return true
		}
	}
	return false
}
