-- Create books table
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
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_books_author ON books(author);
CREATE INDEX IF NOT EXISTS idx_books_genre ON books(genre);
CREATE INDEX IF NOT EXISTS idx_books_available ON books(available);
CREATE INDEX IF NOT EXISTS idx_books_title ON books(title);
CREATE INDEX IF NOT EXISTS idx_books_isbn ON books(isbn);

-- Create a function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$ language 'plpgsql';

-- Create trigger to automatically update updated_at
DROP TRIGGER IF EXISTS update_books_updated_at ON books;
CREATE TRIGGER update_books_updated_at 
    BEFORE UPDATE ON books 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Insert sample data (using correct column names)
INSERT INTO books (title, author, isbn, publisher, publish_year, genre, pages, description) VALUES
('The Go Programming Language', 'Alan Donovan, Brian Kernighan', '978-0134190440', 'Addison-Wesley', 2015, 'Programming', 380, 'The authoritative resource to writing clear and idiomatic Go to solve real-world problems.'),
('Clean Code', 'Robert C. Martin', '978-0132350884', 'Prentice Hall', 2008, 'Programming', 464, 'A handbook of agile software craftsmanship.'),
('Design Patterns', 'Gang of Four', '978-0201633610', 'Addison-Wesley', 1994, 'Programming', 395, 'Elements of reusable object-oriented software.'),
('The Pragmatic Programmer', 'David Thomas, Andrew Hunt', '978-0135957059', 'Addison-Wesley', 2019, 'Programming', 352, 'Your journey to mastery.'),
('Microservices Patterns', 'Chris Richardson', '978-1617294549', 'Manning Publications', 2018, 'Architecture', 520, 'With examples in Java.'),
('Building Microservices', 'Sam Newman', '978-1491950357', 'O''Reilly Media', 2015, 'Architecture', 280, 'Designing fine-grained systems.'),
('Domain-Driven Design', 'Eric Evans', '978-0321125217', 'Addison-Wesley', 2003, 'Architecture', 560, 'Tackling complexity in the heart of software.'),
('The Art of Computer Programming', 'Donald Knuth', '978-0201896831', 'Addison-Wesley', 1997, 'Computer Science', 650, 'Volume 1: Fundamental Algorithms.')
ON CONFLICT (isbn) DO NOTHING;

-- Create full-text search index (PostgreSQL specific)
CREATE INDEX IF NOT EXISTS idx_books_search ON books USING gin(to_tsvector('english', title || ' ' || author || ' ' || description));