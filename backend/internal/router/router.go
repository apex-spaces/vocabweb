package router

import (
	"github.com/apex-spaces/vocabweb/backend/internal/handler"
	"github.com/apex-spaces/vocabweb/backend/internal/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	healthHandler *handler.HealthHandler
	authHandler   *handler.AuthHandler
	wordsHandler  *handler.WordsHandler
	authMiddleware *middleware.AuthMiddleware
}

func New(
	healthHandler *handler.HealthHandler,
	authHandler *handler.AuthHandler,
	wordsHandler *handler.WordsHandler,
	authMiddleware *middleware.AuthMiddleware,
) *Router {
	return &Router{
		healthHandler: healthHandler,
		authHandler: authHandler,
		wordsHandler: wordsHandler,
		authMiddleware: authMiddleware,
	}
}

func (rt *Router) Setup() *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)

	// Health check (public)
	r.Get("/health", rt.healthHandler.Check)

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Get("/health", rt.healthHandler.Check)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(rt.authMiddleware.Verify)

			// Auth
			r.Get("/auth/profile", rt.authHandler.GetProfile)

			// Words
			r.Get("/words", rt.wordsHandler.List)
			r.Get("/words/{id}", rt.wordsHandler.Get)
		})
	})

	return r
}
