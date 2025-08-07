-- Drop trigger
DROP TRIGGER IF EXISTS update_books_updated_at ON books;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_books_search;
DROP INDEX IF EXISTS idx_books_author;
DROP INDEX IF EXISTS idx_books_genre;
DROP INDEX IF EXISTS idx_books_available;
DROP INDEX IF EXISTS idx_books_title;
DROP INDEX IF EXISTS idx_books_isbn;

-- Drop table
DROP TABLE IF EXISTS books;