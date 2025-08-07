package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"library-management/internal/domain"
)

// MockBookRepository implements repository.BookRepository for testing
type MockBookRepository struct {
	books  map[int]*domain.Book
	nextID int
}

func NewMockBookRepository() *MockBookRepository {
	return &MockBookRepository{
		books:  make(map[int]*domain.Book),
		nextID: 1,
	}
}

func (m *MockBookRepository) Create(ctx context.Context, book *domain.Book) (*domain.Book, error) {
	// Check for duplicate ISBN
	for _, existingBook := range m.books {
		if existingBook.ISBN == book.ISBN {
			return nil, fmt.Errorf("book with ISBN %s already exists", book.ISBN)
		}
	}

	book.ID = m.nextID
	m.nextID++
	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()

	m.books[book.ID] = book
	return book, nil
}

func (m *MockBookRepository) GetByID(ctx context.Context, id int) (*domain.Book, error) {
	book, exists := m.books[id]
	if !exists {
		return nil, fmt.Errorf("book with ID %d not found", id)
	}
	return book, nil
}

func (m *MockBookRepository) GetAll(ctx context.Context, filter *domain.BookFilter) ([]*domain.Book, error) {
	var books []*domain.Book
	for _, book := range m.books {
		books = append(books, book)
	}
	return books, nil
}

func (m *MockBookRepository) Update(ctx context.Context, book *domain.Book) (*domain.Book, error) {
	_, exists := m.books[book.ID]
	if !exists {
		return nil, fmt.Errorf("book with ID %d not found", book.ID)
	}

	book.UpdatedAt = time.Now()
	m.books[book.ID] = book
	return book, nil
}

func (m *MockBookRepository) Delete(ctx context.Context, id int) error {
	_, exists := m.books[id]
	if !exists {
		return fmt.Errorf("book with ID %d not found", id)
	}

	delete(m.books, id)
	return nil
}

func (m *MockBookRepository) GetByISBN(ctx context.Context, isbn string) (*domain.Book, error) {
	for _, book := range m.books {
		if book.ISBN == isbn {
			return book, nil
		}
	}
	return nil, fmt.Errorf("book with ISBN %s not found", isbn)
}

func (m *MockBookRepository) Count(ctx context.Context, filter *domain.BookFilter) (int, error) {
	return len(m.books), nil
}

// Tests
func TestBookService_CreateBook(t *testing.T) {
	repo := NewMockBookRepository()
	service := NewBookService(repo)
	ctx := context.Background()

	t.Run("successful creation", func(t *testing.T) {
		req := &domain.CreateBookRequest{
			Title:       "Test Book",
			Author:      "Test Author",
			ISBN:        "978-1234567890",
			Publisher:   "Test Publisher",
			PublishYear: 2024,
			Genre:       "Test",
			Pages:       100,
			Description: "Test description",
		}

		book, err := service.CreateBook(ctx, req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if book.ID == 0 {
			t.Error("Expected book ID to be set")
		}

		if book.Title != req.Title {
			t.Errorf("Expected title %s, got %s", req.Title, book.Title)
		}
	})

	t.Run("duplicate ISBN", func(t *testing.T) {
		req1 := &domain.CreateBookRequest{
			Title:       "Book 1",
			Author:      "Author 1",
			ISBN:        "978-1111111111",
			Publisher:   "Publisher 1",
			PublishYear: 2024,
			Genre:       "Genre 1",
			Pages:       100,
		}

		req2 := &domain.CreateBookRequest{
			Title:       "Book 2",
			Author:      "Author 2",
			ISBN:        "978-1111111111", // Same ISBN
			Publisher:   "Publisher 2",
			PublishYear: 2024,
			Genre:       "Genre 2",
			Pages:       200,
		}

		// Create first book
		_, err := service.CreateBook(ctx, req1)
		if err != nil {
			t.Fatalf("Expected no error for first book, got %v", err)
		}

		// Try to create second book with same ISBN
		_, err = service.CreateBook(ctx, req2)
		if err == nil {
			t.Error("Expected error for duplicate ISBN")
		}
	})

	t.Run("validation error", func(t *testing.T) {
		req := &domain.CreateBookRequest{
			Title:       "", // Empty title should fail validation
			Author:      "Test Author",
			ISBN:        "978-1234567890",
			Publisher:   "Test Publisher",
			PublishYear: 2024,
			Genre:       "Test",
			Pages:       100,
		}

		_, err := service.CreateBook(ctx, req)
		if err == nil {
			t.Error("Expected validation error for empty title")
		}
	})
}

func TestBookService_GetBookByID(t *testing.T) {
	repo := NewMockBookRepository()
	service := NewBookService(repo)
	ctx := context.Background()

	// Create a book first
	req := &domain.CreateBookRequest{
		Title:       "Test Book",
		Author:      "Test Author",
		ISBN:        "978-1234567890",
		Publisher:   "Test Publisher",
		PublishYear: 2024,
		Genre:       "Test",
		Pages:       100,
	}

	createdBook, err := service.CreateBook(ctx, req)
	if err != nil {
		t.Fatalf("Failed to create test book: %v", err)
	}

	t.Run("successful retrieval", func(t *testing.T) {
		book, err := service.GetBookByID(ctx, createdBook.ID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if book.ID != createdBook.ID {
			t.Errorf("Expected ID %d, got %d", createdBook.ID, book.ID)
		}
	})

	t.Run("book not found", func(t *testing.T) {
		_, err := service.GetBookByID(ctx, 999)
		if err == nil {
			t.Error("Expected error for non-existent book")
		}
	})

	t.Run("invalid ID", func(t *testing.T) {
		_, err := service.GetBookByID(ctx, 0)
		if err == nil {
			t.Error("Expected error for invalid book ID")
		}
	})
}

func TestBookService_UpdateBook(t *testing.T) {
	repo := NewMockBookRepository()
	service := NewBookService(repo)
	ctx := context.Background()

	// Create a book first
	req := &domain.CreateBookRequest{
		Title:       "Original Title",
		Author:      "Original Author",
		ISBN:        "978-1234567890",
		Publisher:   "Original Publisher",
		PublishYear: 2024,
		Genre:       "Original Genre",
		Pages:       100,
	}

	createdBook, err := service.CreateBook(ctx, req)
	if err != nil {
		t.Fatalf("Failed to create test book: %v", err)
	}

	t.Run("successful update", func(t *testing.T) {
		newTitle := "Updated Title"
		updateReq := &domain.UpdateBookRequest{
			Title: &newTitle,
		}

		updatedBook, err := service.UpdateBook(ctx, createdBook.ID, updateReq)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if updatedBook.Title != newTitle {
			t.Errorf("Expected title %s, got %s", newTitle, updatedBook.Title)
		}

		// Original author should remain unchanged
		if updatedBook.Author != req.Author {
			t.Errorf("Expected author to remain %s, got %s", req.Author, updatedBook.Author)
		}
	})

	t.Run("book not found", func(t *testing.T) {
		newTitle := "Updated Title"
		updateReq := &domain.UpdateBookRequest{
			Title: &newTitle,
		}

		_, err := service.UpdateBook(ctx, 999, updateReq)
		if err == nil {
			t.Error("Expected error for non-existent book")
		}
	})
}

func TestBookService_DeleteBook(t *testing.T) {
	repo := NewMockBookRepository()
	service := NewBookService(repo)
	ctx := context.Background()

	// Create a book first
	req := &domain.CreateBookRequest{
		Title:       "Test Book",
		Author:      "Test Author",
		ISBN:        "978-1234567890",
		Publisher:   "Test Publisher",
		PublishYear: 2024,
		Genre:       "Test",
		Pages:       100,
	}

	createdBook, err := service.CreateBook(ctx, req)
	if err != nil {
		t.Fatalf("Failed to create test book: %v", err)
	}

	t.Run("successful deletion", func(t *testing.T) {
		err := service.DeleteBook(ctx, createdBook.ID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify book is deleted
		_, err = service.GetBookByID(ctx, createdBook.ID)
		if err == nil {
			t.Error("Expected error when getting deleted book")
		}
	})

	t.Run("book not found", func(t *testing.T) {
		err := service.DeleteBook(ctx, 999)
		if err == nil {
			t.Error("Expected error for non-existent book")
		}
	})
}
