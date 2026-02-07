package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"vocabweb/internal/service"

	"github.com/google/uuid"
)

// ReviewHandler handles review-related HTTP requests
type ReviewHandler struct {
	reviewService *service.ReviewService
}

// NewReviewHandler creates a new review handler instance
func NewReviewHandler(reviewService *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{
		reviewService: reviewService,
	}
}

// GetDueReviews handles GET /api/v1/review/due
// Returns list of words due for review
func (h *ReviewHandler) GetDueReviews(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Get user ID from context (set by auth middleware)
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse limit from query params
	limitStr := r.URL.Query().Get("limit")
	limit := 20 // default
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	// Get due reviews
	dueWords, err := h.reviewService.GetDueReviews(ctx, userID, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"words": dueWords,
			"total": len(dueWords),
		},
	})
}

// SubmitReview handles POST /api/v1/review/submit
// Submits a review result and updates SM-2 parameters
func (h *ReviewHandler) SubmitReview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Get user ID from context
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var req struct {
		UserWordID string `json:"user_word_id"`
		Quality    int    `json:"quality"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Parse UUID
	userWordID, err := uuid.Parse(req.UserWordID)
	if err != nil {
		http.Error(w, "Invalid user_word_id", http.StatusBadRequest)
		return
	}

	// Submit review
	submitReq := service.SubmitReviewRequest{
		UserWordID: userWordID,
		Quality:    req.Quality,
	}

	if err := h.reviewService.SubmitReview(ctx, userID, submitReq); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Review submitted successfully",
	})
}

// GetReviewStats handles GET /api/v1/review/stats
// Returns review statistics for the user
func (h *ReviewHandler) GetReviewStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Get user ID from context
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get review stats
	stats, err := h.reviewService.GetReviewStats(ctx, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    stats,
	})
}
