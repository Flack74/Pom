package config

import (
	"fmt"
	"math"
	"time"
)

type SessionStats struct {
	AverageWorkTime    float64
	AverageBreakTime   float64
	CompletionRate     float64
	PreferredTimeSlots []int // Hours of day
	ProductivityScore  float64
}

type Suggestion struct {
	Type        string `json:"type"`
	Message     string `json:"message"`
	WorkTime    int    `json:"work_time,omitempty"`
	BreakTime   int    `json:"break_time,omitempty"`
	Sessions    int    `json:"sessions,omitempty"`
	Confidence  float64 `json:"confidence"`
}

func AnalyzePerformance() (SessionStats, error) {
	sessions, err := loadSessionHistory()
	if err != nil {
		return SessionStats{}, err
	}

	if len(sessions) == 0 {
		return SessionStats{
			AverageWorkTime:   25,
			AverageBreakTime:  5,
			CompletionRate:    0,
			ProductivityScore: 0,
		}, nil
	}

	var totalWork, totalBreak float64
	var completed int
	hourMap := make(map[int]int)

	for _, session := range sessions {
		totalWork += float64(session.WorkMinutes)
		totalBreak += float64(session.BreakMinutes)
		if session.Completed {
			completed++
		}
		hourMap[session.Date.Hour()]++
	}

	stats := SessionStats{
		AverageWorkTime:  totalWork / float64(len(sessions)),
		AverageBreakTime: totalBreak / float64(len(sessions)),
		CompletionRate:   float64(completed) / float64(len(sessions)),
	}

	// Find preferred time slots
	maxCount := 0
	for hour, count := range hourMap {
		if count > maxCount {
			maxCount = count
			stats.PreferredTimeSlots = []int{hour}
		} else if count == maxCount {
			stats.PreferredTimeSlots = append(stats.PreferredTimeSlots, hour)
		}
	}

	// Calculate productivity score
	stats.ProductivityScore = stats.CompletionRate * 100

	return stats, nil
}

func GenerateSuggestions() ([]Suggestion, error) {
	stats, err := AnalyzePerformance()
	if err != nil {
		return nil, err
	}

	var suggestions []Suggestion

	// Work time suggestions
	if stats.CompletionRate < 0.7 && stats.AverageWorkTime > 25 {
		suggestions = append(suggestions, Suggestion{
			Type:       "work_time",
			Message:    "Consider shorter work sessions to improve completion rate",
			WorkTime:   int(math.Max(15, stats.AverageWorkTime-5)),
			Confidence: 0.8,
		})
	} else if stats.CompletionRate > 0.9 && stats.AverageWorkTime < 45 {
		suggestions = append(suggestions, Suggestion{
			Type:       "work_time",
			Message:    "You're doing great! Try longer sessions for deeper focus",
			WorkTime:   int(math.Min(45, stats.AverageWorkTime+5)),
			Confidence: 0.7,
		})
	}

	// Break time suggestions
	if stats.CompletionRate < 0.6 {
		suggestions = append(suggestions, Suggestion{
			Type:       "break_time",
			Message:    "Longer breaks might help you stay focused",
			BreakTime:  int(stats.AverageBreakTime + 2),
			Confidence: 0.6,
		})
	}

	// Time-based suggestions
	currentHour := time.Now().Hour()
	isPreferredTime := false
	for _, hour := range stats.PreferredTimeSlots {
		if abs(currentHour-hour) <= 1 {
			isPreferredTime = true
			break
		}
	}

	if !isPreferredTime && len(stats.PreferredTimeSlots) > 0 {
		suggestions = append(suggestions, Suggestion{
			Type:       "timing",
			Message:    fmt.Sprintf("You're most productive around %d:00. Consider scheduling sessions then.", stats.PreferredTimeSlots[0]),
			Confidence: 0.5,
		})
	}

	// Productivity suggestions
	if stats.ProductivityScore < 50 {
		suggestions = append(suggestions, Suggestion{
			Type:       "general",
			Message:    "Try the 'quick' profile for easier wins, then gradually increase session length",
			Confidence: 0.7,
		})
	}

	return suggestions, nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}