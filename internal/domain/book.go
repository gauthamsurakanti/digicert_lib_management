package domain

import (
	"errors"
	"time"
)

// Book represents a book in the library
type Book struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Author      string    `json:"author" db:"author"`
	ISBN        string    `json:"isbn" db:"isbn"`
	Publisher   string    `json:"publisher" db:"publisher"`
	PublishYear int       `json:"publish_year" db:"publish_year"`
	Genre       string    `json:"genre" db:"genre"`
	Pages       int       `json:"pages" db:"pages"`
	Available   bool      `json:"available" db:"available"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateBookRequest represents the request payload for creating a book
type CreateBookRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Author      string `json:"author" validate:"required,min=1,max=255"`
	ISBN        string `json:"isbn" validate:"required,isbn"`
	Publisher   string `json:"publisher" validate:"required,min=1,max=255"`
	PublishYear int    `json:"publish_year" validate:"required,min=1000,max=2030"`
	Genre       string `json:"genre" validate:"required,min=1,max=100"`
	Pages       int    `json:"pages" validate:"required,min=1"`
	Description string `json:"description" validate:"max=1000"`
}

// UpdateBookRequest represents the request payload for updating a book
type UpdateBookRequest struct {
	Title       *string `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Author      *string `json:"author,omitempty" validate:"omitempty,min=1,max=255"`
	ISBN        *string `json:"isbn,omitempty" validate:"omitempty,isbn"`
	Publisher   *string `json:"publisher,omitempty" validate:"omitempty,min=1,max=255"`
	PublishYear *int    `json:"publish_year,omitempty" validate:"omitempty,min=1000,max=2030"`
	Genre       *string `json:"genre,omitempty" validate:"omitempty,min=1,max=100"`
	Pages       *int    `json:"pages,omitempty" validate:"omitempty,min=1"`
	Available   *bool   `json:"available,omitempty"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
}

// Validate validates the CreateBookRequest
func (r *CreateBookRequest) Validate() error {
	if r.Title == "" {
		return errors.New("title is required")
	}
	if r.Author == "" {
		return errors.New("author is required")
	}
	if r.ISBN == "" {
		return errors.New("ISBN is required")
	}
	if r.Publisher == "" {
		return errors.New("publisher is required")
	}
	if r.Genre == "" {
		return errors.New("genre is required")
	}
	if r.PublishYear < 1000 || r.PublishYear > 2030 {
		return errors.New("publish year must be between 1000 and 2030")
	}
	if r.Pages < 1 {
		return errors.New("pages must be greater than 0")
	}
	return nil
}

// ToBook converts CreateBookRequest to Book domain model
func (r *CreateBookRequest) ToBook() *Book {
	now := time.Now()
	return &Book{
		Title:       r.Title,
		Author:      r.Author,
		ISBN:        r.ISBN,
		Publisher:   r.Publisher,
		PublishYear: r.PublishYear,
		Genre:       r.Genre,
		Pages:       r.Pages,
		Available:   true, // Default to available
		Description: r.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// ApplyTo applies UpdateBookRequest changes to existing Book
func (r *UpdateBookRequest) ApplyTo(book *Book) {
	if r.Title != nil {
		book.Title = *r.Title
	}
	if r.Author != nil {
		book.Author = *r.Author
	}
	if r.ISBN != nil {
		book.ISBN = *r.ISBN
	}
	if r.Publisher != nil {
		book.Publisher = *r.Publisher
	}
	if r.PublishYear != nil {
		book.PublishYear = *r.PublishYear
	}
	if r.Genre != nil {
		book.Genre = *r.Genre
	}
	if r.Pages != nil {
		book.Pages = *r.Pages
	}
	if r.Available != nil {
		book.Available = *r.Available
	}
	if r.Description != nil {
		book.Description = *r.Description
	}
	book.UpdatedAt = time.Now()
}

// BookFilter represents filtering options for books
type BookFilter struct {
	Author    string `json:"author,omitempty"`
	Genre     string `json:"genre,omitempty"`
	Available *bool  `json:"available,omitempty"`
	Search    string `json:"search,omitempty"` // Search in title, author, or description
}