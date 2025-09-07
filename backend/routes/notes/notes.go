package notes

import (
	"database/sql"
	"encoding/base64"
	"net/http"
	"path/filepath"
	"strings"

	"studypartner/db"
	"studypartner/middleware"
	"studypartner/services"

	"github.com/gin-gonic/gin"
)

type UploadRequest struct {
	File string `json:"file" binding:"required"` // Base64 encoded file
	Name string `json:"name" binding:"required"`
}

func SetupNotesRoutes(router *gin.RouterGroup, database *sql.DB) {
	notes := router.Group("/notes")
	notes.Use(middleware.AuthRequired())
	{
		notes.POST("/upload", uploadNote(database))
		notes.GET("/", getUserNotes(database))
		notes.GET("/:id", getNote(database))
		notes.DELETE("/:id", deleteNote(database))
		notes.POST("/search", searchNotes(database))
	}
}

func uploadNote(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UploadRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, _ := c.Get("userID")

		// Decode base64 file
		fileData, err := base64.StdEncoding.DecodeString(req.File)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file data"})
			return
		}

		// Determine file type
		fileType := strings.ToLower(filepath.Ext(req.Name))
		if fileType == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File type not supported"})
			return
		}

		// Extract text content based on file type
		var content string
		switch fileType {
		case ".txt":
			content = string(fileData)
		case ".pdf":
			content, err = services.ExtractPDFText(fileData)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to extract PDF text"})
				return
			}
		case ".docx":
			content, err = services.ExtractDOCXText(fileData)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to extract DOCX text"})
				return
			}
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported file type"})
			return
		}

		// Generate embedding
		embedding, err := services.GenerateEmbedding(content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate embedding"})
			return
		}

		// Save note to database
		var note db.Note
		err = database.QueryRow(
			`INSERT INTO notes (user_id, title, content, file_type, file_name, file_size, embedding) 
			 VALUES ($1, $2, $3, $4, $5, $6, $7) 
			 RETURNING id, user_id, title, content, file_type, file_name, file_size, created_at, updated_at`,
			userID, req.Name, content, fileType, req.Name, len(fileData), embedding,
		).Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.FileType, &note.FileName, &note.FileSize, &note.CreatedAt, &note.UpdatedAt)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save note"})
			return
		}

		c.JSON(http.StatusCreated, note)
	}
}

func getUserNotes(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		rows, err := database.Query(
			"SELECT id, user_id, title, content, file_type, file_name, file_size, created_at, updated_at FROM notes WHERE user_id = $1 ORDER BY created_at DESC",
			userID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes"})
			return
		}
		defer rows.Close()

		var notes []db.Note
		for rows.Next() {
			var note db.Note
			err := rows.Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.FileType, &note.FileName, &note.FileSize, &note.CreatedAt, &note.UpdatedAt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan note"})
				return
			}
			notes = append(notes, note)
		}

		c.JSON(http.StatusOK, notes)
	}
}

func getNote(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		noteID := c.Param("id")

		var note db.Note
		err := database.QueryRow(
			"SELECT id, user_id, title, content, file_type, file_name, file_size, created_at, updated_at FROM notes WHERE id = $1 AND user_id = $2",
			noteID, userID,
		).Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.FileType, &note.FileName, &note.FileSize, &note.CreatedAt, &note.UpdatedAt)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
			return
		}

		c.JSON(http.StatusOK, note)
	}
}

func deleteNote(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		noteID := c.Param("id")

		result, err := database.Exec("DELETE FROM notes WHERE id = $1 AND user_id = $2", noteID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete note"})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
	}
}

func searchNotes(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Query string `json:"query" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, _ := c.Get("userID")

		// Generate embedding for search query
		queryEmbedding, err := services.GenerateEmbedding(req.Query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate query embedding"})
			return
		}

		// Search using cosine similarity
		rows, err := database.Query(
			`SELECT id, user_id, title, content, file_type, file_name, file_size, created_at, updated_at,
			 1 - (embedding <=> $1) as similarity
			 FROM notes 
			 WHERE user_id = $2 
			 ORDER BY similarity DESC 
			 LIMIT 10`,
			queryEmbedding, userID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search notes"})
			return
		}
		defer rows.Close()

		var results []struct {
			db.Note
			Similarity float64 `json:"similarity"`
		}

		for rows.Next() {
			var result struct {
				db.Note
				Similarity float64 `json:"similarity"`
			}
			err := rows.Scan(&result.ID, &result.UserID, &result.Title, &result.Content, &result.FileType, &result.FileName, &result.FileSize, &result.CreatedAt, &result.UpdatedAt, &result.Similarity)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan search result"})
				return
			}
			results = append(results, result)
		}

		c.JSON(http.StatusOK, results)
	}
}
