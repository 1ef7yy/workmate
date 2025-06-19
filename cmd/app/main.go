package main

import (
	"net/http"

	"github.com/1ef7yy/workmate/internal/handlers"
	"github.com/1ef7yy/workmate/internal/routes"
	"github.com/1ef7yy/workmate/internal/services"
	"github.com/1ef7yy/workmate/internal/storage"
	"github.com/1ef7yy/workmate/pkg/logger"
)

const (
	HTTP_PORT = ":8080"
)

func main() {
	taskStorage := storage.NewTaskStorage()
	taskService := services.NewTaskService(taskStorage)
	taskHandler := handlers.NewTaskHandler(taskService)

	routes.SetupRoutes(taskHandler)

	logger := logger.NewLogger()

	logger.Infof("starting service on port %s", HTTP_PORT)
	if err := http.ListenAndServe(HTTP_PORT, nil); err != nil {
		logger.Fatalf("error running service: %s", err.Error())
	}
}
