package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"library-management/internal/domain"
	"library-management/internal/repository"
)

type bookRepository struct {
	db *sql.DB
}

// NewBookRepository creates a new PostgreSQL book repository
func NewBookRepository(db *sql.DB) repository.BookRepository {
	return &bookRepository{db: db}
}

// Create creates a new book
func (r *bookRepository) Create(ctx context.Context, book *domain.Book) (*domain.Book, error) {
	query := `
		INSERT INTO books (title, author, isbn, publisher, publish_year, genre, pages, available, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(
		ctx, query,
		book.Title, book.Author, book.ISBN, book.Publisher,
		book.PublishYear, book.Genre, book.Pages, book.Available,
		book.Description, book.CreatedAt, book.UpdatedAt,
	).Scan(&book.ID, &book.CreatedAt, &book.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create book: %w", err)
	}

	return book, nil
}

// GetByID retrieves a book by its ID
func (r *bookRepository) GetByID(ctx context.Context, id int) (*domain.Book, error) {
	query := `
		SELECT id, title, author, isbn, publisher, publish_year, genre, 
		       pages, available, description, created_at, updated_at
		FROM books 
		WHERE id = $1`

	book := &domain.Book{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&book.ID, &book.Title, &book.Author, &book.ISBN,
		&book.Publisher, &book.PublishYear, &book.Genre,
		&book.Pages, &book.Available, &book.Description,
		&book.CreatedAt, &book.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("book with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get book: %w", err)
	}

	return book, nil
}

// GetAll retrieves all books with optional filtering
func (r *bookRepository) GetAll(ctx context.Context, filter *domain.BookFilter) ([]*domain.Book, error) {
	query := `
		SELECT id, title, author, isbn, publisher, publish_year, genre, 
		       pages, available, description, created_at, updated_at
		FROM books`

	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter != nil {
		if filter.Author != "" {
			conditions = append(conditions, fmt.Sprintf("LOWER(author) LIKE LOWER($%d)", argIndex))
			args = append(args, "%"+filter.Author+"%")
			argIndex++
		}

		if filter.Genre != "" {
			conditions = append(conditions, fmt.Sprintf("LOWER(genre) = LOWER($%d)", argIndex))
			args = append(args, filter.Genre)
			argIndex++
		}

		if filter.Available != nil {
			conditions = append(conditions, fmt.Sprintf("available = $%d", argIndex))
			args = append(args, *filter.Available)
			argIndex++
		}

		if filter.Search != "" {
			searchCondition := fmt.Sprintf(`(
				LOWER(title) LIKE LOWER($%d) OR 
				LOWER(author) LIKE LOWER($%d) OR 
				LOWER(description) LIKE LOWER($%d)
			)`, argIndex, argIndex, argIndex)
			conditions = append(conditions, searchCondition)
			args = append(args, "%"+filter.Search+"%")
			argIndex++
		}

		if len(conditions) > 0 {
			query += " WHERE " + strings.Join(conditions, " AND ")
		}
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query books: %w", err)
	}
	defer rows.Close()

	var books []*domain.Book
	for rows.Next() {
		book := &domain.Book{}
		err := rows.Scan(
			&book.ID, &book.Title, &book.Author, &book.ISBN,
			&book.Publisher, &book.PublishYear, &book.Genre,
			&book.Pages, &book.Available, &book.Description,
			&book.CreatedAt, &book.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan book: %w", err)
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return books, nil
}

// Update updates an existing book
func (r *bookRepository) Update(ctx context.Context, book *domain.Book) (*domain.Book, error) {
	query := `
		UPDATE books 
		SET title = $2, author = $3, isbn = $4, publisher = $5, 
		    publish_year = $6, genre = $7, pages = $8, available = $9, 
		    description = $10, updated_at = $11
		WHERE id = $1
		RETURNING updated_at`

	err := r.db.QueryRowContext(
		ctx, query,
		book.ID, book.Title, book.Author, book.ISBN,
		book.Publisher, book.PublishYear, book.Genre,
		book.Pages, book.Available, book.Description, book.UpdatedAt,
	).Scan(&book.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("book with ID %d not found", book.ID)
		}
		return nil, fmt.Errorf("failed to update book: %w", err)
	}

	return book, nil
}

// Delete deletes a book by its ID
func (r *bookRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM books WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete book: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("book with ID %d not found", id)
	}

	return nil
}

// GetByISBN retrieves a book by its ISBN
func (r *bookRepository) GetByISBN(ctx context.Context, isbn string) (*domain.Book, error) {
	query := `
		SELECT id, title, author, isbn, publisher, publish_year, genre, 
		       pages, available, description, created_at, updated_at
		FROM books 
		WHERE isbn = $1`

	book := &domain.Book{}
	err := r.db.QueryRowContext(ctx, query, isbn).Scan(
		&book.ID, &book.Title, &book.Author, &book.ISBN,
		&book.Publisher, &book.PublishYear, &book.Genre,
		&book.Pages, &book.Available, &book.Description,
		&book.CreatedAt, &book.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("book with ISBN %s not found", isbn)
		}
		return nil, fmt.Errorf("failed to get book by ISBN: %w", err)
	}

	return book, nil
}

// Count returns the total number of books with optional filtering
func (r *bookRepository) Count(ctx context.Context, filter *domain.BookFilter) (int, error) {
	query := "SELECT COUNT(*) FROM books"

	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter != nil {
		if filter.Author != "" {
			conditions = append(conditions, fmt.Sprintf("LOWER(author) LIKE LOWER($%d)", argIndex))
			args = append(args, "%"+filter.Author+"%")
			argIndex++
		}

		if filter.Genre != "" {
			conditions = append(conditions, fmt.Sprintf("LOWER(genre) = LOWER($%d)", argIndex))
			args = append(args, filter.Genre)
			argIndex++
		}

		if filter.Available != nil {
			conditions = append(conditions, fmt.Sprintf("available = $%d", argIndex))
			args = append(args, *filter.Available)
			argIndex++
		}

		if filter.Search != "" {
			searchCondition := fmt.Sprintf(`(
				LOWER(title) LIKE LOWER($%d) OR 
				LOWER(author) LIKE LOWER($%d) OR 
				LOWER(description) LIKE LOWER($%d)
			)`, argIndex, argIndex, argIndex)
			conditions = append(conditions, searchCondition)
			args = append(args, "%"+filter.Search+"%")
			argIndex++
		}

		if len(conditions) > 0 {
			query += " WHERE " + strings.Join(conditions, " AND ")
		}
	}

	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count books: %w", err)
	}

	return count, nil
}