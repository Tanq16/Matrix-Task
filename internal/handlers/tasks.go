package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/tanq16/matrix-task/internal/models"
	"github.com/tanq16/matrix-task/internal/storage"
)

// TaskHandler handles all task-related HTTP requests
type TaskHandler struct {
	store    storage.Storage
	tmpl     *template.Template
	basePath string
}

// NewTaskHandler creates a new TaskHandler instance
func NewTaskHandler(store storage.Storage, basePath string) (*TaskHandler, error) {
	tmpl, err := template.ParseGlob(filepath.Join(basePath, "internal/templates/*.html"))
	if err != nil {
		return nil, err
	}

	return &TaskHandler{
		store:    store,
		tmpl:     tmpl,
		basePath: basePath,
	}, nil
}

type response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *TaskHandler) RenderMatrix(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	q1Tasks, _ := h.store.GetTasksByQuadrant(models.QuadrantUrgentImportant)
	q2Tasks, _ := h.store.GetTasksByQuadrant(models.QuadrantNotUrgentImportant)
	q3Tasks, _ := h.store.GetTasksByQuadrant(models.QuadrantUrgentNotImportant)
	q4Tasks, _ := h.store.GetTasksByQuadrant(models.QuadrantNotUrgentNotImportant)

	data := map[string]interface{}{
		"Q1Tasks": q1Tasks,
		"Q2Tasks": q2Tasks,
		"Q3Tasks": q3Tasks,
		"Q4Tasks": q4Tasks,
	}

	h.tmpl.ExecuteTemplate(w, "index.html", data)
}

func (h *TaskHandler) RenderArchive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	archivedTasks, _ := h.store.GetArchivedTasks()
	data := map[string]interface{}{
		"Tasks": archivedTasks,
	}

	h.tmpl.ExecuteTemplate(w, "archive.html", data)
}

func (h *TaskHandler) AddTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Content  string          `json:"content"`
		Quadrant models.Quadrant `json:"quadrant"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, response{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	task := models.NewTask(req.Content, req.Quadrant)
	if err := h.store.AddTask(*task); err != nil {
		sendJSON(w, http.StatusInternalServerError, response{
			Success: false,
			Error:   "Failed to create task",
		})
		return
	}

	sendJSON(w, http.StatusCreated, response{
		Success: true,
		Data:    task,
	})
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		sendJSON(w, http.StatusBadRequest, response{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	if err := h.store.UpdateTask(task); err != nil {
		status := http.StatusInternalServerError
		if _, ok := err.(storage.ErrTaskNotFound); ok {
			status = http.StatusNotFound
		}
		sendJSON(w, status, response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	sendJSON(w, http.StatusOK, response{
		Success: true,
		Data:    task,
	})
}

func (h *TaskHandler) CompleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, response{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	task, err := h.store.GetTask(req.ID)
	if err != nil {
		sendJSON(w, http.StatusNotFound, response{
			Success: false,
			Error:   "Task not found",
		})
		return
	}

	task.Completed = true
	if err := h.store.UpdateTask(task); err != nil {
		sendJSON(w, http.StatusInternalServerError, response{
			Success: false,
			Error:   "Failed to complete task",
		})
		return
	}

	sendJSON(w, http.StatusOK, response{
		Success: true,
		Data:    task,
	})
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		sendJSON(w, http.StatusBadRequest, response{
			Success: false,
			Error:   "Task ID is required",
		})
		return
	}

	if err := h.store.DeleteTask(taskID); err != nil {
		status := http.StatusInternalServerError
		if _, ok := err.(storage.ErrTaskNotFound); ok {
			status = http.StatusNotFound
		}
		sendJSON(w, status, response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	sendJSON(w, http.StatusOK, response{
		Success: true,
	})
}
