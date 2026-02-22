// CareerManifest Backend - Main Entry Point
// AI-Powered Career Decision Platform for Indian Students
package main

import (
	"log"

	"github.com/careermanifest/backend/internal/config"
	"github.com/careermanifest/backend/internal/database"
	"github.com/careermanifest/backend/internal/engine"
	"github.com/careermanifest/backend/internal/handler"
	"github.com/careermanifest/backend/internal/repository"
	"github.com/careermanifest/backend/internal/router"
	"github.com/careermanifest/backend/internal/seed"
	"github.com/careermanifest/backend/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	// Connect to PostgreSQL database (Neon)
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run database migrations
	if err := database.Migrate(db); err != nil {
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	assessmentRepo := repository.NewAssessmentRepository(db)
	questionRepo := repository.NewQuestionRepository(db)

	// Initialize scoring engine
	scoringEngine := engine.NewScoringEngine()

	// Initialize LLM service (works without API key)
	llmService := service.NewLLMService(cfg)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg)
	assessmentService := service.NewAssessmentService(assessmentRepo, questionRepo, scoringEngine, llmService)
	questionService := service.NewQuestionService(questionRepo)
	adminService := service.NewAdminService(userRepo, assessmentRepo, questionRepo)

	// Seed default admin user
	if err := authService.SeedAdmin(); err != nil {
		log.Printf("⚠️ Failed to seed admin: %v", err)
	}

	// Seed assessment questions
	if err := seed.SeedQuestions(questionRepo); err != nil {
		log.Printf("⚠️ Failed to seed questions: %v", err)
	}

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	assessmentHandler := handler.NewAssessmentHandler(assessmentService)
	questionHandler := handler.NewQuestionHandler(questionService)
	adminHandler := handler.NewAdminHandler(adminService)

	// Setup router
	r := router.Setup(
		struct {
			JWTSecret      string
			AllowedOrigins string
			GinMode        string
		}{
			JWTSecret:      cfg.JWTSecret,
			AllowedOrigins: cfg.AllowedOrigins,
			GinMode:        cfg.GinMode,
		},
		authHandler,
		assessmentHandler,
		questionHandler,
		adminHandler,
	)

	// Start server
	log.Printf("🚀 CareerManifest API starting on port %s", cfg.Port)
	log.Printf("📊 LLM Integration: %v", cfg.IsLLMEnabled())
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
