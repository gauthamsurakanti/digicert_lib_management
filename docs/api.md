# Library Management API

## Base URL
`/api/v1`

## Endpoints

### List Books
- **Method**: GET
- **Path**: `/books`
- **Response**: 200 OK
- **Body**: Array of Book objects
```json
[
  {
    "id": "uuid",
    "title": "string",
    "author": "string",
    "isbn": "string",
    "published_at": "YYYY-MM-DD",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
]