package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"vocabweb/internal/middleware"
	"vocabweb/internal/repository"
)

type DashboardHandler struct {
	statsRepo *repository.StatsRepository
}

func NewDashboardHandler(statsRepo *repository.StatsRepository) *DashboardHandler {
	return &DashboardHandler{statsRepo: statsRepo}
}

// DashboardResponse represents the complete dashboard data
type DashboardResponse struct {
	TodayDue      int                        `json:"today_due"`
	TodayNew      int                        `json:"today_new"`
	TotalMastered int                        `json:"total_mastered"`
	StreakDays    int                        `json:"streak_days"`
	RecentWords   []repository.RecentWord    `json:"recent_words"`
	WeeklyStats   []repository.DailyStat     `json:"weekly_stats"`
}

// GetDashboard returns all dashboard statistics for the current user
func (h *DashboardHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Fetch all statistics
	todayDue, err := h.statsRepo.GetTodayDueCount(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to get today due count", http.StatusInternalServerError)
		return
	}

	todayNew, err := h.statsRepo.GetTodayNewCount(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to get today new count", http.StatusInternalServerError)
		return
	}

	totalMastered, err := h.statsRepo.GetMasteredCount(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to get mastered count", http.StatusInternalServerError)
		return
	}

	streakDays, err := h.statsRepo.GetStreakDays(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to get streak days", http.StatusInternalServerError)
		return
	}

	recentWords, err := h.statsRepo.GetRecentWords(r.Context(), userID, 5)
	if err != nil {
		http.Error(w, "Failed to get recent words", http.StatusInternalServerError)
		return
	}

	weeklyStats, err := h.statsRepo.GetWeeklyStats(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to get weekly stats", http.StatusInternalServerError)
		return
	}

	// Build response
	response := DashboardResponse{
		TodayDue:      todayDue,
		TodayNew:      todayNew,
		TotalMastered: totalMastered,
		StreakDays:    streakDays,
		RecentWords:   recentWords,
		WeeklyStats:   weeklyStats,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
