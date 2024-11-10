package models

import "github.com/google/uuid"

// Quadrant represents the four quadrants of the Eisenhower Matrix
type Quadrant int

const (
	QuadrantUrgentImportant       Quadrant = iota + 1 // Q1
	QuadrantNotUrgentImportant                        // Q2
	QuadrantUrgentNotImportant                        // Q3
	QuadrantNotUrgentNotImportant                     // Q4
)

// Task represents a single task in the system
type Task struct {
	ID        string   `json:"id"`
	Content   string   `json:"content"`
	Quadrant  Quadrant `json:"quadrant"`
	Completed bool     `json:"completed"`
}

// NewTask creates a new task with the given content and quadrant
func NewTask(content string, quadrant Quadrant) *Task {
	return &Task{
		ID:        uuid.New().String(),
		Content:   content,
		Quadrant:  quadrant,
		Completed: false,
	}
}

// QuadrantName returns a human-readable name for each quadrant
func (q Quadrant) String() string {
	switch q {
	case QuadrantUrgentImportant:
		return "Urgent & Important"
	case QuadrantNotUrgentImportant:
		return "Not Urgent & Important"
	case QuadrantUrgentNotImportant:
		return "Urgent & Not Important"
	case QuadrantNotUrgentNotImportant:
		return "Not Urgent & Not Important"
	default:
		return "Unknown"
	}
}
