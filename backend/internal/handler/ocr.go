package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/apex-spaces/vocabweb/backend/internal/service"
)

const (
	maxUploadSize = 10 << 20 // 10 MB
)

// OCRHandler handles OCR-related requests
type OCRHandler struct {
	ocrService    *service.OCRService
	geminiService *service.GeminiService
}

// NewOCRHandler creates a new OCR handler
func NewOCRHandler(ocrService *service.OCRService, geminiService *service.GeminiService) *OCRHandler {
	return &OCRHandler{
		ocrService:    ocrService,
		geminiService: geminiService,
	}
}

// OCRAnalyzeResponse represents the response for OCR analysis
type OCRAnalyzeResponse struct {
	ExtractedText string               `json:"extracted_text"`
	Words         []service.VocabWord  `json:"words"`
}

// AnalyzeImage handles POST /api/v1/ocr/analyze
func (h *OCRHandler) AnalyzeImage(w http.ResponseWriter, r *http.Request) {
	// Limit request body size
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	// Parse multipart form
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(w, "File too large or invalid form data", http.StatusBadRequest)
		return
	}

	// Get uploaded file
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Missing or invalid image file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read image data into memory
	imageData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read image data", http.StatusInternalServerError)
		return
	}

	// Get optional user level parameter
	userLevel := r.FormValue("level")
	if userLevel == "" {
		userLevel = "A2" // Default level
	}

	// Step 1: Extract text using Cloud Vision OCR
	extractedText, err := h.ocrService.ExtractText(r.Context(), imageData)
	if err != nil {
		http.Error(w, fmt.Sprintf("OCR failed: %v", err), http.StatusInternalServerError)
		return
	}

	if extractedText == "" {
		http.Error(w, "No text detected in image", http.StatusBadRequest)
		return
	}

	// Step 2: Analyze vocabulary using Gemini
	words, err := h.geminiService.AnalyzeVocabulary(r.Context(), extractedText, userLevel)
	if err != nil {
		http.Error(w, fmt.Sprintf("Vocabulary analysis failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Step 3: Return response
	response := OCRAnalyzeResponse{
		ExtractedText: extractedText,
		Words:         words,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
