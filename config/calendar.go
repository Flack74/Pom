package config

import (
	"fmt"
	"strings"
	"time"
)

type CalendarDay struct {
	Date     time.Time
	Sessions int
	Minutes  int
	Level    int // 0-4 intensity level for heatmap
}

func GenerateCalendarView(months int) (string, error) {
	sessions, err := loadSessionHistory()
	if err != nil {
		return "", err
	}

	// Create map of date -> activity
	activityMap := make(map[string]CalendarDay)
	
	// Process sessions
	for _, session := range sessions {
		dateKey := session.Date.Format("2006-01-02")
		if day, exists := activityMap[dateKey]; exists {
			day.Sessions += session.Sessions
			day.Minutes += session.WorkMinutes * session.Sessions
			activityMap[dateKey] = day
		} else {
			activityMap[dateKey] = CalendarDay{
				Date:     session.Date,
				Sessions: session.Sessions,
				Minutes:  session.WorkMinutes * session.Sessions,
			}
		}
	}

	// Calculate intensity levels
	maxMinutes := 0
	for _, day := range activityMap {
		if day.Minutes > maxMinutes {
			maxMinutes = day.Minutes
		}
	}

	for dateKey, day := range activityMap {
		if maxMinutes > 0 {
			day.Level = int((float64(day.Minutes) / float64(maxMinutes)) * 4)
		}
		activityMap[dateKey] = day
	}

	// Generate calendar
	var result strings.Builder
	result.WriteString("📅 Focus Session Calendar (Last " + fmt.Sprintf("%d", months) + " months)\n\n")
	
	// Generate heatmap legend
	result.WriteString("Less ")
	for i := 0; i <= 4; i++ {
		result.WriteString(getHeatmapChar(i))
	}
	result.WriteString(" More\n\n")

	// Generate monthly view
	now := time.Now()
	for m := months - 1; m >= 0; m-- {
		monthStart := now.AddDate(0, -m, 0)
		monthStart = time.Date(monthStart.Year(), monthStart.Month(), 1, 0, 0, 0, 0, monthStart.Location())
		
		result.WriteString(fmt.Sprintf("%s %d\n", monthStart.Format("January"), monthStart.Year()))
		result.WriteString("Mo Tu We Th Fr Sa Su\n")
		
		// Add padding for first week
		firstWeekday := int(monthStart.Weekday())
		if firstWeekday == 0 {
			firstWeekday = 7 // Sunday = 7
		}
		for i := 1; i < firstWeekday; i++ {
			result.WriteString("   ")
		}
		
		// Add days of month
		daysInMonth := time.Date(monthStart.Year(), monthStart.Month()+1, 0, 0, 0, 0, 0, monthStart.Location()).Day()
		for day := 1; day <= daysInMonth; day++ {
			date := time.Date(monthStart.Year(), monthStart.Month(), day, 0, 0, 0, 0, monthStart.Location())
			dateKey := date.Format("2006-01-02")
			
			level := 0
			if dayData, exists := activityMap[dateKey]; exists {
				level = dayData.Level
			}
			
			result.WriteString(fmt.Sprintf("%s%2d", getHeatmapChar(level), day))
			
			if (firstWeekday+day-1)%7 == 0 {
				result.WriteString("\n")
			} else {
				result.WriteString(" ")
			}
		}
		result.WriteString("\n\n")
	}

	// Add summary stats
	totalSessions := 0
	totalMinutes := 0
	for _, day := range activityMap {
		totalSessions += day.Sessions
		totalMinutes += day.Minutes
	}
	
	result.WriteString(fmt.Sprintf("📊 Total: %d sessions, %d minutes (%.1f hours)\n", 
		totalSessions, totalMinutes, float64(totalMinutes)/60))

	return result.String(), nil
}

func getHeatmapChar(level int) string {
	chars := []string{"⬜", "🟩", "🟨", "🟧", "🟥"}
	if level < 0 || level >= len(chars) {
		return chars[0]
	}
	return chars[level]
}

func GetTodayStats() (int, int, error) {
	sessions, err := loadSessionHistory()
	if err != nil {
		return 0, 0, err
	}

	today := time.Now().Format("2006-01-02")
	totalSessions := 0
	totalMinutes := 0

	for _, session := range sessions {
		if session.Date.Format("2006-01-02") == today {
			totalSessions += session.Sessions
			totalMinutes += session.WorkMinutes * session.Sessions
		}
	}

	return totalSessions, totalMinutes, nil
}