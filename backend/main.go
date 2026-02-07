package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/apex-spaces/vocabweb/backend/internal/config"
	"github.com/apex-spaces/vocabweb/backend/internal/handler"
	"github.com/apex-spaces/vocabweb/backend/internal/middleware"
	"github.com/apex-spaces/vocabweb/backend/internal/repository"
	"github.com/apex-spaces/vocabweb/backend/internal/router"
)

func main() {
	ctx := context.Background()

	// Load configuration
	cfg := config.Load()
	log.Printf("Starting VocabWeb backend on port %s (env: %s)", cfg.Port, cfg.Env)

	// Initialize database (optional, skip if DATABASE_URL is empty)
	var db *repository.DB
	if cfg.DatabaseURL != "" {
		var err error
		db, err = repository.NewDB(ctx, cfg.DatabaseURL)
		if err != nil {
			log.Printf("Warning: Failed to connect to database: %v", err)
		} else {
			defer db.Close()
		}
	}

	// Initialize auth middleware
	authMiddleware, err := middleware.NewAuthMiddleware(ctx, cfg.FirebaseProjectID)
	if err != nil {
		log.Printf("Warning: Failed to initialize auth middleware: %v", err)
	}

	// Initialize repositories
	var statsRepo *repository.StatsRepository
	if db != nil {
		statsRepo = repository.NewStatsRepository(db)
	}

	// Initialize handlers
	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler()
	wordsHandler := handler.NewWordsHandler()
	dashboardHandler := handler.NewDashboardHandler(statsRepo)

	// Setup router
	rt := router.New(healthHandler, authHandler, wordsHandler, dashboardHandler, authMiddleware)
	r := rt.Setup()

	// Apply CORS middleware
	corsMiddleware := middleware.CORS(cfg.AllowedOrigins)
	handler := corsMiddleware(r)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
