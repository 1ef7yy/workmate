package routes

import (
	"net/http"

	"github.com/1ef7yy/workmate/internal/handlers"
)

func SetupRoutes(taskHandler *handlers.TaskHandler) {
	http.HandleFunc("GET /tasks", taskHandler.ListTasks)
	http.HandleFunc("POST /tasks", taskHandler.CreateTask)
	http.HandleFunc("GET /tasks/{id}", taskHandler.GetTask)
	http.HandleFunc("DELETE /tasks/{id}", taskHandler.DeleteTask)
}
