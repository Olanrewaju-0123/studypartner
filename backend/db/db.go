package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

func Initialize(databaseURL string) (*sql.DB, error) {
	var err error
	DB, err = sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Enable pgvector extension
	if _, err = DB.Exec("CREATE EXTENSION IF NOT EXISTS vector"); err != nil {
		log.Printf("Warning: Could not create vector extension: %v", err)
	}

	log.Println("Database connected successfully")
	return DB, nil
}

func RunMigrations(db *sql.DB) error {
	// Check if vector extension is available
	var vectorAvailable bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_extension WHERE extname = 'vector')").Scan(&vectorAvailable)
	if err != nil {
		log.Printf("Warning: Could not check vector extension: %v", err)
		vectorAvailable = false
	}

	migrations := []string{
		createUsersTable,
		createSummariesTable,
		createFlashcardsTable,
		createQuizzesTable,
		createStudySessionsTable,
	}

	// Add notes table with or without vector support
	if vectorAvailable {
		migrations = append([]string{createNotesTableWithVector}, migrations...)
		migrations = append(migrations, createIndexesWithVector)
	} else {
		migrations = append([]string{createNotesTableWithoutVector}, migrations...)
		migrations = append(migrations, createIndexesWithoutVector)
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("failed to run migration: %w", err)
		}
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// SeedTestData creates a test user for demo purposes
func SeedTestData(db *sql.DB) error {
	// Check if test user already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", "demo@studypartner.com").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check test user: %w", err)
	}

	if count > 0 {
		log.Println("Test user already exists")
		return nil
	}

	// Hash password for test user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("demo123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash test password: %w", err)
	}

	// Create test user (password: "demo123")
	_, err = db.Exec(`
		INSERT INTO users (email, password, name) 
		VALUES ($1, $2, $3)
	`, "demo@studypartner.com", string(hashedPassword), "Demo User")

	if err != nil {
		return fmt.Errorf("failed to create test user: %w", err)
	}

	log.Println("Test user created successfully")
	return nil
}

const createUsersTable = `
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

const createNotesTableWithVector = `
CREATE TABLE IF NOT EXISTS notes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    file_type VARCHAR(50) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    embedding vector(384), -- For all-MiniLM-L6-v2 model
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

const createNotesTableWithoutVector = `
CREATE TABLE IF NOT EXISTS notes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    file_type VARCHAR(50) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

const createSummariesTable = `
CREATE TABLE IF NOT EXISTS summaries (
    id SERIAL PRIMARY KEY,
    note_id INTEGER UNIQUE REFERENCES notes(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

const createFlashcardsTable = `
CREATE TABLE IF NOT EXISTS flashcards (
    id SERIAL PRIMARY KEY,
    note_id INTEGER REFERENCES notes(id) ON DELETE CASCADE,
    question TEXT NOT NULL,
    answer TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

const createQuizzesTable = `
CREATE TABLE IF NOT EXISTS quizzes (
    id SERIAL PRIMARY KEY,
    note_id INTEGER REFERENCES notes(id) ON DELETE CASCADE,
    question TEXT NOT NULL,
    options TEXT[] NOT NULL,
    answer INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

const createStudySessionsTable = `
CREATE TABLE IF NOT EXISTS study_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    note_id INTEGER REFERENCES notes(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    score INTEGER,
    completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

const createIndexesWithVector = `
-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_notes_user_id ON notes(user_id);
CREATE INDEX IF NOT EXISTS idx_notes_embedding ON notes USING ivfflat (embedding vector_cosine_ops);
CREATE INDEX IF NOT EXISTS idx_summaries_note_id ON summaries(note_id);
CREATE INDEX IF NOT EXISTS idx_flashcards_note_id ON flashcards(note_id);
CREATE INDEX IF NOT EXISTS idx_quizzes_note_id ON quizzes(note_id);
CREATE INDEX IF NOT EXISTS idx_study_sessions_user_id ON study_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_study_sessions_note_id ON study_sessions(note_id);
`

const createIndexesWithoutVector = `
-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_notes_user_id ON notes(user_id);
CREATE INDEX IF NOT EXISTS idx_summaries_note_id ON summaries(note_id);
CREATE INDEX IF NOT EXISTS idx_flashcards_note_id ON flashcards(note_id);
CREATE INDEX IF NOT EXISTS idx_quizzes_note_id ON quizzes(note_id);
CREATE INDEX IF NOT EXISTS idx_study_sessions_user_id ON study_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_study_sessions_note_id ON study_sessions(note_id);
`
