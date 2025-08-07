package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"library-management/internal/domain"
	"library-management/internal/service"
	"library-management/pkg/logger"
)

type BookHandler struct {
	service service.BookService
	logger  logger.Logger
}

type Handlers struct {
	Book *BookHandler
}

// NewHandlers creates a new handlers instance
func NewHandlers(bookService service.BookService, log logger.Logger) *Handlers {
	return &Handlers{
		Book: &BookHandler{
			service: bookService,
			logger:  log,
		},
	}
}

// Response represents a standard API response
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// CreateBook handles POST /api/v1/books
func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateBookRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	book, err := h.service.CreateBook(r.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create book", "error", err)
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusCreated, "Book created successfully", book)
}

// GetBook handles GET /api/v1/books/{id}
func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	book, err := h.service.GetBookByID(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get book", "error", err, "id", id)
		h.respondError(w, http.StatusNotFound, "Book not found")
		return
	}

	h.respondSuccess(w, http.StatusOK, "Book retrieved successfully", book)
}

// GetBooks handles GET /api/v1/books
func (h *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for filtering
	filter := &domain.BookFilter{
		Author: r.URL.Query().Get("author"),
		Genre:  r.URL.Query().Get("genre"),
		Search: r.URL.Query().Get("search"),
	}

	// Parse available filter
	if availableStr := r.URL.Query().Get("available"); availableStr != "" {
		if available, err := strconv.ParseBool(availableStr); err == nil {
			filter.Available = &available
		}
	}

	books, err := h.service.GetAllBooks(r.Context(), filter)
	if err != nil {
		h.logger.Error("Failed to get books", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to retrieve books")
		return
	}

	// Get count for metadata
	count, err := h.service.GetBooksCount(r.Context(), filter)
	if err != nil {
		h.logger.Warn("Failed to get books count", "error", err)
		count = len(books) // Fallback to actual count
	}

	response := map[string]interface{}{
		"books": books,
		"meta": map[string]interface{}{
			"total": count,
			"count": len(books),
		},
	}

	h.respondSuccess(w, http.StatusOK, "Books retrieved successfully", response)
}

// UpdateBook handles PUT /api/v1/books/{id}
func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	var req domain.UpdateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	book, err := h.service.UpdateBook(r.Context(), id, &req)
	if err != nil {
		h.logger.Error("Failed to update book", "error", err, "id", id)
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusOK, "Book updated successfully", book)
}

// DeleteBook handles DELETE /api/v1/books/{id}
func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	err = h.service.DeleteBook(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to delete book", "error", err, "id", id)
		h.respondError(w, http.StatusNotFound, "Book not found")
		return
	}

	h.respondSuccess(w, http.StatusOK, "Book deleted successfully", nil)
}

// GetBookByISBN handles GET /api/v1/books/isbn/{isbn}
func (h *BookHandler) GetBookByISBN(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	isbn := vars["isbn"]

	book, err := h.service.GetBookByISBN(r.Context(), isbn)
	if err != nil {
		h.logger.Error("Failed to get book by ISBN", "error", err, "isbn", isbn)
		h.respondError(w, http.StatusNotFound, "Book not found")
		return
	}

	h.respondSuccess(w, http.StatusOK, "Book retrieved successfully", book)
}

// HealthCheck handles GET /health
func (h *BookHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.respondSuccess(w, http.StatusOK, "Service is healthy", map[string]string{
		"status": "ok",
		"service": "library-management-api",
	})
}

// respondSuccess sends a success response
func (h *BookHandler) respondSuccess(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	// Ensure JSON content type is set
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	
	response := Response{
		Status:  "success",
		Message: message,
		Data:    data,
	}
	
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("Failed to encode JSON response", "error", err)
	}
}

// respondError sends an error response
func (h *BookHandler) respondError(w http.ResponseWriter, statusCode int, message string) {
	// Ensure JSON content type is set
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	
	response := Response{
		Status: "error",
		Error:  message,
	}
	
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("Failed to encode JSON error response", "error", err)
	}
}