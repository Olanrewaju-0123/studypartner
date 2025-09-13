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
		databaseURL = "postgres://user:password@localhost/studypartner?sslmode=disable"
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

	// Check if summaries table exists
	var tableExists bool
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM information_schema.tables 
			WHERE table_name = 'summaries'
		)
	`).Scan(&tableExists)

	if err != nil {
		log.Fatalf("Failed to check if summaries table exists: %v", err)
	}

	if !tableExists {
		fmt.Println("❌ Summaries table does not exist")
		return
	}

	fmt.Println("✅ Summaries table exists")

	// Check if UNIQUE constraint exists
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
		fmt.Println("✅ UNIQUE constraint exists on summaries.note_id")
	} else {
		fmt.Println("❌ UNIQUE constraint missing on summaries.note_id")
	}

	// Test the ON CONFLICT query
	fmt.Println("\nTesting ON CONFLICT query...")
	
	// First, let's see if there are any notes
	var noteCount int
	err = db.QueryRow("SELECT COUNT(*) FROM notes").Scan(&noteCount)
	if err != nil {
		log.Printf("Warning: Could not count notes: %v", err)
		noteCount = 0
	}

	fmt.Printf("Found %d notes in database\n", noteCount)

	if noteCount > 0 {
		// Get the first note ID
		var noteID int
		err = db.QueryRow("SELECT id FROM notes LIMIT 1").Scan(&noteID)
		if err != nil {
			log.Printf("Warning: Could not get note ID: %v", err)
		} else {
			fmt.Printf("Testing with note ID: %d\n", noteID)
			
			// Test the INSERT with ON CONFLICT
			var summaryID int
			err = db.QueryRow(`
				INSERT INTO summaries (note_id, content) VALUES ($1, $2) 
				ON CONFLICT (note_id) DO UPDATE SET content = $2, updated_at = CURRENT_TIMESTAMP
				RETURNING id`,
				noteID, "Test summary content",
			).Scan(&summaryID)

			if err != nil {
				fmt.Printf("❌ ON CONFLICT query failed: %v\n", err)
			} else {
				fmt.Printf("✅ ON CONFLICT query succeeded, summary ID: %d\n", summaryID)
			}
		}
	}

	fmt.Println("\nTest completed!")
}
