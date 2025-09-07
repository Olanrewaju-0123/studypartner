package routes

import (
	"database/sql"

	"studypartner/routes/auth"
	"studypartner/routes/notes"
	"studypartner/routes/study"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, db *sql.DB) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := router.Group("/api")
	{
		// Auth routes
		auth.SetupAuthRoutes(api, db)
		
		// Notes routes
		notes.SetupNotesRoutes(api, db)
		
		// Study routes
		study.SetupStudyRoutes(api, db)
	}
}
