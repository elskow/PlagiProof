package service

import (
	"github.com/elskow/Code-Plagiarism-Detector/dto"
	"github.com/elskow/Code-Plagiarism-Detector/entity"
	"github.com/elskow/Code-Plagiarism-Detector/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HTTPHandler struct {
	FileUseCase *usecase.FileUseCase
}

func (h *HTTPHandler) PingHandler(c *gin.Context) {
	response := dto.ResponseDTO{
		Status:  http.StatusOK,
		Error:   false,
		Message: map[string]string{"description": "pong"},
	}
	c.JSON(http.StatusOK, response)
}

func (h *HTTPHandler) UploadHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response := dto.ResponseDTO{
			Status:  http.StatusBadRequest,
			Error:   true,
			Message: map[string]string{"description": "No file is received"},
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	src, err := file.Open()
	if err != nil {
		response := dto.ResponseDTO{
			Status:  http.StatusInternalServerError,
			Error:   true,
			Message: map[string]string{"description": err.Error()},
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	defer src.Close()

	fileEntity := entity.File{Name: file.Filename, Size: file.Size}
	err = h.FileUseCase.UploadFile(c, fileEntity, src)
	if err != nil {
		response := dto.ResponseDTO{
			Status:  http.StatusInternalServerError,
			Error:   true,
			Message: map[string]string{"description": err.Error()},
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := dto.ResponseDTO{
		Status:  http.StatusOK,
		Error:   false,
		Message: map[string]string{"description": "File uploaded successfully"},
	}
	c.JSON(http.StatusOK, response)
}
