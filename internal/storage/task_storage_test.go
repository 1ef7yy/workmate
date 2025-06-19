package storage

import (
	"sync"
	"testing"

	"github.com/1ef7yy/workmate/internal/apperrors"
	"github.com/1ef7yy/workmate/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryTaskQueue(t *testing.T) {
	storage := NewTaskStorage()

	t.Run("storage and get", func(t *testing.T) {
		task := models.NewTask(1)
		storage.Save(task)

		storagedTask, err := storage.Get(1)
		assert.NoError(t, err)
		assert.Equal(t, task, storagedTask)
	})

	t.Run("get non existent", func(t *testing.T) {
		_, err := storage.Get(999)
		assert.Error(t, err)
	})

	taskNum := 1000
	var wg sync.WaitGroup

	t.Run("concurrent storage access", func(t *testing.T) {
		wg.Add(taskNum)
		for i := 0; i < taskNum; i++ {
			go func(id int) {
				defer wg.Done()
				task := models.NewTask(id)
				storage.Save(task)
			}(i)
		}
		wg.Wait()

		for i := 0; i < taskNum; i++ {
			_, err := storage.Get(i)

			if err != nil {
				if err == apperrors.ErrTaskNotFound {
					t.Errorf("task %d not found in queue", i)
				} else {
					t.Errorf("error getting task from queue concurrently: %s", err.Error())
				}
			}
		}
	})

	t.Run("concurrent RW", func(t *testing.T) {
		wg.Add(taskNum * 2)
		for i := 0; i < taskNum; i++ {
			go func(id int) {
				defer wg.Done()
				task := models.Task{ID: id, Status: models.StatusRunning}
				storage.Save(task)
			}(i)

			go func(id int) {
				defer wg.Done()
				if _, err := storage.Get(id); err != nil {
					t.Errorf("failed to read task %d: %v", id, err)
				}
			}(i)
		}
		wg.Wait()
	})

	t.Run("concurrent delete", func(t *testing.T) {
		for i := 0; i < taskNum; i++ {
			storage.Save(models.Task{ID: i})
		}

		wg.Add(taskNum)
		for i := 0; i < taskNum; i++ {
			go func(id int) {
				defer wg.Done()
				if err := storage.Delete(id); err != nil {
					t.Errorf("failed to delete task %d: %v", id, err)
				}
			}(i)
		}
		wg.Wait()

		for i := 0; i < taskNum; i++ {
			if _, err := storage.Get(i); err == nil {
				t.Errorf("task %d still exists after concurrent delete", i)
			}
		}
	})
}
