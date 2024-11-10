package storage

import (
	"sync"

	"github.com/tanq16/matrix-task/internal/models"
)

// MemoryStorage implements Storage interface using in-memory maps
type MemoryStorage struct {
	tasks map[string]*models.Task
	mu    sync.RWMutex
}

// NewMemoryStorage creates a new instance of MemoryStorage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		tasks: make(map[string]*models.Task),
	}
}

func (m *MemoryStorage) AddTask(task *models.Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.tasks[task.ID] = task
	return nil
}

func (m *MemoryStorage) GetTask(id string) (*models.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	task, exists := m.tasks[id]
	if !exists {
		return nil, ErrTaskNotFound{TaskID: id}
	}
	return task, nil
}

func (m *MemoryStorage) UpdateTask(task *models.Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.tasks[task.ID]; !exists {
		return ErrTaskNotFound{TaskID: task.ID}
	}

	m.tasks[task.ID] = task
	return nil
}

func (m *MemoryStorage) DeleteTask(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.tasks[id]; !exists {
		return ErrTaskNotFound{TaskID: id}
	}

	delete(m.tasks, id)
	return nil
}

func (m *MemoryStorage) GetTasksByQuadrant(quadrant models.Quadrant) ([]*models.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var tasks []*models.Task
	for _, task := range m.tasks {
		if task.Quadrant == quadrant && !task.Completed {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

func (m *MemoryStorage) GetArchivedTasks() ([]*models.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var tasks []*models.Task
	for _, task := range m.tasks {
		if task.Completed {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}
