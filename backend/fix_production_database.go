package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	// Get database URL from environment
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Connect to database
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Connected to database successfully")

	// Check if UNIQUE constraint already exists
	var constraintExists bool
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM information_schema.table_constraints 
			WHERE table_name = 'summaries' 
			AND constraint_name = 'summaries_note_id_unique'
		)
	`).Scan(&constraintExists)

	if err != nil {
		log.Fatalf("Failed to check constraint: %v", err)
	}

	if constraintExists {
		fmt.Println("UNIQUE constraint already exists on summaries.note_id")
		return
	}

	// Remove duplicate entries (keep the latest one)
	fmt.Println("Removing duplicate summary entries...")
	result, err := db.Exec(`
		DELETE FROM summaries 
		WHERE id NOT IN (
			SELECT MAX(id) 
			FROM summaries 
			GROUP BY note_id
		)
	`)
	if err != nil {
		log.Fatalf("Failed to remove duplicates: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("Removed %d duplicate summary entries\n", rowsAffected)

	// Add UNIQUE constraint
	fmt.Println("Adding UNIQUE constraint to summaries.note_id...")
	_, err = db.Exec("ALTER TABLE summaries ADD CONSTRAINT summaries_note_id_unique UNIQUE (note_id)")
	if err != nil {
		log.Fatalf("Failed to add UNIQUE constraint: %v", err)
	}

	fmt.Println("Successfully added UNIQUE constraint to summaries.note_id")

	// Add index for better performance
	fmt.Println("Adding index for better performance...")
	_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_summaries_note_id_unique ON summaries(note_id)")
	if err != nil {
		log.Printf("Warning: Failed to add index: %v", err)
	} else {
		fmt.Println("Successfully added index")
	}

	fmt.Println("Database fix completed successfully!")
}
