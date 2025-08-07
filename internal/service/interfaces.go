package service

import (
	"context"
	"library-management/internal/domain"
)

// BookService defines the interface for book business logic
type BookService interface {
	// CreateBook creates a new book
	CreateBook(ctx context.Context, req *domain.CreateBookRequest) (*domain.Book, error)
	
	// GetBookByID retrieves a book by its ID
	GetBookByID(ctx context.Context, id int) (*domain.Book, error)
	
	// GetAllBooks retrieves all books with optional filtering
	GetAllBooks(ctx context.Context, filter *domain.BookFilter) ([]*domain.Book, error)
	
	// UpdateBook updates an existing book
	UpdateBook(ctx context.Context, id int, req *domain.UpdateBookRequest) (*domain.Book, error)
	
	// DeleteBook deletes a book by its ID
	DeleteBook(ctx context.Context, id int) error
	
	// GetBookByISBN retrieves a book by its ISBN
	GetBookByISBN(ctx context.Context, isbn string) (*domain.Book, error)
	
	// GetBooksCount returns the total number of books with optional filtering
	GetBooksCount(ctx context.Context, filter *domain.BookFilter) (int, error)
}