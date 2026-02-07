package repository

import (
	"context"
	"fmt"
	"strings"

	"vocabweb/internal/model"
)

type WordRepository struct {
	db *DB
}

func NewWordRepository(db *DB) *WordRepository {
	return &WordRepository{db: db}
}

// CreateWord adds a new word to the global dictionary
func (r *WordRepository) CreateWord(ctx context.Context, word *model.Word) error {
	query := `
		INSERT INTO words (word, language, part_of_speech, definition, pronunciation, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	
	err := r.db.Pool.QueryRow(
		ctx,
		query,
		word.Word,
		word.Language,
		word.PartOfSpeech,
		word.Definition,
		word.Pronunciation,
	).Scan(&word.ID, &word.CreatedAt, &word.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create word: %w", err)
	}
	
	return nil
}

// GetWordByText retrieves a word by its text
func (r *WordRepository) GetWordByText(ctx context.Context, text string) (*model.Word, error) {
	query := `
		SELECT id, word, language, part_of_speech, definition, pronunciation, created_at, updated_at
		FROM words
		WHERE LOWER(word) = LOWER($1)
		LIMIT 1
	`
	
	word := &model.Word{}
	err := r.db.Pool.QueryRow(ctx, query, text).Scan(
		&word.ID,
		&word.Word,
		&word.Language,
		&word.PartOfSpeech,
		&word.Definition,
		&word.Pronunciation,
		&word.CreatedAt,
		&word.UpdatedAt,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get word by text: %w", err)
	}
	
	return word, nil
}

// SearchWords performs fuzzy search on words
func (r *WordRepository) SearchWords(ctx context.Context, query string, limit, offset int) ([]*model.Word, error) {
	sqlQuery := `
		SELECT id, word, language, part_of_speech, definition, pronunciation, created_at, updated_at
		FROM words
		WHERE LOWER(word) LIKE LOWER($1) OR LOWER(definition) LIKE LOWER($1)
		ORDER BY word ASC
		LIMIT $2 OFFSET $3
	`
	
	searchPattern := "%" + strings.TrimSpace(query) + "%"
	
	rows, err := r.db.Pool.Query(ctx, sqlQuery, searchPattern, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search words: %w", err)
	}
	defer rows.Close()
	
	var words []*model.Word
	for rows.Next() {
		word := &model.Word{}
		err := rows.Scan(
			&word.ID,
			&word.Word,
			&word.Language,
			&word.PartOfSpeech,
			&word.Definition,
			&word.Pronunciation,
			&word.CreatedAt,
			&word.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan word: %w", err)
		}
		words = append(words, word)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating words: %w", err)
	}
	
	return words, nil
}
