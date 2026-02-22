// Package router sets up all API routes for CareerManifest.
package router

import (
	"github.com/careermanifest/backend/internal/handler"
	"github.com/careermanifest/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

// Setup configures all API routes and returns the Gin engine.
func Setup(
	cfg struct {
		JWTSecret      string
		AllowedOrigins string
		GinMode        string
	},
	authHandler *handler.AuthHandler,
	assessmentHandler *handler.AssessmentHandler,
	questionHandler *handler.QuestionHandler,
	adminHandler *handler.AdminHandler,
) *gin.Engine {
	gin.SetMode(cfg.GinMode)

	r := gin.Default()

	// Global middleware
	r.Use(middleware.CORSMiddleware(cfg.AllowedOrigins))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "CareerManifest API"})
	})

	api := r.Group("/api")
	{
		// ============================================================
		// PUBLIC ROUTES
		// ============================================================
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// ============================================================
		// PROTECTED ROUTES (require JWT)
		// ============================================================
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			// User profile
			protected.GET("/auth/profile", authHandler.Profile)

			// Questions (for assessment form)
			protected.GET("/questions", questionHandler.GetActiveQuestions)

			// Assessments
			protected.POST("/assessment", assessmentHandler.Submit)
			protected.GET("/assessment", assessmentHandler.ListByUser)
			protected.GET("/assessment/:id", assessmentHandler.GetByID)
		}

		// ============================================================
		// ADMIN ROUTES (require JWT + admin role)
		// ============================================================
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		admin.Use(middleware.AdminMiddleware())
		{
			admin.GET("/stats", adminHandler.GetStats)
			admin.GET("/questions", questionHandler.GetAllQuestions)
			admin.POST("/questions", questionHandler.CreateQuestion)
			admin.PUT("/questions/:id", questionHandler.UpdateQuestion)
		}
	}

	return r
}
