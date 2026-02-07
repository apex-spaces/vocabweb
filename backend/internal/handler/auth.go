package handler

import (
	"encoding/json"
	"net/http"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// GetProfile returns the current user's profile (placeholder)
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "Auth profile endpoint - to be implemented",
		"user": map[string]string{
			"id": "placeholder",
			"email": "user@example.com",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
