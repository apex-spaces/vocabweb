package model

import "time"

// User represents a user in the system
type User struct {
	ID              int64     `json:"id"`
	FirebaseUID     string    `json:"firebase_uid"`
	Email           string    `json:"email"`
	DisplayName     string    `json:"display_name"`
	PhotoURL        string    `json:"photo_url,omitempty"`
	PreferredLocale string    `json:"preferred_locale"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
