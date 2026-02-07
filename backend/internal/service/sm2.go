package service

import (
	"math"
	"time"
)

// SM2Result represents the result of SM-2 algorithm calculation
type SM2Result struct {
	EasinessFactor float64
	Interval       int
	Repetitions    int
	NextReviewAt   time.Time
}

// CalculateSM2 implements the SuperMemo 2 (SM-2) spaced repetition algorithm
// 
// Parameters:
//   - quality: User's recall quality (0-5)
//     0-2: Incorrect response (don't recognize)
//     3: Correct but difficult (vague)
//     4-5: Correct response (recognize)
//   - easinessFactor: Current easiness factor (minimum 1.3)
//   - interval: Current interval in days
//   - repetitions: Number of consecutive correct reviews
//
// Returns:
//   - SM2Result with updated values and next review timestamp
func CalculateSM2(quality int, easinessFactor float64, interval int, repetitions int) SM2Result {
	// Validate quality range
	if quality < 0 {
		quality = 0
	}
	if quality > 5 {
		quality = 5
	}

	// Calculate new easiness factor
	// Formula: EF' = EF + (0.1 - (5-q) * (0.08 + (5-q) * 0.02))
	newEF := easinessFactor + (0.1 - float64(5-quality)*(0.08+float64(5-quality)*0.02))
	
	// Ensure EF doesn't go below 1.3
	if newEF < 1.3 {
		newEF = 1.3
	}

	var newInterval int
	var newRepetitions int

	// If quality < 3, reset the learning process
	if quality < 3 {
		newRepetitions = 0
		newInterval = 1
	} else {
		// Correct response (quality >= 3)
		newRepetitions = repetitions + 1

		switch newRepetitions {
		case 1:
			// First correct review
			newInterval = 1
		case 2:
			// Second correct review
			newInterval = 6
		default:
			// Subsequent reviews: multiply previous interval by EF
			newInterval = int(math.Round(float64(interval) * newEF))
		}
	}

	// Calculate next review timestamp
	nextReviewAt := time.Now().Add(time.Duration(newInterval) * 24 * time.Hour)

	return SM2Result{
		EasinessFactor: newEF,
		Interval:       newInterval,
		Repetitions:    newRepetitions,
		NextReviewAt:   nextReviewAt,
	}
}
