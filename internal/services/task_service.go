package services

import (
	"context"
	"fmt"
	"time"

	"github.com/1ef7yy/workmate/internal/models"
	"github.com/1ef7yy/workmate/internal/storage"
	"github.com/1ef7yy/workmate/pkg/logger"
)

type TaskService interface {
	CreateTask(ctx context.Context) (models.Task, error)
	GetTask(ctx context.Context, id int) (models.Task, error)
	ListTasks(ctx context.Context) ([]models.Task, error)
	DeleteTask(ctx context.Context, id int) error
}

type taskService struct {
	storage storage.TaskStorage
	log     logger.Logger
}

func NewTaskService(storage storage.TaskStorage) TaskService {
	return &taskService{
		storage: storage,
		log:     logger.NewLogger(),
	}
}
func (s *taskService) CreateTask(ctx context.Context) (models.Task, error) {
	task := models.NewTask(s.storage.NextID())

	s.storage.Save(task)

	go s.processTask(task.ID)

	return task, nil
}
func (s *taskService) GetTask(ctx context.Context, id int) (models.Task, error) {
	task, err := s.storage.Get(id)
	if err != nil {
		return models.Task{}, err
	}

	// duration is fixed, if task is not running or pending
	if task.Status == models.StatusFailed || task.Status == models.StatusCompleted {
		return task, nil
	}

	elapsed := time.Now().UTC().Sub(*task.StartedAt)
	task.Duration = elapsed.String()

	return task, nil
}
func (s *taskService) ListTasks(ctx context.Context) ([]models.Task, error) {
	return s.storage.GetAll(), nil
}
func (s *taskService) DeleteTask(ctx context.Context, id int) error {
	return s.storage.Delete(id)
}

func (s *taskService) processTask(id int) {
	task, err := s.storage.Get(id)
	if err != nil {
		s.log.Errorf("error proccessing task with id %d: %s", task.ID, err.Error())
		task.ErrorMsg = err.Error()
		task.Status = models.StatusFailed
		s.storage.Save(task)
		return
	}

	task.Status = models.StatusRunning
	now := time.Now().UTC()
	task.StartedAt = &now
	s.storage.Save(task)

	time.Sleep(time.Minute * 1) // 1 min for development purposes

	task.Status = models.StatusCompleted
	completedAt := time.Now().UTC()
	task.CompletedAt = &completedAt
	duration := completedAt.Sub(*task.StartedAt)
	task.Duration = duration.String()
	task.Result = fmt.Sprintf("task %d completed successfully", id)
	s.storage.Save(task)
}
