package storage

import (
	"sync"

	"github.com/1ef7yy/workmate/internal/apperrors"
	"github.com/1ef7yy/workmate/internal/models"
)

type TaskStorage interface {
	Save(task models.Task)
	Get(id int) (models.Task, error)
	Delete(id int) error
	GetAll() []models.Task
	NextID() int
}

type InMemoryTaskQueue struct {
	tasks  map[int]models.Task
	nextID int
	mu     sync.RWMutex
}

func NewTaskStorage() *InMemoryTaskQueue {
	return &InMemoryTaskQueue{
		tasks:  make(map[int]models.Task),
		nextID: 1,
	}
}
func (q *InMemoryTaskQueue) Save(task models.Task) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.tasks[task.ID] = task
}
func (q *InMemoryTaskQueue) Get(id int) (models.Task, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	task, ok := q.tasks[id]

	if !ok {
		return models.Task{}, apperrors.ErrTaskNotFound
	}

	return task, nil
}
func (q *InMemoryTaskQueue) Delete(id int) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	_, ok := q.tasks[id]
	if !ok {
		return apperrors.ErrTaskNotFound
	}

	delete(q.tasks, id)

	return nil
}
func (q *InMemoryTaskQueue) GetAll() []models.Task {
	q.mu.RLock()
	defer q.mu.RUnlock()

	tasks := make([]models.Task, 0, len(q.tasks))

	for _, task := range q.tasks {
		tasks = append(tasks, task)
	}

	return tasks
}
func (q *InMemoryTaskQueue) NextID() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	id := q.nextID
	q.nextID++

	return id
}
