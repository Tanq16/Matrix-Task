package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/tanq16/matrix-task/internal/handlers"
	"github.com/tanq16/matrix-task/internal/storage"
)

func main() {
	// Command line flags
	port := flag.Int("port", 8080, "Port to serve the application")
	flag.Parse()

	// Initialize storage
	store := storage.NewMemoryStorage()

	// Initialize handlers with embedded templates
	taskHandler, err := handlers.NewTaskHandler(store)
	if err != nil {
		log.Fatal("Failed to initialize task handler:", err)
	}

	// Get the embedded static file system
	staticFS, err := handlers.GetStaticFileSystem()
	if err != nil {
		log.Fatal("Failed to get static file system:", err)
	}

	// Create router
	mux := http.NewServeMux()

	// Static file server using embedded files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticFS)))

	// Register routes
	// Page routes
	mux.HandleFunc("/", taskHandler.RenderMatrix)
	mux.HandleFunc("/archive", taskHandler.RenderArchive)

	// API routes
	mux.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			taskHandler.AddTask(w, r)
		case http.MethodDelete:
			taskHandler.DeleteTask(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/tasks/complete", taskHandler.CompleteTask)

	// Create custom server with timeouts
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", *port),
		Handler:        logMiddleware(corsMiddleware(mux)),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server
	log.Printf("Starting server on http://localhost:%d", *port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

// Middleware for logging requests
func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom response writer to capture the status code
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		// Log request details
		log.Printf(
			"%s %s %d %v",
			r.Method,
			r.URL.Path,
			rw.statusCode,
			duration,
		)
	})
}

// Middleware for CORS (useful for local development)
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers for all responses
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Custom response writer to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
