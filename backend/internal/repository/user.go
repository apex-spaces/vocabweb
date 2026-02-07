package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type User struct {
	ID              string    `json:"id"`
	Email           string    `json:"email"`
	DisplayName     string    `json:"display_name"`
	Timezone        string    `json:"timezone"`
	DailyReviewGoal int       `json:"daily_review_goal"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetUserByID retrieves a user by their ID
func (r *UserRepository) GetUserByID(ctx context.Context, userID string) (*User, error) {
	query := `
		SELECT id, email, display_name, timezone, daily_review_goal, created_at, updated_at
		FROM profiles
		WHERE id = $1
	`

	var user User
	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.DisplayName,
		&user.Timezone,
		&user.DailyReviewGoal,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil // User not found
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// CreateUser creates a new user profile
func (r *UserRepository) CreateUser(ctx context.Context, userID, email string) (*User, error) {
	query := `
		INSERT INTO profiles (id, email, display_name, timezone, daily_review_goal)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, email, display_name, timezone, daily_review_goal, created_at, updated_at
	`

	var user User
	err := r.db.Pool.QueryRow(ctx, query, userID, email, "", "UTC", 20).Scan(
		&user.ID,
		&user.Email,
		&user.DisplayName,
		&user.Timezone,
		&user.DailyReviewGoal,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

// UpdateUser updates user profile information
func (r *UserRepository) UpdateUser(ctx context.Context, userID string, displayName, timezone string, dailyReviewGoal int) (*User, error) {
	query := `
		UPDATE profiles
		SET display_name = $2, timezone = $3, daily_review_goal = $4, updated_at = NOW()
		WHERE id = $1
		RETURNING id, email, display_name, timezone, daily_review_goal, created_at, updated_at
	`

	var user User
	err := r.db.Pool.QueryRow(ctx, query, userID, displayName, timezone, dailyReviewGoal).Scan(
		&user.ID,
		&user.Email,
		&user.DisplayName,
		&user.Timezone,
		&user.DailyReviewGoal,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &user, nil
}
