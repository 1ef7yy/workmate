package models

import "time"

type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusRunning   TaskStatus = "running"
	StatusCompleted TaskStatus = "completed"
	StatusFailed    TaskStatus = "failed"
)

type Task struct {
	ID          int        `json:"id"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at,omitempty"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Duration    string     `json:"duration,omitempty"`
	Result      string     `json:"result,omitempty"`
	ErrorMsg    string     `json:"error_msg,omitempty"`
}

func NewTask(id int) Task {
	return Task{
		ID:        id,
		Status:    StatusPending,
		CreatedAt: time.Now().UTC(),
	}
}
