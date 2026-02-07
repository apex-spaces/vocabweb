package model

import "time"

// Word represents a word in the dictionary
type Word struct {
	ID          int64     `json:"id"`
	Word        string    `json:"word"`
	Language    string    `json:"language"`
	PartOfSpeech string   `json:"part_of_speech,omitempty"`
	Definition  string    `json:"definition,omitempty"`
	Pronunciation string  `json:"pronunciation,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserWord represents a user's relationship with a word
type UserWord struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	WordID          int64     `json:"word_id"`
	Status          string    `json:"status"` // learning, mastered, reviewing
	Proficiency     int       `json:"proficiency"`
	LastReviewedAt  *time.Time `json:"last_reviewed_at,omitempty"`
	NextReviewAt    *time.Time `json:"next_review_at,omitempty"`
	ReviewCount     int       `json:"review_count"`
	CorrectCount    int       `json:"correct_count"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
