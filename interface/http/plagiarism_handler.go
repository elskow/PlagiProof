package http

import (
	"github.com/elskow/Code-Plagiarism-Detector/application/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PlagiarismHandler struct {
	CheckPlagiarismUseCase *usecase.UploadFileUseCase
}

func (h *PlagiarismHandler) GetStatusHandler(c *gin.Context) {
	checkID := c.Param("check_id")
	status, err := h.CheckPlagiarismUseCase.GetStatus(c.Request.Context(), checkID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": status})
}
