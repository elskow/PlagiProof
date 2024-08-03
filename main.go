package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/elskow/Code-Plagiarism-Detector/config"
	"github.com/elskow/Code-Plagiarism-Detector/repository"
	"github.com/elskow/Code-Plagiarism-Detector/service"
	"github.com/elskow/Code-Plagiarism-Detector/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	minioClient, err := config.InitMinio()
	if err != nil {
		fmt.Printf("Failed to initialize MinIO: %v\n", err)
		return
	}

	fileRepository := &repository.MinioFileRepository{
		Client:     minioClient,
		BucketName: "uploads",
	}
	fileUseCase := &usecase.FileUseCase{FileRepository: fileRepository}
	httpHandler := &service.HTTPHandler{FileUseCase: fileUseCase}

	router := gin.New()
	router.Use(gin.Logger())
	router.GET("/ping", httpHandler.PingHandler)
	router.POST("/upload", httpHandler.UploadHandler)

	server := &http.Server{
		Addr:    "localhost:8000",
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("listen: %s\n", err)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server Shutdown: %v\n", err)
	}
}
