package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *mux.Router, handlers *Handlers) {
	// Add CORS and logging middleware
	router.Use(corsMiddleware)
	router.Use(loggingMiddleware)

	// Health check endpoint
	router.HandleFunc("/health", handlers.Book.HealthCheck).Methods("GET")

	// API routes - ensure these are registered first
	api := router.PathPrefix("/api/v1").Subrouter()
	api.Use(jsonMiddleware)

	// Book API routes
	books := api.PathPrefix("/books").Subrouter()
	books.HandleFunc("", handlers.Book.CreateBook).Methods("POST")
	books.HandleFunc("", handlers.Book.GetBooks).Methods("GET")
	books.HandleFunc("/{id:[0-9]+}", handlers.Book.GetBook).Methods("GET")
	books.HandleFunc("/{id:[0-9]+}", handlers.Book.UpdateBook).Methods("PUT")
	books.HandleFunc("/{id:[0-9]+}", handlers.Book.DeleteBook).Methods("DELETE")
	books.HandleFunc("/isbn/{isbn}", handlers.Book.GetBookByISBN).Methods("GET")

	// Web UI routes - these should come last to not interfere with API
	router.HandleFunc("/", serveWebUI).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/"))))
	
	// Catch-all for SPA routing - this ensures the web app works for all routes
	router.PathPrefix("/").HandlerFunc(serveWebUI).Methods("GET")
}

// serveWebUI serves the web interface
func serveWebUI(w http.ResponseWriter, r *http.Request) {
	// Set proper content type for HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "./web/templates/index.html")
}