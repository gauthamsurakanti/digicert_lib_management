package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Connect establishes a connection to PostgreSQL database
func Connect(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return db, nil
}

// InitializeDatabase creates the database schema and sample data
func InitializeDatabase(db *sql.DB) error {
	fmt.Println("Initializing database schema...")

	// Create books table
	if err := createBooksTable(db); err != nil {
		return fmt.Errorf("failed to create books table: %w", err)
	}

	// Create indexes
	if err := createIndexes(db); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	// Create triggers for auto-updating timestamps
	if err := createTriggers(db); err != nil {
		fmt.Printf("Warning: failed to create triggers: %v\n", err)
	}

	// Insert sample data if table is empty
	if err := insertSampleData(db); err != nil {
		return fmt.Errorf("failed to insert sample data: %w", err)
	}

	fmt.Println("Database initialization completed successfully")
	return nil
}

// createBooksTable creates the books table
func createBooksTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS books (
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
	);`

	if _, err := db.Exec(query); err != nil {
		return err
	}

	fmt.Println("Books table created successfully")
	return nil
}

// createIndexes creates database indexes for better performance
func createIndexes(db *sql.DB) error {
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_books_author ON books(author);",
		"CREATE INDEX IF NOT EXISTS idx_books_genre ON books(genre);",
		"CREATE INDEX IF NOT EXISTS idx_books_available ON books(available);",
		"CREATE INDEX IF NOT EXISTS idx_books_title ON books(title);",
		"CREATE INDEX IF NOT EXISTS idx_books_isbn ON books(isbn);",
	}

	for _, indexQuery := range indexes {
		if _, err := db.Exec(indexQuery); err != nil {
			fmt.Printf("Warning: failed to create index: %v\n", err)
		}
	}

	fmt.Println("Database indexes created successfully")
	return nil
}

// createTriggers creates database triggers for automatic timestamp updates
func createTriggers(db *sql.DB) error {
	// Create trigger function
	functionQuery := `
	CREATE OR REPLACE FUNCTION update_updated_at_column()
	RETURNS TRIGGER AS '
	BEGIN
		NEW.updated_at = CURRENT_TIMESTAMP;
		RETURN NEW;
	END;
	' LANGUAGE plpgsql;`

	if _, err := db.Exec(functionQuery); err != nil {
		return err
	}

	// Create trigger
	triggerQuery := `
	DROP TRIGGER IF EXISTS update_books_updated_at ON books;
	CREATE TRIGGER update_books_updated_at 
		BEFORE UPDATE ON books 
		FOR EACH ROW 
		EXECUTE FUNCTION update_updated_at_column();`

	if _, err := db.Exec(triggerQuery); err != nil {
		return err
	}

	fmt.Println("Database triggers created successfully")
	return nil
}

// insertSampleData inserts sample books if the table is empty
func insertSampleData(db *sql.DB) error {
	// Check if sample data already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM books").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		fmt.Printf("Found %d existing books, skipping sample data insertion\n", count)
		return nil
	}

	fmt.Println("Inserting sample data...")

	// Define sample books
	sampleBooks := []struct {
		title, author, isbn, publisher, genre, description string
		publishYear, pages                                 int
	}{
		{
			title:       "The Go Programming Language",
			author:      "Alan Donovan, Brian Kernighan",
			isbn:        "978-0134190440",
			publisher:   "Addison-Wesley",
			publishYear: 2015,
			genre:       "Programming",
			pages:       380,
			description: "The authoritative resource to writing clear and idiomatic Go to solve real-world problems.",
		},
		{
			title:       "Clean Code",
			author:      "Robert C. Martin",
			isbn:        "978-0132350884",
			publisher:   "Prentice Hall",
			publishYear: 2008,
			genre:       "Programming",
			pages:       464,
			description: "A handbook of agile software craftsmanship.",
		},
		{
			title:       "Design Patterns",
			author:      "Gang of Four",
			isbn:        "978-0201633610",
			publisher:   "Addison-Wesley",
			publishYear: 1994,
			genre:       "Programming",
			pages:       395,
			description: "Elements of reusable object-oriented software.",
		},
		{
			title:       "The Pragmatic Programmer",
			author:      "David Thomas, Andrew Hunt",
			isbn:        "978-0135957059",
			publisher:   "Addison-Wesley",
			publishYear: 2019,
			genre:       "Programming",
			pages:       352,
			description: "Your journey to mastery.",
		},
		{
			title:       "Microservices Patterns",
			author:      "Chris Richardson",
			isbn:        "978-1617294549",
			publisher:   "Manning Publications",
			publishYear: 2018,
			genre:       "Architecture",
			pages:       520,
			description: "With examples in Java.",
		},
		{
			title:       "Building Microservices",
			author:      "Sam Newman",
			isbn:        "978-1491950357",
			publisher:   "O'Reilly Media",
			publishYear: 2015,
			genre:       "Architecture",
			pages:       280,
			description: "Designing fine-grained systems.",
		},
		{
			title:       "Domain-Driven Design",
			author:      "Eric Evans",
			isbn:        "978-0321125217",
			publisher:   "Addison-Wesley",
			publishYear: 2003,
			genre:       "Architecture",
			pages:       560,
			description: "Tackling complexity in the heart of software.",
		},
		{
			title:       "The Art of Computer Programming",
			author:      "Donald Knuth",
			isbn:        "978-0201896831",
			publisher:   "Addison-Wesley",
			publishYear: 1997,
			genre:       "Computer Science",
			pages:       650,
			description: "Volume 1: Fundamental Algorithms.",
		},
	}

	// Insert each book
	insertQuery := `
	INSERT INTO books (title, author, isbn, publisher, publish_year, genre, pages, description) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	for _, book := range sampleBooks {
		_, err := db.Exec(insertQuery,
			book.title,
			book.author,
			book.isbn,
			book.publisher,
			book.publishYear,
			book.genre,
			book.pages,
			book.description,
		)
		if err != nil {
			fmt.Printf("Warning: failed to insert book '%s': %v\n", book.title, err)
		}
	}

	fmt.Printf("Sample data inserted successfully (%d books)\n", len(sampleBooks))
	return nil
}
