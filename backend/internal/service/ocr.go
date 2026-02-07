package service

import (
	"context"
	"fmt"

	vision "cloud.google.com/go/vision/v2/apiv1"
	"cloud.google.com/go/vision/v2/apiv1/visionpb"
)

// OCRService handles image text extraction using Cloud Vision API
type OCRService struct {
	client *vision.ImageAnnotatorClient
}

// NewOCRService creates a new OCR service instance
func NewOCRService(ctx context.Context) (*OCRService, error) {
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create vision client: %w", err)
	}
	
	return &OCRService{
		client: client,
	}, nil
}

// Close closes the OCR service client
func (s *OCRService) Close() error {
	return s.client.Close()
}

// ExtractText extracts text from image bytes using DOCUMENT_TEXT_DETECTION
func (s *OCRService) ExtractText(ctx context.Context, imageData []byte) (string, error) {
	image := &visionpb.Image{
		Content: imageData,
	}

	// Use DOCUMENT_TEXT_DETECTION for better text extraction
	feature := &visionpb.Feature{
		Type: visionpb.Feature_DOCUMENT_TEXT_DETECTION,
	}

	request := &visionpb.AnnotateImageRequest{
		Image:    image,
		Features: []*visionpb.Feature{feature},
	}

	response, err := s.client.AnnotateImage(ctx, request)
	if err != nil {
		return "", fmt.Errorf("failed to annotate image: %w", err)
	}

	if response.Error != nil {
		return "", fmt.Errorf("vision api error: %s", response.Error.Message)
	}

	// Extract full text annotation
	if response.FullTextAnnotation != nil {
		return response.FullTextAnnotation.Text, nil
	}

	return "", fmt.Errorf("no text detected in image")
}
