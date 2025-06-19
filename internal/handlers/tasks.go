package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/1ef7yy/workmate/internal/apperrors"
	"github.com/1ef7yy/workmate/internal/services"
	"github.com/1ef7yy/workmate/pkg/logger"
)

type TaskHandler struct {
	service services.TaskService
	log     logger.Logger
}

func NewTaskHandler(service services.TaskService) *TaskHandler {
	return &TaskHandler{
		service: service,
		log:     logger.NewLogger(),
	}
}
func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.ListTasks(r.Context())
	if err != nil {
		h.log.Errorf("error listing tasks: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(tasks) == 0 {
		http.Error(w, "tasks not found", http.StatusNotFound)
		return
	}

	h.respondJSON(w, http.StatusOK, tasks)
}
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	task, err := h.service.CreateTask(r.Context())

	if err != nil {
		h.log.Errorf("error creating task: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusCreated, task)
}
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		http.Error(w, "invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := h.service.GetTask(r.Context(), id)

	if err == apperrors.ErrTaskNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.respondJSON(w, http.StatusOK, task)

}
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid task ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteTask(r.Context(), id)

	if err == apperrors.ErrTaskNotFound {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	if err != nil {
		h.log.Errorf("error deleting task with id %d: %s", id, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) respondJSON(w http.ResponseWriter, statusCode int, payload any) {
	resp, err := json.Marshal(payload)
	if err != nil {
		h.log.Errorf("error marshalling payload: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)

	_, err = w.Write(resp)
	if err != nil {
		h.log.Errorf("error writing to client: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
