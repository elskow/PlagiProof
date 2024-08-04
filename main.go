package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/elskow/Code-Plagiarism-Detector/application/usecase"
	"github.com/elskow/Code-Plagiarism-Detector/infrastructure/config"
	"github.com/elskow/Code-Plagiarism-Detector/infrastructure/service"
	http2 "github.com/elskow/Code-Plagiarism-Detector/interface/http"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type App struct {
	Server *http.Server
	Log    *logrus.Logger
}

func NewApp() (*App, error) {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	minioClient, err := config.InitMinio()
	if err != nil {
		log.Fatalf("Failed to initialize MinIO: %v", err)
	}

	db, err := config.InitPostgres()
	if err != nil {
		log.Fatalf("Failed to initialize Postgres: %v", err)
	}

	fileService := &service.MinioFileService{
		Client:     minioClient,
		BucketName: "uploads",
	}
	queueService := &service.PostgresQueueService{DB: db}

	fileUseCase := &usecase.UploadFileUseCase{
		FileRepository: fileService,
		QueueService:   queueService,
	}

	fileHandler := &http2.FileHandler{UploadFileUseCase: fileUseCase}
	plagiarismHandler := &http2.PlagiarismHandler{CheckPlagiarismUseCase: fileUseCase}

	router := gin.New()
	router.Use(gin.Logger())
	router.POST("/plagiarism", fileHandler.UploadHandler)
	router.GET("/plagiarism/:check_id", plagiarismHandler.GetStatusHandler)

	router.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "")
	})

	server := &http.Server{
		Addr:    "localhost:8000",
		Handler: router,
	}

	return &App{
		Server: server,
		Log:    log,
	}, nil
}

func (app *App) Run() {
	go func() {
		if err := app.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.Log.Fatalf("listen: %s\n", err)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	if err := app.Server.Shutdown(ctx); err != nil {
		app.Log.Fatalf("Server Shutdown: %v", err)
	}
	app.Log.Println("Server exiting")
}

func main() {
	app, err := NewApp()
	if err != nil {
		fmt.Printf("Failed to initialize application: %v\n", err)
		return
	}
	app.Run()
}
