package db

import (
	"time"

	"github.com/pgvector/pgvector-go"
)

type User struct {
	ID        int       `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"` // Hidden from JSON
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Note struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	Title       string    `json:"title" db:"title"`
	Content     string    `json:"content" db:"content"`
	FileType    string    `json:"file_type" db:"file_type"`
	FileName    string    `json:"file_name" db:"file_name"`
	FileSize    int64     `json:"file_size" db:"file_size"`
	Embedding   pgvector.Vector `json:"-" db:"embedding"` // Hidden from JSON
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Summary struct {
	ID        int       `json:"id" db:"id"`
	NoteID    int       `json:"note_id" db:"note_id"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Flashcard struct {
	ID        int       `json:"id" db:"id"`
	NoteID    int       `json:"note_id" db:"note_id"`
	Question  string    `json:"question" db:"question"`
	Answer    string    `json:"answer" db:"answer"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Quiz struct {
	ID        int       `json:"id" db:"id"`
	NoteID    int       `json:"note_id" db:"note_id"`
	Question  string    `json:"question" db:"question"`
	Options   []string  `json:"options" db:"options"`
	Answer    int       `json:"answer" db:"answer"` // Index of correct option
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type StudySession struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	NoteID    int       `json:"note_id" db:"note_id"`
	Type      string    `json:"type" db:"type"` // "flashcard", "quiz", "summary"
	Score     *int      `json:"score,omitempty" db:"score"`
	Completed bool      `json:"completed" db:"completed"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
