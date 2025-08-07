# 📚 Library Management REST API

A production-ready library management system built with **Go**, featuring clean architecture, comprehensive testing, and full containerization.

## 🌟 Features

- **Full CRUD Operations** for books
- **Clean Architecture** with proper separation of concerns
- **PostgreSQL** database with migrations
- **Docker & Docker Compose** for easy deployment
- **RESTful API** with proper HTTP status codes
- **Input Validation** and error handling
- **Structured Logging** with JSON output
- **Health Check** endpoints
- **CORS Support** for web frontends
- **Database Indexing** for optimal performance
- **Graceful Shutdown** handling

## 🚀 Quick Start

### Prerequisites
- Go 1.23+
- Docker & Docker Compose
- Make (optional, for convenience commands)

### 1. Clone and Setup
```bash
git clone https://github.com/gauthamsurakanti/digicert_lib_management.git
cd digicert_lib_management
```

### 2. Start Docker
1. Start Docker Desktop from /Applications (Recommended)

### 3. Start with Docker Compose (Recommended)
```bash
docker-compose up -d
#or 
make docker-up
```

### 3. Verify Installation
```bash
curl http://localhost:8080/health
```

## 📡 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/api/v1/books` | List all books |
| POST | `/api/v1/books` | Create a new book |
| GET | `/api/v1/books/{id}` | Get book by ID |
| PUT | `/api/v1/books/{id}` | Update book |
| DELETE | `/api/v1/books/{id}` | Delete book |
| GET | `/api/v1/books/isbn/{isbn}` | Get book by ISBN |

### Query Parameters (for GET /api/v1/books)
- `author` - Filter by author (partial match)
- `genre` - Filter by genre (exact match)
- `available` - Filter by availability (true/false)
- `search` - Search in title, author, or description

## 📝 API Examples

### Create a Book
```bash
curl -X POST http://localhost:8080/api/v1/books \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Applied Cryptography",
    "author": "Bruce Schneier",
    "isbn": "978-1119096726",
    "publisher": "Wiley",
    "publish_year": 2015,
    "genre": "Cybersecurity",
    "pages": 784,
    "description": "Classic guide covering practical cryptographic protocols and algorithms."
  }'
```

### Get All Books
```bash
curl http://localhost:8080/api/v1/books
```

### Filter Books
```bash
# Get programming books
curl "http://localhost:8080/api/v1/books?genre=Programming"

# Get available books by specific author
curl "http://localhost:8080/api/v1/books?author=Robert&available=true"

# Search books
curl "http://localhost:8080/api/v1/books?search=programming"
```

### Update a Book
```bash
curl -X PUT http://localhost:8080/api/v1/books/1 \
  -H "Content-Type: application/json" \
  -d '{
    "available": false,
    "description": "Updated description"
  }'
```

### Delete a Book
```bash
curl -X DELETE http://localhost:8080/api/v1/books/1
```

## 🏗️ Architecture

This project follows **Clean Architecture** principles:

```
├── cmd/api/                 # Application entry point
├── internal/
│   ├── domain/             # Business entities & validation
│   ├── service/            # Business logic layer
│   ├── repository/         # Data access layer
│   │   └── postgres/       # PostgreSQL implementation
│   ├── handler/            # HTTP handlers (controllers)
│   ├── config/             # Configuration management
│   └── database/           # Database connection & migrations
├── pkg/                    # Shared packages
├── migrations/             # Database migrations
└── web/                    # Web UI assets
```

### Key Design Decisions

1. **Interface-Based Design** - All layers communicate through interfaces for better testability
2. **Dependency Injection** - Clean dependency management from main.go
3. **Repository Pattern** - Abstracts data access for easy database switching
4. **Service Layer** - Contains business logic and validation
5. **Middleware Chain** - CORS, logging, and JSON content handling
6. **Graceful Shutdown** - Handles SIGINT/SIGTERM properly

## 🐳 Docker Setup

### Development with Docker Compose
```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Production Deployment
```bash
# Build optimized image
docker build -t library-management:prod .

# Run with production settings
docker run -p 8080:8080 \
  -e DATABASE_URL=your-prod-db-url \
  -e ENVIRONMENT=production \
  library-management:prod
```

## 🗄️ Database

### Migrations
```bash
# Run migrations up
make migrate-up

# Run migrations down  
make migrate-down

# Create new migration
make migrate-create name=add_new_field
```

### Schema
The `books` table includes:
- **Primary Key**: Auto-incrementing ID
- **Unique Constraint**: ISBN
- **Indexes**: On author, genre, availability, title, ISBN
- **Full-Text Search**: PostgreSQL GIN index
- **Timestamps**: Automatic created_at/updated_at handling

## 🧪 Testing

### Run Tests
```bash
make test
make test-coverage
```

### API Testing
```bash
# Test all endpoints
make test-api

# Create test book
make create-book
```

## 🔧 Development

### Local Development
```bash
# Setup environment
make dev-setup

# Install dependencies
make deps

# Format code
make fmt

# Run linter (requires golangci-lint)
make lint

# Start development server
make run
```

### Adding New Features
1. Define domain models in `internal/domain/`
2. Create repository interfaces in `internal/repository/interfaces.go`
3. Implement repository in `internal/repository/postgres/`
4. Add business logic in `internal/service/`
5. Create HTTP handlers in `internal/handler/`
6. Register routes in `internal/handler/routes.go`

## 🚀 Production Considerations

### Security
- Input validation on all endpoints
- SQL injection prevention with parameterized queries
- Non-root user in Docker container
- Environment-based configuration

### Performance
- Database connection pooling
- Proper indexing strategy
- Efficient query patterns
- Graceful shutdown handling

### Monitoring
- Structured JSON logging
- Health check endpoints
- Request/response logging middleware

## 📊 Sample Data

The application comes with sample programming books including:
- The Go Programming Language
- Clean Code
- Design Patterns
- The Pragmatic Programmer
- Microservices Patterns

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

**Built with ❤️ using Go, PostgreSQL, and Docker**