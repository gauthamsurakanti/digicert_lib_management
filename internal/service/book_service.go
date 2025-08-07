package service

import (
	"context"
	"fmt"

	"library-management/internal/domain"
	"library-management/internal/repository"
)

type bookService struct {
	repo repository.BookRepository
}

// NewBookService creates a new book service
func NewBookService(repo repository.BookRepository) BookService {
	return &bookService{
		repo: repo,
	}
}

// CreateBook creates a new book
func (s *bookService) CreateBook(ctx context.Context, req *domain.CreateBookRequest) (*domain.Book, error) {
	// Validate the request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Check if a book with this ISBN already exists
	existingBook, err := s.repo.GetByISBN(ctx, req.ISBN)
	if err == nil && existingBook != nil {
		return nil, fmt.Errorf("book with ISBN %s already exists", req.ISBN)
	}

	// Convert request to domain model
	book := req.ToBook()

	// Create the book
	createdBook, err := s.repo.Create(ctx, book)
	if err != nil {
		return nil, fmt.Errorf("failed to create book: %w", err)
	}

	return createdBook, nil
}

// GetBookByID retrieves a book by its ID
func (s *bookService) GetBookByID(ctx context.Context, id int) (*domain.Book, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid book ID: %d", id)
	}

	book, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get book: %w", err)
	}

	return book, nil
}

// GetAllBooks retrieves all books with optional filtering
func (s *bookService) GetAllBooks(ctx context.Context, filter *domain.BookFilter) ([]*domain.Book, error) {
	books, err := s.repo.GetAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get books: %w", err)
	}

	// If no books found, return empty slice instead of nil
	if books == nil {
		books = []*domain.Book{}
	}

	return books, nil
}

// UpdateBook updates an existing book
func (s *bookService) UpdateBook(ctx context.Context, id int, req *domain.UpdateBookRequest) (*domain.Book, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid book ID: %d", id)
	}

	// Get the existing book
	existingBook, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing book: %w", err)
	}

	// Check if ISBN is being updated and conflicts with another book
	if req.ISBN != nil && *req.ISBN != existingBook.ISBN {
		conflictingBook, err := s.repo.GetByISBN(ctx, *req.ISBN)
		if err == nil && conflictingBook != nil && conflictingBook.ID != id {
			return nil, fmt.Errorf("book with ISBN %s already exists", *req.ISBN)
		}
	}

	// Apply updates to the existing book
	req.ApplyTo(existingBook)

	// Update the book
	updatedBook, err := s.repo.Update(ctx, existingBook)
	if err != nil {
		return nil, fmt.Errorf("failed to update book: %w", err)
	}

	return updatedBook, nil
}

// DeleteBook deletes a book by its ID
func (s *bookService) DeleteBook(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid book ID: %d", id)
	}

	// Check if book exists before attempting to delete
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("book not found: %w", err)
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete book: %w", err)
	}

	return nil
}

// GetBookByISBN retrieves a book by its ISBN
func (s *bookService) GetBookByISBN(ctx context.Context, isbn string) (*domain.Book, error) {
	if isbn == "" {
		return nil, fmt.Errorf("ISBN cannot be empty")
	}

	book, err := s.repo.GetByISBN(ctx, isbn)
	if err != nil {
		return nil, fmt.Errorf("failed to get book by ISBN: %w", err)
	}

	return book, nil
}

// GetBooksCount returns the total number of books with optional filtering
func (s *bookService) GetBooksCount(ctx context.Context, filter *domain.BookFilter) (int, error) {
	count, err := s.repo.Count(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to get books count: %w", err)
	}

	return count, nil
}