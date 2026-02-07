package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ReviewRepository handles review-related database operations
type ReviewRepository struct {
	db *DB
}

// NewReviewRepository creates a new review repository instance
func NewReviewRepository(db *DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

// DueWord represents a word that needs review with its details
type DueWord struct {
	UserWordID      uuid.UUID  `json:"user_word_id"`
	WordID          uuid.UUID  `json:"word_id"`
	Word            string     `json:"word"`
	Phonetic        string     `json:"phonetic"`
	Definitions     string     `json:"definitions"` // JSONB as string
	EasinessFactor  float64    `json:"easiness_factor"`
	Interval        int        `json:"interval"`
	Repetitions     int        `json:"repetitions"`
	LastReviewedAt  *time.Time `json:"last_reviewed_at"`
	NextReviewAt    *time.Time `json:"next_review_at"`
	ContextSentence string     `json:"context_sentence"`
}

// GetDueWords retrieves words that are due for review, sorted by forgetting probability
// Higher forgetting probability = more urgent to review
func (r *ReviewRepository) GetDueWords(ctx context.Context, userID string, limit int) ([]DueWord, error) {
	query := `
		SELECT 
			uw.id as user_word_id,
			w.id as word_id,
			w.word,
			w.phonetic,
			w.definitions::text,
			COALESCE(uw.easiness_factor, 2.5) as easiness_factor,
			COALESCE(uw.interval, 0) as interval,
			COALESCE(uw.repetitions, 0) as repetitions,
			uw.last_reviewed_at,
			uw.next_review_at,
			COALESCE(uw.context_sentence, '') as context_sentence
		FROM user_words uw
		JOIN words w ON uw.word_id = w.id
		WHERE uw.user_id = $1
		  AND (uw.next_review_at IS NULL OR uw.next_review_at <= NOW())
		ORDER BY 
			-- Forgetting probability: time elapsed / expected interval
			-- Higher value = more overdue = higher priority
			CASE 
				WHEN uw.last_reviewed_at IS NOT NULL AND uw.interval > 0 THEN
					EXTRACT(EPOCH FROM (NOW() - uw.last_reviewed_at)) / (uw.interval * 86400.0)
				ELSE 999999 -- New words (never reviewed) get highest priority
			END DESC,
			uw.easiness_factor ASC,  -- Harder words first
			uw.repetitions ASC       -- Less practiced words first
		LIMIT $2
	`

	rows, err := r.db.Pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query due words: %w", err)
	}
	defer rows.Close()

	var dueWords []DueWord
	for rows.Next() {
		var word DueWord
		err := rows.Scan(
			&word.UserWordID,
			&word.WordID,
			&word.Word,
			&word.Phonetic,
			&word.Definitions,
			&word.EasinessFactor,
			&word.Interval,
			&word.Repetitions,
			&word.LastReviewedAt,
			&word.NextReviewAt,
			&word.ContextSentence,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan due word: %w", err)
		}
		dueWords = append(dueWords, word)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating due words: %w", err)
	}

	return dueWords, nil
}

// ReviewLog represents a review log entry
type ReviewLog struct {
	ID             uuid.UUID `json:"id"`
	UserWordID     uuid.UUID `json:"user_word_id"`
	Quality        int       `json:"quality"`
	EasinessFactor float64   `json:"easiness_factor"`
	Interval       int       `json:"interval"`
	Repetitions    int       `json:"repetitions"`
	NextReviewAt   time.Time `json:"next_review_at"`
	ReviewedAt     time.Time `json:"reviewed_at"`
}

// CreateReviewLog records a review session in the review_logs table
func (r *ReviewRepository) CreateReviewLog(ctx context.Context, log ReviewLog) error {
	query := `
		INSERT INTO review_logs (
			user_word_id, quality, easiness_factor, interval, 
			repetitions, next_review_at, reviewed_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		log.UserWordID,
		log.Quality,
		log.EasinessFactor,
		log.Interval,
		log.Repetitions,
		log.NextReviewAt,
		log.ReviewedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create review log: %w", err)
	}

	return nil
}

// UpdateUserWordSM2 updates the SM-2 parameters for a user word
func (r *ReviewRepository) UpdateUserWordSM2(
	ctx context.Context,
	userWordID uuid.UUID,
	easinessFactor float64,
	interval int,
	repetitions int,
	nextReviewAt time.Time,
) error {
	query := `
		UPDATE user_words
		SET 
			easiness_factor = $2,
			interval = $3,
			repetitions = $4,
			last_reviewed_at = NOW(),
			next_review_at = $5,
			updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Pool.Exec(ctx, query,
		userWordID,
		easinessFactor,
		interval,
		repetitions,
		nextReviewAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update user word SM-2: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user word not found")
	}

	return nil
}

// TodayStats represents today's review statistics
type TodayStats struct {
	TotalDue      int `json:"total_due"`
	Reviewed      int `json:"reviewed"`
	NewWords      int `json:"new_words"`
	MasteredToday int `json:"mastered_today"`
}

// GetTodayStats retrieves today's review statistics for a user
func (r *ReviewRepository) GetTodayStats(ctx context.Context, userID string) (*TodayStats, error) {
	stats := &TodayStats{}

	// Get total due words count
	dueQuery := `
		SELECT COUNT(*)
		FROM user_words
		WHERE user_id = $1
		  AND (next_review_at IS NULL OR next_review_at <= NOW())
	`
	err := r.db.Pool.QueryRow(ctx, dueQuery, userID).Scan(&stats.TotalDue)
	if err != nil {
		return nil, fmt.Errorf("failed to get due count: %w", err)
	}

	// Get reviewed count today
	reviewedQuery := `
		SELECT COUNT(DISTINCT rl.user_word_id)
		FROM review_logs rl
		JOIN user_words uw ON rl.user_word_id = uw.id
		WHERE uw.user_id = $1
		  AND rl.reviewed_at >= CURRENT_DATE
	`
	err = r.db.Pool.QueryRow(ctx, reviewedQuery, userID).Scan(&stats.Reviewed)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviewed count: %w", err)
	}

	// Get new words added today
	newWordsQuery := `
		SELECT COUNT(*)
		FROM user_words
		WHERE user_id = $1
		  AND collected_at >= CURRENT_DATE
	`
	err = r.db.Pool.QueryRow(ctx, newWordsQuery, userID).Scan(&stats.NewWords)
	if err != nil {
		return nil, fmt.Errorf("failed to get new words count: %w", err)
	}

	// Get words mastered today (repetitions >= 5 and updated today)
	masteredQuery := `
		SELECT COUNT(*)
		FROM user_words
		WHERE user_id = $1
		  AND repetitions >= 5
		  AND is_mastered = true
		  AND updated_at >= CURRENT_DATE
	`
	err = r.db.Pool.QueryRow(ctx, masteredQuery, userID).Scan(&stats.MasteredToday)
	if err != nil {
		return nil, fmt.Errorf("failed to get mastered count: %w", err)
	}

	return stats, nil
}
