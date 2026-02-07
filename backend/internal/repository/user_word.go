package repository

import (
	"context"
	"fmt"
	"strings"

	"vocabweb/internal/model"
)

type UserWordRepository struct {
	db *DB
}

func NewUserWordRepository(db *DB) *UserWordRepository {
	return &UserWordRepository{db: db}
}

// AddUserWord adds a word to user's collection
func (r *UserWordRepository) AddUserWord(ctx context.Context, userID, wordID int64, source, contextText string) (*model.UserWord, error) {
	query := `
		INSERT INTO user_words (user_id, word_id, status, proficiency, review_count, correct_count, source, context, created_at, updated_at)
		VALUES ($1, $2, 'learning', 0, 0, 0, $3, $4, NOW(), NOW())
		RETURNING id, user_id, word_id, status, proficiency, last_reviewed_at, next_review_at, review_count, correct_count, created_at, updated_at
	`
	
	userWord := &model.UserWord{}
	err := r.db.Pool.QueryRow(ctx, query, userID, wordID, source, contextText).Scan(
		&userWord.ID,
		&userWord.UserID,
		&userWord.WordID,
		&userWord.Status,
		&userWord.Proficiency,
		&userWord.LastReviewedAt,
		&userWord.NextReviewAt,
		&userWord.ReviewCount,
		&userWord.CorrectCount,
		&userWord.CreatedAt,
		&userWord.UpdatedAt,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to add user word: %w", err)
	}
	
	return userWord, nil
}

// ListUserWords retrieves user's words with filters
func (r *UserWordRepository) ListUserWords(ctx context.Context, userID int64, filters map[string]interface{}) ([]*model.UserWord, error) {
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(`
		SELECT uw.id, uw.user_id, uw.word_id, uw.status, uw.proficiency, 
		       uw.last_reviewed_at, uw.next_review_at, uw.review_count, uw.correct_count,
		       uw.created_at, uw.updated_at
		FROM user_words uw
		WHERE uw.user_id = $1
	`)
	
	args := []interface{}{userID}
	argCount := 1
	
	// Apply filters
	if status, ok := filters["status"].(string); ok && status != "" {
		argCount++
		queryBuilder.WriteString(fmt.Sprintf(" AND uw.status = $%d", argCount))
		args = append(args, status)
	}
	
	// Sorting
	sortBy := "created_at"
	if sort, ok := filters["sort"].(string); ok && sort != "" {
		sortBy = sort
	}
	
	order := "DESC"
	if ord, ok := filters["order"].(string); ok && strings.ToUpper(ord) == "ASC" {
		order = "ASC"
	}
	
	queryBuilder.WriteString(fmt.Sprintf(" ORDER BY uw.%s %s", sortBy, order))
	
	// Pagination
	limit := 20
	if l, ok := filters["limit"].(int); ok && l > 0 {
		limit = l
	}
	
	offset := 0
	if o, ok := filters["offset"].(int); ok && o >= 0 {
		offset = o
	}
	
	argCount++
	queryBuilder.WriteString(fmt.Sprintf(" LIMIT $%d", argCount))
	args = append(args, limit)
	
	argCount++
	queryBuilder.WriteString(fmt.Sprintf(" OFFSET $%d", argCount))
	args = append(args, offset)
	
	rows, err := r.db.Pool.Query(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list user words: %w", err)
	}
	defer rows.Close()
	
	var userWords []*model.UserWord
	for rows.Next() {
		uw := &model.UserWord{}
		err := rows.Scan(
			&uw.ID,
			&uw.UserID,
			&uw.WordID,
			&uw.Status,
			&uw.Proficiency,
			&uw.LastReviewedAt,
			&uw.NextReviewAt,
			&uw.ReviewCount,
			&uw.CorrectCount,
			&uw.CreatedAt,
			&uw.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user word: %w", err)
		}
		userWords = append(userWords, uw)
	}
	
	return userWords, nil
}

// GetUserWord retrieves a single user word
func (r *UserWordRepository) GetUserWord(ctx context.Context, userID, wordID int64) (*model.UserWord, error) {
	query := `
		SELECT id, user_id, word_id, status, proficiency, last_reviewed_at, next_review_at, 
		       review_count, correct_count, created_at, updated_at
		FROM user_words
		WHERE user_id = $1 AND word_id = $2
		LIMIT 1
	`
	
	userWord := &model.UserWord{}
	err := r.db.Pool.QueryRow(ctx, query, userID, wordID).Scan(
		&userWord.ID,
		&userWord.UserID,
		&userWord.WordID,
		&userWord.Status,
		&userWord.Proficiency,
		&userWord.LastReviewedAt,
		&userWord.NextReviewAt,
		&userWord.ReviewCount,
		&userWord.CorrectCount,
		&userWord.CreatedAt,
		&userWord.UpdatedAt,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get user word: %w", err)
	}
	
	return userWord, nil
}

// DeleteUserWord removes a word from user's collection
func (r *UserWordRepository) DeleteUserWord(ctx context.Context, userID, wordID int64) error {
	query := `DELETE FROM user_words WHERE user_id = $1 AND word_id = $2`
	
	result, err := r.db.Pool.Exec(ctx, query, userID, wordID)
	if err != nil {
		return fmt.Errorf("failed to delete user word: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user word not found")
	}
	
	return nil
}

// BatchAddUserWords adds multiple words to user's collection
func (r *UserWordRepository) BatchAddUserWords(ctx context.Context, userID int64, words []struct {
	WordID  int64
	Source  string
	Context string
}) error {
	if len(words) == 0 {
		return nil
	}
	
	// Use transaction for batch insert
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	
	query := `
		INSERT INTO user_words (user_id, word_id, status, proficiency, review_count, correct_count, source, context, created_at, updated_at)
		VALUES ($1, $2, 'learning', 0, 0, 0, $3, $4, NOW(), NOW())
		ON CONFLICT (user_id, word_id) DO NOTHING
	`
	
	for _, word := range words {
		_, err := tx.Exec(ctx, query, userID, word.WordID, word.Source, word.Context)
		if err != nil {
			return fmt.Errorf("failed to insert word %d: %w", word.WordID, err)
		}
	}
	
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	return nil
}
