package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/Flack74/pom/config"
	"github.com/Flack74/pom/logs"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "📊 View Pomodoro session statistics",
	Long: `📊 View Your Pomodoro Statistics

Track your productivity with detailed statistics about your Pomodoro sessions:
  • Today's progress toward daily goals
  • All-time session totals
  • Focus time tracking
  • Daily averages
  • Goal completion streaks
  • Task-specific statistics

Examples:
  pom stats            View all statistics
  pom stats --today    Show only today's stats
  pom stats --task     Show task-specific stats`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load theme
		theme, err := config.LoadTheme()
		if err != nil {
			theme = config.DefaultTheme
		}

		// Get all-time stats
		totalSessions, totalFocusMinutes, avgSessionsPerDay, err := logs.GetSessionStats()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s⚠️  Error getting stats: %v%s\n", theme.WarningColor, err, theme.TextColor)
			os.Exit(1)
		}

		// Get today's stats
		todaySessions, todayMinutes, err := logs.GetDailyStats()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s⚠️  Error getting daily stats: %v%s\n", theme.WarningColor, err, theme.TextColor)
			os.Exit(1)
		}

		// Get current goals
		goal, err := config.LoadGoal()
		if err != nil {
			goal = config.Goal{} // Empty goal if none set
		}

		// Print statistics header
		fmt.Printf("\n%s📊 Pomodoro Statistics%s\n\n", theme.HighlightColor, theme.TextColor)

		// Today's Progress
		fmt.Printf("%s🌅 Today's Progress%s\n", theme.SuccessColor, theme.TextColor)
		fmt.Printf("   Sessions completed: %d", todaySessions)
		if goal.DailySessionTarget > 0 {
			fmt.Printf(" / %d (%.1f%%)", goal.DailySessionTarget, float64(todaySessions)/float64(goal.DailySessionTarget)*100)
		}
		fmt.Println()
		fmt.Printf("   Focus time: %d minutes", todayMinutes)
		if goal.DailyMinutes > 0 {
			fmt.Printf(" / %d (%.1f%%)", goal.DailyMinutes, float64(todayMinutes)/float64(goal.DailyMinutes)*100)
		}
		fmt.Println()

		// All-time Stats
		fmt.Printf("\n%s🏆 All-time Statistics%s\n", theme.SuccessColor, theme.TextColor)
		fmt.Printf("   Total sessions: %d\n", totalSessions)
		fmt.Printf("   Total focus time: %.0f minutes (%.1f hours)\n", totalFocusMinutes, totalFocusMinutes/60)
		fmt.Printf("   Average sessions per day: %.1f\n", avgSessionsPerDay)

		// Get current streak if goals are set
		if goal.DailySessionTarget > 0 || goal.DailyMinutes > 0 {
			progress, err := config.LoadProgress()
			if err == nil {
				fmt.Printf("\n%s🔥 Goal Streaks%s\n", theme.SuccessColor, theme.TextColor)
				fmt.Printf("   Current streak: %d days\n", progress.CurrentStreak)
				fmt.Printf("   Longest streak: %d days\n", progress.LongestStreak)
			}
		}

		// Print motivational message based on stats
		fmt.Printf("\n%s%s%s\n", theme.HighlightColor, getMotivationalMessage(totalSessions, todaySessions), theme.TextColor)
	},
}

func getMotivationalMessage(totalSessions, todaySessions int) string {
	if totalSessions == 0 {
		return "🌱 Start your Pomodoro journey today!"
	}

	if todaySessions == 0 {
		return "🌅 A new day, new opportunities to focus!"
	}

	messages := []string{
		"🚀 Keep up the great work!",
		"💪 You're building great habits!",
		"⭐ Your dedication is inspiring!",
		"🎯 Stay focused, you're doing great!",
		"🌟 Every session counts!",
	}

	// Use current time as seed for variety
	seed := time.Now().UnixNano()
	return messages[seed%int64(len(messages))]
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
