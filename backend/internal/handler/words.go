package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"vocabweb/internal/model"
	"vocabweb/internal/repository"
	"vocabweb/internal/service"

	"github.com/gorilla/mux"
)

type WordsHandler struct {
	wordRepo     *repository.WordRepository
	userWordRepo *repository.UserWordRepository
	analyzer     *service.AnalyzerService
}

func NewWordsHandler(wordRepo *repository.WordRepository, userWordRepo *repository.UserWordRepository, analyzer *service.AnalyzerService) *WordsHandler {
	return &WordsHandler{
		wordRepo:     wordRepo,
		userWordRepo: userWordRepo,
		analyzer:     analyzer,
	}
}

// Helper functions
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

// getUserID extracts user ID from request context (set by auth middleware)
func getUserID(r *http.Request) int64 {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		return 0
	}
	return userID
}

// Create adds a new word manually
// POST /api/v1/words
func (h *WordsHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == 0 {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		Word        string `json:"word"`
		Language    string `json:"language"`
		Definition  string `json:"definition"`
		Source      string `json:"source"`
		Context     string `json:"context"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate required fields
	if req.Word == "" {
		respondError(w, http.StatusBadRequest, "word is required")
		return
	}
	if req.Language == "" {
		req.Language = "en" // Default to English
	}

	ctx := r.Context()

	// Check if word exists in global dictionary
	existingWord, err := h.wordRepo.GetWordByText(ctx, req.Word)
	var wordID int64

	if err != nil || existingWord == nil {
		// Create new word in global dictionary
		newWord := &model.Word{
			Word:       req.Word,
			Language:   req.Language,
			Definition: req.Definition,
		}
		if err := h.wordRepo.CreateWord(ctx, newWord); err != nil {
			respondError(w, http.StatusInternalServerError, "failed to create word")
			return
		}
		wordID = newWord.ID
	} else {
		wordID = existingWord.ID
	}

	// Add to user's collection
	userWord, err := h.userWordRepo.AddUserWord(ctx, userID, wordID, req.Source, req.Context)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to add word to collection")
		return
	}

	respondJSON(w, http.StatusCreated, userWord)
}

// BatchCreate adds multiple words at once
// POST /api/v1/words/batch
func (h *WordsHandler) BatchCreate(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == 0 {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		Words []struct {
			Word       string `json:"word"`
			Language   string `json:"language"`
			Definition string `json:"definition"`
			Source     string `json:"source"`
			Context    string `json:"context"`
		} `json:"words"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if len(req.Words) == 0 {
		respondError(w, http.StatusBadRequest, "words array is required")
		return
	}

	ctx := r.Context()
	var batchWords []struct {
		WordID  int64
		Source  string
		Context string
	}

	for _, wordReq := range req.Words {
		if wordReq.Word == "" {
			continue
		}

		language := wordReq.Language
		if language == "" {
			language = "en"
		}

		// Check if word exists
		existingWord, err := h.wordRepo.GetWordByText(ctx, wordReq.Word)
		var wordID int64

		if err != nil || existingWord == nil {
			// Create new word
			newWord := &model.Word{
				Word:       wordReq.Word,
				Language:   language,
				Definition: wordReq.Definition,
			}
			if err := h.wordRepo.CreateWord(ctx, newWord); err != nil {
				continue // Skip on error
			}
			wordID = newWord.ID
		} else {
			wordID = existingWord.ID
		}

		batchWords = append(batchWords, struct {
			WordID  int64
			Source  string
			Context string
		}{
			WordID:  wordID,
			Source:  wordReq.Source,
			Context: wordReq.Context,
		})
	}

	// Batch add to user collection
	if err := h.userWordRepo.BatchAddUserWords(ctx, userID, batchWords); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to batch add words")
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "words added successfully",
		"count":   len(batchWords),
	})
}

// List returns user's word collection with pagination and filters
// GET /api/v1/words?page=1&limit=20&sort=created_at&order=desc&status=learning
func (h *WordsHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == 0 {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse query parameters
	query := r.URL.Query()
	
	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}
	
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	
	offset := (page - 1) * limit
	
	sort := query.Get("sort")
	if sort == "" {
		sort = "created_at"
	}
	
	order := strings.ToUpper(query.Get("order"))
	if order != "ASC" && order != "DESC" {
		order = "DESC"
	}
	
	status := query.Get("status")

	// Build filters
	filters := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
		"sort":   sort,
		"order":  order,
	}
	
	if status != "" {
		filters["status"] = status
	}

	ctx := r.Context()
	userWords, err := h.userWordRepo.ListUserWords(ctx, userID, filters)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list words")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"words": userWords,
		"page":  page,
		"limit": limit,
	})
}

// Get returns a single word detail
// GET /api/v1/words/:id
func (h *WordsHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == 0 {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	vars := mux.Vars(r)
	wordID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid word id")
		return
	}

	ctx := r.Context()
	userWord, err := h.userWordRepo.GetUserWord(ctx, userID, wordID)
	if err != nil {
		respondError(w, http.StatusNotFound, "word not found")
		return
	}

	respondJSON(w, http.StatusOK, userWord)
}

// Delete removes a word from user's collection
// DELETE /api/v1/words/:id
func (h *WordsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == 0 {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	vars := mux.Vars(r)
	wordID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid word id")
		return
	}

	ctx := r.Context()
	if err := h.userWordRepo.DeleteUserWord(ctx, userID, wordID); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to delete word")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "word deleted successfully",
	})
}

// AnalyzeText analyzes pasted text and returns new word candidates
// POST /api/v1/words/analyze
func (h *WordsHandler) AnalyzeText(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	if userID == 0 {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		Text string `json:"text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Text == "" {
		respondError(w, http.StatusBadRequest, "text is required")
		return
	}

	ctx := r.Context()
	candidates, err := h.analyzer.AnalyzeText(ctx, req.Text, userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to analyze text")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"candidates": candidates,
		"count":      len(candidates),
	})
}
