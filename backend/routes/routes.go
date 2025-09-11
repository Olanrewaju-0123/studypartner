package routes

import (
	"database/sql"

	"studypartner/routes/auth"
	"studypartner/routes/notes"
	"studypartner/routes/study"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, db *sql.DB) {
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
