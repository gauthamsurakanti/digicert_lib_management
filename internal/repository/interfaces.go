package repository

import (
	"context"
	"library-management/internal/domain"
)

// BookRepository defines the interface for book data operations
type BookRepository interface {
	// Create creates a new book
	Create(ctx context.Context, book *domain.Book) (*domain.Book, error)
	
	// GetByID retrieves a book by its ID
	GetByID(ctx context.Context, id int) (*domain.Book, error)
	
	// GetAll retrieves all books with optional filtering
	GetAll(ctx context.Context, filter *domain.BookFilter) ([]*domain.Book, error)
	
	// Update updates an existing book
	Update(ctx context.Context, book *domain.Book) (*domain.Book, error)
	
	// Delete deletes a book by its ID
	Delete(ctx context.Context, id int) error
	
	// GetByISBN retrieves a book by its ISBN
	GetByISBN(ctx context.Context, isbn string) (*domain.Book, error)
	
	// Count returns the total number of books with optional filtering
	Count(ctx context.Context, filter *domain.BookFilter) (int, error)
}