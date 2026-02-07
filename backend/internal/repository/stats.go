package repository

import (
	"context"
	"fmt"
	"time"
)

type StatsRepository struct {
	db *DB
}

func NewStatsRepository(db *DB) *StatsRepository {
	return &StatsRepository{db: db}
}

// GetTodayDueCount returns the number of words due for review today
func (r *StatsRepository) GetTodayDueCount(ctx context.Context, userID int64) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM user_words
		WHERE user_id = $1 
		  AND status = 'learning'
		  AND (next_review_at IS NULL OR next_review_at <= NOW())
	`
	
	var count int
	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get today due count: %w", err)
	}
	
	return count, nil
}

// GetTodayNewCount returns the number of words added today
func (r *StatsRepository) GetTodayNewCount(ctx context.Context, userID int64) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM user_words
		WHERE user_id = $1 
		  AND DATE(created_at) = CURRENT_DATE
	`
	
	var count int
	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get today new count: %w", err)
	}
	
	return count, nil
}

// GetMasteredCount returns the number of mastered words (review_count >= 5)
func (r *StatsRepository) GetMasteredCount(ctx context.Context, userID int64) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM user_words
		WHERE user_id = $1 
		  AND review_count >= 5
	`
	
	var count int
	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get mastered count: %w", err)
	}
	
	return count, nil
}

// GetStreakDays returns the number of consecutive days the user has studied
func (r *StatsRepository) GetStreakDays(ctx context.Context, userID int64) (int, error) {
	query := `
		WITH daily_activity AS (
			SELECT DISTINCT DATE(last_reviewed_at) as activity_date
			FROM user_words
			WHERE user_id = $1 
			  AND last_reviewed_at IS NOT NULL
			ORDER BY activity_date DESC
		)
		SELECT COUNT(*) as streak
		FROM (
			SELECT activity_date,
			       activity_date - (ROW_NUMBER() OVER (ORDER BY activity_date DESC))::int as grp
			FROM daily_activity
		) grouped
		WHERE grp = (
			SELECT activity_date - (ROW_NUMBER() OVER (ORDER BY activity_date DESC))::int
			FROM daily_activity
			LIMIT 1
		)
	`
	
	var streak int
	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(&streak)
	if err != nil {
		// If no activity found, return 0 instead of error
		return 0, nil
	}
	
	return streak, nil
}

// RecentWord represents a recently collected word with its details
type RecentWord struct {
	WordID      int64     `json:"word_id"`
	Word        string    `json:"word"`
	Definition  string    `json:"definition"`
	CreatedAt   time.Time `json:"created_at"`
}

// GetRecentWords returns the most recently collected words
func (r *StatsRepository) GetRecentWords(ctx context.Context, userID int64, limit int) ([]RecentWord, error) {
	query := `
		SELECT uw.word_id, w.word, w.definition, uw.created_at
		FROM user_words uw
		JOIN words w ON uw.word_id = w.id
		WHERE uw.user_id = $1
		ORDER BY uw.created_at DESC
		LIMIT $2
	`
	
	rows, err := r.db.Pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent words: %w", err)
	}
	defer rows.Close()
	
	var words []RecentWord
	for rows.Next() {
		var word RecentWord
		err := rows.Scan(&word.WordID, &word.Word, &word.Definition, &word.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan recent word: %w", err)
		}
		words = append(words, word)
	}
	
	return words, nil
}

// DailyStat represents daily learning statistics
type DailyStat struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// GetWeeklyStats returns the last 7 days of learning statistics
func (r *StatsRepository) GetWeeklyStats(ctx context.Context, userID int64) ([]DailyStat, error) {
	query := `
		WITH date_series AS (
			SELECT generate_series(
				CURRENT_DATE - INTERVAL '6 days',
				CURRENT_DATE,
				'1 day'::interval
			)::date as date
		)
		SELECT 
			TO_CHAR(ds.date, 'Mon') as day_name,
			COALESCE(COUNT(uw.id), 0) as count
		FROM date_series ds
		LEFT JOIN user_words uw ON DATE(uw.last_reviewed_at) = ds.date AND uw.user_id = $1
		GROUP BY ds.date, day_name
		ORDER BY ds.date
	`
	
	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get weekly stats: %w", err)
	}
	defer rows.Close()
	
	var stats []DailyStat
	for rows.Next() {
		var stat DailyStat
		err := rows.Scan(&stat.Date, &stat.Count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan daily stat: %w", err)
		}
		stats = append(stats, stat)
	}
	
	return stats, nil
}
