package service

import (
	"context"
	"fmt"
	"time"

	"vocabweb/internal/repository"

	"github.com/google/uuid"
)

// ReviewService handles review business logic
type ReviewService struct {
	reviewRepo *repository.ReviewRepository
}

// NewReviewService creates a new review service instance
func NewReviewService(reviewRepo *repository.ReviewRepository) *ReviewService {
	return &ReviewService{
		reviewRepo: reviewRepo,
	}
}

// GetDueReviews retrieves words that need review, sorted by forgetting probability
func (s *ReviewService) GetDueReviews(ctx context.Context, userID string, limit int) ([]repository.DueWord, error) {
	if limit <= 0 {
		limit = 20 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit to prevent overload
	}

	dueWords, err := s.reviewRepo.GetDueWords(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get due reviews: %w", err)
	}

	return dueWords, nil
}

// SubmitReviewRequest represents a review submission request
type SubmitReviewRequest struct {
	UserWordID uuid.UUID `json:"user_word_id"`
	Quality    int       `json:"quality"` // 0-5
}

// SubmitReview processes a review submission and updates SM-2 parameters
func (s *ReviewService) SubmitReview(ctx context.Context, userID string, req SubmitReviewRequest) error {
	// Validate quality range
	if req.Quality < 0 || req.Quality > 5 {
		return fmt.Errorf("quality must be between 0 and 5")
	}

	// Get current word state
	dueWords, err := s.reviewRepo.GetDueWords(ctx, userID, 1000)
	if err != nil {
		return fmt.Errorf("failed to get word state: %w", err)
	}

	// Find the word being reviewed
	var currentWord *repository.DueWord
	for i := range dueWords {
		if dueWords[i].UserWordID == req.UserWordID {
			currentWord = &dueWords[i]
			break
		}
	}

	if currentWord == nil {
		return fmt.Errorf("word not found or not due for review")
	}

	// Calculate new SM-2 parameters
	sm2Result := CalculateSM2(
		req.Quality,
		currentWord.EasinessFactor,
		currentWord.Interval,
		currentWord.Repetitions,
	)

	// Create review log
	reviewLog := repository.ReviewLog{
		UserWordID:     req.UserWordID,
		Quality:        req.Quality,
		EasinessFactor: sm2Result.EasinessFactor,
		Interval:       sm2Result.Interval,
		Repetitions:    sm2Result.Repetitions,
		NextReviewAt:   sm2Result.NextReviewAt,
		ReviewedAt:     time.Now(),
	}

	if err := s.reviewRepo.CreateReviewLog(ctx, reviewLog); err != nil {
		return fmt.Errorf("failed to create review log: %w", err)
	}

	// Update user word SM-2 parameters
	if err := s.reviewRepo.UpdateUserWordSM2(
		ctx,
		req.UserWordID,
		sm2Result.EasinessFactor,
		sm2Result.Interval,
		sm2Result.Repetitions,
		sm2Result.NextReviewAt,
	); err != nil {
		return fmt.Errorf("failed to update word SM-2: %w", err)
	}

	return nil
}

// GetReviewStats retrieves review statistics for a user
func (s *ReviewService) GetReviewStats(ctx context.Context, userID string) (*repository.TodayStats, error) {
	stats, err := s.reviewRepo.GetTodayStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get review stats: %w", err)
	}

	return stats, nil
}
