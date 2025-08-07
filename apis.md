# Library Management API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
Currently, no authentication is required. In production, consider implementing JWT or API key authentication.

## Error Handling

All API responses follow this standard format:

### Success Response
```json
{
  "status": "success",
  "message": "Operation completed successfully",
  "data": { ... }
}
```

### Error Response
```json
{
  "status": "error",
  "error": "Error description"
}
```

## Endpoints

### 1. Health Check

**GET** `/health`

Check if the API service is running.

**Response:**
```json
{
  "status": "success",
  "message": "Service is healthy",
  "data": {
    "status": "ok",
    "service": "library-management"
  }
}
```

---

### 2. List All Books

**GET** `/api/v1/books`

Retrieve all books with optional filtering.

**Query Parameters:**
- `author` (string, optional) - Filter by author (partial match, case-insensitive)
- `genre` (string, optional) - Filter by genre (exact match, case-insensitive)
- `available` (boolean, optional) - Filter by availability (true/false)
- `search` (string, optional) - Search in title, author, or description

**Examples:**
```bash
# Get all books
GET /api/v1/books

# Filter by genre
GET /api/v1/books?genre=Programming

# Search for books
GET /api/v1/books?search=golang

# Get available books by author
GET /api/v1/books?author=Martin&available=true
```

**Response:**
```json
{
  "status": "success",
  "message": "Books retrieved successfully",
  "data": {
    "books": [
      {
        "id": 1,
        "title": "The Go Programming Language",
        "author": "Alan Donovan, Brian Kernighan",
        "isbn": "978-0134190440",
        "publisher": "Addison-Wesley",
        "publish_year": 2015,
        "genre": "Programming",
        "pages": 380,
        "available": true,
        "description": "The authoritative resource...",
        "created_at": "2024-01-01T10:00:00Z",
        "updated_at": "2024-01-01T10:00:00Z"
      }
    ],
    "meta": {
      "total": 8,
      "count": 1
    }
  }
}
```

---

### 3. Get Book by ID

**GET** `/api/v1/books/{id}`

Retrieve a specific book by its ID.

**Path Parameters:**
- `id` (integer, required) - Book ID

**Response:**
```json
{
  "status": "success",
  "message": "Book retrieved successfully",
  "data": {
    "id": 1,
    "title": "The Go Programming Language",
    "author": "Alan Donovan, Brian Kernighan",
    "isbn": "978-0134190440",
    "publisher": "Addison-Wesley",
    "publish_year": 2015,
    "genre": "Programming",
    "pages": 380,
    "available": true,
    "description": "The authoritative resource...",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

**Error Response (404):**
```json
{
  "status": "error",
  "error": "Book not found"
}
```

---

### 4. Create New Book

**POST** `/api/v1/books`

Create a new book in the library.

**Request Body:**
```json
{
  "title": "Book Title",
  "author": "Author Name",
  "isbn": "978-1234567890",
  "publisher": "Publisher Name",
  "publish_year": 2024,
  "genre": "Genre",
  "pages": 250,
  "description": "Book description (optional)"
}
```

**Validation Rules:**
- `title`: Required, 1-255 characters
- `author`: Required, 1-255 characters
- `isbn`: Required, must be unique
- `publisher`: Required, 1-255 characters
- `publish_year`: Required, between 1000-2030
- `genre`: Required, 1-100 characters
- `pages`: Required, must be > 0
- `description`: Optional, max 1000 characters

**Response (201):**
```json
{
  "status": "success",
  "message": "Book created successfully",
  "data": {
    "id": 9,
    "title": "Book Title",
    "author": "Author Name",
    "isbn": "978-1234567890",
    "publisher": "Publisher Name",
    "publish_year": 2024,
    "genre": "Genre",
    "pages": 250,
    "available": true,
    "description": "Book description",
    "created_at": "2024-01-02T10:00:00Z",
    "updated_at": "2024-01-02T10:00:00Z"
  }
}
```

**Error Response (400):**
```json
{
  "status": "error",
  "error": "book with ISBN 978-1234567890 already exists"
}
```

---

### 5. Update Book

**PUT** `/api/v1/books/{id}`

Update an existing book. Only provided fields will be updated.

**Path Parameters:**
- `id` (integer, required) - Book ID

**Request Body (all fields optional):**
```json
{
  "title": "Updated Title",
  "author": "Updated Author",
  "isbn": "978-0987654321",
  "publisher": "Updated Publisher",
  "publish_year": 2025,
  "genre": "Updated Genre",
  "pages": 300,
  "available": false,
  "description": "Updated description"
}
```

**Response (200):**
```json
{
  "status": "success",
  "message": "Book updated successfully",
  "data": {
    "id": 1,
    "title": "Updated Title",
    "author": "Updated Author",
    "isbn": "978-0987654321",
    "publisher": "Updated Publisher",
    "publish_year": 2025,
    "genre": "Updated Genre",
    "pages": 300,
    "available": false,
    "description": "Updated description",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-02T15:30:00Z"
  }
}
```

---

### 6. Delete Book

**DELETE** `/api/v1/books/{id}`

Delete a book from the library.

**Path Parameters:**
- `id` (integer, required) - Book ID

**Response (200):**
```json
{
  "status": "success",
  "message": "Book deleted successfully"
}
```

**Error Response (404):**
```json
{
  "status": "error",
  "error": "Book not found"
}
```

---

### 7. Get Book by ISBN

**GET** `/api/v1/books/isbn/{isbn}`

Retrieve a book by its ISBN.

**Path Parameters:**
- `isbn` (string, required) - Book ISBN

**Response:**
```json
{
  "status": "success",
  "message": "Book retrieved successfully",
  "data": {
    "id": 1,
    "title": "The Go Programming Language",
    "author": "Alan Donovan, Brian Kernighan",
    "isbn": "978-0134190440",
    "publisher": "Addison-Wesley",
    "publish_year": 2015,
    "genre": "Programming",
    "pages": 380,
    "available": true,
    "description": "The authoritative resource...",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

## HTTP Status Codes

| Status Code | Description |
|-------------|-------------|
| 200 | OK - Request successful |
| 201 | Created - Resource created successfully |
| 400 | Bad Request - Invalid input or validation error |
| 404 | Not Found - Resource not found |
| 500 | Internal Server Error - Server error |

## Rate Limiting

Currently no rate limiting is implemented. For production use, consider implementing rate limiting middleware.

## Examples with cURL

### Create a Book
```bash
curl -X POST http://localhost:8080/api/v1/books \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Effective Go",
    "author": "The Go Team",
    "isbn": "978-1234567890",
    "publisher": "Google",
    "publish_year": 2022,
    "genre": "Programming",
    "pages": 200,
    "description": "Best practices for Go programming"
  }'
```

### Update Book Availability
```bash
curl -X PUT http://localhost:8080/api/v1/books/1 \
  -H "Content-Type: application/json" \
  -d '{"available": false}'
```

### Search Books
```bash
# Search for Go books
curl "http://localhost:8080/api/v1/books?search=go"

# Get programming books
curl "http://localhost:8080/api/v1/books?genre=Programming"

# Get available books by author
curl "http://localhost:8080/api/v1/books?author=Martin&available=true"
```

## Database Schema

### Books Table
```sql
CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    isbn VARCHAR(20) UNIQUE NOT NULL,
    publisher VARCHAR(255) NOT NULL,
    publish_year INTEGER NOT NULL CHECK (publish_year >= 1000 AND publish_year <= 2030),
    genre VARCHAR(100) NOT NULL,
    pages INTEGER NOT NULL CHECK (pages > 0),
    available BOOLEAN NOT NULL DEFAULT true,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### Indexes
- Primary key on `id`
- Unique index on `isbn`
- Indexes on `author`, `genre`, `available`, `title`
- Full-text search index on `title`, `author`, `description`

## Testing

### Unit Tests
```bash
make test
```

### API Integration Tests
```bash
# Start the server first
make docker-up

# Test endpoints
make test-api
```

### Manual Testing
Use the provided web interface at `http://localhost:8080` or use tools like Postman/Insomnia with the API endpoints.