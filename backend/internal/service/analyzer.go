package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"vocabweb/internal/repository"
)

type AnalyzerService struct {
	wordRepo     *repository.WordRepository
	userWordRepo *repository.UserWordRepository
}

func NewAnalyzerService(wordRepo *repository.WordRepository, userWordRepo *repository.UserWordRepository) *AnalyzerService {
	return &AnalyzerService{
		wordRepo:     wordRepo,
		userWordRepo: userWordRepo,
	}
}

// WordCandidate represents a word candidate from text analysis
type WordCandidate struct {
	Word        string `json:"word"`
	Frequency   int    `json:"frequency"`
	IsCollected bool   `json:"is_collected"`
	WordID      int64  `json:"word_id,omitempty"`
}

// Common English stop words
var stopWords = map[string]bool{
	"a": true, "an": true, "and": true, "are": true, "as": true, "at": true,
	"be": true, "by": true, "for": true, "from": true, "has": true, "he": true,
	"in": true, "is": true, "it": true, "its": true, "of": true, "on": true,
	"that": true, "the": true, "to": true, "was": true, "will": true, "with": true,
	"i": true, "you": true, "we": true, "they": true, "this": true, "but": true,
	"or": true, "not": true, "can": true, "have": true, "do": true, "does": true,
	"did": true, "been": true, "being": true, "am": true, "were": true, "would": true,
	"could": true, "should": true, "may": true, "might": true, "must": true,
}

// AnalyzeText analyzes text and extracts new words for the user
func (s *AnalyzerService) AnalyzeText(ctx context.Context, text string, userID int64) ([]*WordCandidate, error) {
	// Step 1: Tokenize text
	words := s.tokenize(text)
	
	// Step 2: Count word frequency
	wordFreq := make(map[string]int)
	for _, word := range words {
		word = strings.ToLower(word)
		
		// Filter out stop words and short words
		if stopWords[word] || len(word) < 3 {
			continue
		}
		
		wordFreq[word]++
	}
	
	// Step 3: Get user's existing words
	userWords, err := s.userWordRepo.ListUserWords(ctx, userID, map[string]interface{}{
		"limit": 10000, // Get all user words
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user words: %w", err)
	}
	
	// Create a map of user's word IDs
	userWordIDs := make(map[int64]bool)
	for _, uw := range userWords {
		userWordIDs[uw.WordID] = true
	}
	
	// Step 4: Build candidate list
	var candidates []*WordCandidate
	for word, freq := range wordFreq {
		// Check if word exists in global dictionary
		wordObj, err := s.wordRepo.GetWordByText(ctx, word)
		
		candidate := &WordCandidate{
			Word:      word,
			Frequency: freq,
		}
		
		if err == nil && wordObj != nil {
			candidate.WordID = wordObj.ID
			candidate.IsCollected = userWordIDs[wordObj.ID]
		}
		
		candidates = append(candidates, candidate)
	}
	
	// Step 5: Sort by frequency (descending)
	for i := 0; i < len(candidates); i++ {
		for j := i + 1; j < len(candidates); j++ {
			if candidates[i].Frequency < candidates[j].Frequency {
				candidates[i], candidates[j] = candidates[j], candidates[i]
			}
		}
	}
	
	return candidates, nil
}

// tokenize splits text into words
func (s *AnalyzerService) tokenize(text string) []string {
	// Remove punctuation and split by whitespace
	reg := regexp.MustCompile(`[^a-zA-Z\s]+`)
	cleaned := reg.ReplaceAllString(text, " ")
	
	words := strings.Fields(cleaned)
	return words
}
