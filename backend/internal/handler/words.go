package handler

import (
	"encoding/json"
	"net/http"
)

type WordsHandler struct{}

func NewWordsHandler() *WordsHandler {
	return &WordsHandler{}
}

// List returns a list of words (placeholder)
func (h *WordsHandler) List(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "List words endpoint - to be implemented",
		"words": []map[string]string{
			{"id": "1", "word": "example", "language": "en"},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Get returns a single word by ID (placeholder)
func (h *WordsHandler) Get(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "Get word endpoint - to be implemented",
		"word": map[string]string{
			"id": "1",
			"word": "example",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
