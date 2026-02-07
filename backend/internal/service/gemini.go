package service

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/vertexai/genai"
)

// GeminiService handles vocabulary analysis using Vertex AI Gemini
type GeminiService struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

// NewGeminiService creates a new Gemini service instance
func NewGeminiService(ctx context.Context, projectID, location string) (*GeminiService, error) {
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	model := client.GenerativeModel("gemini-1.5-flash")
	model.SetTemperature(0.2) // Lower temperature for more consistent output
	
	return &GeminiService{
		client: client,
		model:  model,
	}, nil
}

// Close closes the Gemini service client
func (s *GeminiService) Close() error {
	return s.client.Close()
}

// VocabWord represents a vocabulary word with its analysis
type VocabWord struct {
	Word            string `json:"word"`
	Definition      string `json:"definition"`
	PartOfSpeech    string `json:"pos"`
	CEFRLevel       string `json:"cefr_level"`
	ContextSentence string `json:"context_sentence"`
}

// buildPrompt creates the analysis prompt for Gemini
func (s *GeminiService) buildPrompt(text, userLevel string) string {
	return fmt.Sprintf(`You are a vocabulary analysis assistant. Analyze the following text and identify vocabulary words that would be challenging for a learner at %s level.

Text to analyze:
%s

Instructions:
1. Identify words above the user's current level
2. For each word, provide: word, definition (concise), part of speech, CEFR level, and a context sentence from the text
3. Return ONLY a valid JSON array, no other text
4. Limit to maximum 20 words
5. Focus on useful vocabulary (skip proper nouns, numbers, basic words)

Output format (JSON array):
[
  {
    "word": "example",
    "definition": "a thing characteristic of its kind",
    "pos": "noun",
    "cefr_level": "B2",
    "context_sentence": "This is an example sentence."
  }
]`, userLevel, text)
}

// AnalyzeVocabulary analyzes text and returns vocabulary words
func (s *GeminiService) AnalyzeVocabulary(ctx context.Context, text, userLevel string) ([]VocabWord, error) {
	if text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}

	// Default to A2 if no level specified
	if userLevel == "" {
		userLevel = "A2"
	}

	prompt := s.buildPrompt(text, userLevel)
	
	resp, err := s.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no response from gemini")
	}

	// Extract text from response
	var responseText string
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			responseText += string(txt)
		}
	}

	// Parse JSON response
	var words []VocabWord
	if err := json.Unmarshal([]byte(responseText), &words); err != nil {
		return nil, fmt.Errorf("failed to parse gemini response: %w (response: %s)", err, responseText)
	}

	return words, nil
}
