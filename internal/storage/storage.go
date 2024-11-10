package storage

import "github.com/tanq16/matrix-task/internal/models"

// Storage defines the interface for task storage operations
type Storage interface {
	// AddTask adds a new task to storage
	AddTask(task models.Task) error

	// GetTask retrieves a task by ID
	GetTask(id string) (models.Task, error)

	// UpdateTask updates an existing task
	UpdateTask(task models.Task) error

	// DeleteTask removes a task from storage
	DeleteTask(id string) error

	// GetTasksByQuadrant retrieves all active (non-completed) tasks for a specific quadrant
	GetTasksByQuadrant(quadrant models.Quadrant) ([]models.Task, error)

	// GetArchivedTasks retrieves all completed tasks
	GetArchivedTasks() ([]models.Task, error)
}

// ErrTaskNotFound is returned when a task cannot be found in storage
type ErrTaskNotFound struct {
	TaskID string
}

func (e ErrTaskNotFound) Error() string {
	return "task not found: " + e.TaskID
}
