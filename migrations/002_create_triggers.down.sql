-- Drop trigger
DROP TRIGGER IF EXISTS update_books_updated_at ON books;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();