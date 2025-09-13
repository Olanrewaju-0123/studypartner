package study

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"studypartner/db"
	"studypartner/middleware"
	"studypartner/services"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func SetupStudyRoutes(router *gin.RouterGroup, database *sql.DB) {
	study := router.Group("/study")
	study.Use(middleware.AuthRequired())
	{
		study.GET("/notes/:id/summary", getSummary(database))
		study.POST("/notes/:id/summary", generateSummary(database))
		study.GET("/notes/:id/flashcards", getFlashcards(database))
		study.POST("/notes/:id/flashcards", generateFlashcards(database))
		study.GET("/notes/:id/quiz", getQuiz(database))
		study.POST("/notes/:id/quiz", generateQuiz(database))
		study.POST("/sessions", createStudySession(database))
		study.PUT("/sessions/:id", updateStudySession(database))
	}
}

// GetSummary godoc
// @Summary Get note summary
// @Description Get the AI-generated summary for a specific note
// @Tags Study Materials
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Note ID"
// @Success 200 {object} db.Summary "Note summary"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Note or summary not found"
// @Router /study/notes/{id}/summary [get]
func getSummary(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		noteID := c.Param("id")

		// Check if note belongs to user
		var note db.Note
		err := database.QueryRow(
			"SELECT id FROM notes WHERE id = $1 AND user_id = $2",
			noteID, userID,
		).Scan(&note.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
			return
		}

		// Get existing summary
		var summary db.Summary
		err = database.QueryRow(
			"SELECT id, note_id, content, created_at, updated_at FROM summaries WHERE note_id = $1",
			noteID,
		).Scan(&summary.ID, &summary.NoteID, &summary.Content, &summary.CreatedAt, &summary.UpdatedAt)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Summary not found"})
			return
		}

		c.JSON(http.StatusOK, summary)
	}
}

func generateSummary(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		noteID := c.Param("id")

		// Check if note belongs to user and get content
		var note db.Note
		err := database.QueryRow(
			"SELECT id, content FROM notes WHERE id = $1 AND user_id = $2",
			noteID, userID,
		).Scan(&note.ID, &note.Content)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
			return
		}

		// Generate summary using AI
		summaryContent, err := services.GenerateSummary(note.Content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate summary"})
			return
		}

		// Save or update summary
		var summary db.Summary
		err = database.QueryRow(
			`INSERT INTO summaries (note_id, content) VALUES ($1, $2) 
			 ON CONFLICT (note_id) DO UPDATE SET content = $2, updated_at = CURRENT_TIMESTAMP
			 RETURNING id, note_id, content, created_at, updated_at`,
			noteID, summaryContent,
		).Scan(&summary.ID, &summary.NoteID, &summary.Content, &summary.CreatedAt, &summary.UpdatedAt)

		if err != nil {
			// Log the actual error for debugging
			fmt.Printf("Failed to save summary for note %s: %v\n", noteID, err)
			
			// Check if it's a constraint error
			if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database constraint error - please contact support"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save summary"})
			}
			return
		}

		c.JSON(http.StatusOK, summary)
	}
}

func getFlashcards(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		noteID := c.Param("id")

		// Check if note belongs to user
		var note db.Note
		err := database.QueryRow(
			"SELECT id FROM notes WHERE id = $1 AND user_id = $2",
			noteID, userID,
		).Scan(&note.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
			return
		}

		// Get existing flashcards
		rows, err := database.Query(
			"SELECT id, note_id, question, answer, created_at FROM flashcards WHERE note_id = $1 ORDER BY created_at",
			noteID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch flashcards"})
			return
		}
		defer rows.Close()

		var flashcards []db.Flashcard
		for rows.Next() {
			var flashcard db.Flashcard
			err := rows.Scan(&flashcard.ID, &flashcard.NoteID, &flashcard.Question, &flashcard.Answer, &flashcard.CreatedAt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan flashcard"})
				return
			}
			flashcards = append(flashcards, flashcard)
		}

		c.JSON(http.StatusOK, flashcards)
	}
}

func generateFlashcards(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		noteID := c.Param("id")

		// Check if note belongs to user and get content
		var note db.Note
		err := database.QueryRow(
			"SELECT id, content FROM notes WHERE id = $1 AND user_id = $2",
			noteID, userID,
		).Scan(&note.ID, &note.Content)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
			return
		}

		// Generate flashcards using AI
		flashcards, err := services.GenerateFlashcards(note.Content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate flashcards"})
			return
		}

		// Clear existing flashcards for this note
		_, err = database.Exec("DELETE FROM flashcards WHERE note_id = $1", noteID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear existing flashcards"})
			return
		}

		// Insert new flashcards
		var insertedFlashcards []db.Flashcard
		for _, fc := range flashcards {
			var flashcard db.Flashcard
			err = database.QueryRow(
				"INSERT INTO flashcards (note_id, question, answer) VALUES ($1, $2, $3) RETURNING id, note_id, question, answer, created_at",
				noteID, fc.Question, fc.Answer,
			).Scan(&flashcard.ID, &flashcard.NoteID, &flashcard.Question, &flashcard.Answer, &flashcard.CreatedAt)
			if err != nil {
				// Log the actual error for debugging
				fmt.Printf("Failed to save flashcard for note %s: %v\n", noteID, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save flashcard"})
				return
			}
			insertedFlashcards = append(insertedFlashcards, flashcard)
		}

		c.JSON(http.StatusOK, insertedFlashcards)
	}
}

func getQuiz(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		noteID := c.Param("id")

		// Check if note belongs to user
		var note db.Note
		err := database.QueryRow(
			"SELECT id FROM notes WHERE id = $1 AND user_id = $2",
			noteID, userID,
		).Scan(&note.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
			return
		}

		// Get existing quiz questions
		rows, err := database.Query(
			"SELECT id, note_id, question, options, answer, created_at FROM quizzes WHERE note_id = $1 ORDER BY created_at",
			noteID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch quiz"})
			return
		}
		defer rows.Close()

		var quizQuestions []db.Quiz
		for rows.Next() {
			var quiz db.Quiz
			err := rows.Scan(&quiz.ID, &quiz.NoteID, &quiz.Question, pq.Array(&quiz.Options), &quiz.Answer, &quiz.CreatedAt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan quiz question"})
				return
			}
			// Debug logging
			fmt.Printf("Retrieved quiz question: ID=%d, Question=%s, Options=%v, Answer=%d\n", 
				quiz.ID, quiz.Question, quiz.Options, quiz.Answer)
			quizQuestions = append(quizQuestions, quiz)
		}

		fmt.Printf("Returning %d quiz questions\n", len(quizQuestions))
		c.JSON(http.StatusOK, quizQuestions)
	}
}

func generateQuiz(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		noteID := c.Param("id")

		// Check if note belongs to user and get content
		var note db.Note
		err := database.QueryRow(
			"SELECT id, content FROM notes WHERE id = $1 AND user_id = $2",
			noteID, userID,
		).Scan(&note.ID, &note.Content)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
			return
		}

		// Generate quiz using AI
		quizQuestions, err := services.GenerateQuiz(note.Content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate quiz"})
			return
		}

		// Clear existing quiz for this note
		_, err = database.Exec("DELETE FROM quizzes WHERE note_id = $1", noteID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear existing quiz"})
			return
		}

		// Insert new quiz questions
		var insertedQuiz []db.Quiz
		for i, q := range quizQuestions {
			fmt.Printf("Saving quiz question %d: Question=%s, Options=%v, Answer=%d\n", 
				i+1, q.Question, q.Options, q.Answer)
			
			var quiz db.Quiz
			err = database.QueryRow(
				"INSERT INTO quizzes (note_id, question, options, answer) VALUES ($1, $2, $3, $4) RETURNING id, note_id, question, options, answer, created_at",
				noteID, q.Question, pq.Array(q.Options), q.Answer,
			).Scan(&quiz.ID, &quiz.NoteID, &quiz.Question, pq.Array(&quiz.Options), &quiz.Answer, &quiz.CreatedAt)
			if err != nil {
				// Log the actual error for debugging
				fmt.Printf("Failed to save quiz question for note %s: %v\n", noteID, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save quiz question"})
				return
			}
			
			fmt.Printf("Successfully saved quiz question: ID=%d, Options=%v\n", quiz.ID, quiz.Options)
			insertedQuiz = append(insertedQuiz, quiz)
		}

		c.JSON(http.StatusOK, insertedQuiz)
	}
}

func createStudySession(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			NoteID int    `json:"note_id" binding:"required"`
			Type   string `json:"type" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, _ := c.Get("userID")

		// Check if note belongs to user
		var note db.Note
		err := database.QueryRow(
			"SELECT id FROM notes WHERE id = $1 AND user_id = $2",
			req.NoteID, userID,
		).Scan(&note.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
			return
		}

		// Create study session
		var session db.StudySession
		err = database.QueryRow(
			"INSERT INTO study_sessions (user_id, note_id, type) VALUES ($1, $2, $3) RETURNING id, user_id, note_id, type, score, completed, created_at",
			userID, req.NoteID, req.Type,
		).Scan(&session.ID, &session.UserID, &session.NoteID, &session.Type, &session.Score, &session.Completed, &session.CreatedAt)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create study session"})
			return
		}

		c.JSON(http.StatusCreated, session)
	}
}

func updateStudySession(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.Param("id")
		userID, _ := c.Get("userID")

		var req struct {
			Score     *int  `json:"score"`
			Completed bool  `json:"completed"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update study session
		var session db.StudySession
		err := database.QueryRow(
			"UPDATE study_sessions SET score = $1, completed = $2 WHERE id = $3 AND user_id = $4 RETURNING id, user_id, note_id, type, score, completed, created_at",
			req.Score, req.Completed, sessionID, userID,
		).Scan(&session.ID, &session.UserID, &session.NoteID, &session.Type, &session.Score, &session.Completed, &session.CreatedAt)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Study session not found"})
			return
		}

		c.JSON(http.StatusOK, session)
	}
}
