package cmd

import (
	"fmt"

	"github.com/Flack74/pom/config"
	"github.com/spf13/cobra"
)

var insightsCmd = &cobra.Command{
	Use:   "insights",
	Short: "🧠 AI-powered insights and suggestions",
	Long: `🧠 AI-Powered Insights

Get personalized suggestions based on your Pomodoro history:
  • Optimal session lengths
  • Best times to focus
  • Productivity patterns
  • Performance improvements

Examples:
  pom insights suggest          Get AI suggestions
  pom insights calendar         View session calendar
  pom insights today           Today's statistics`,
}

var suggestCmd = &cobra.Command{
	Use:   "suggest",
	Short: "Get AI-powered suggestions for better productivity",
	Run: func(cmd *cobra.Command, args []string) {
		suggestions, err := config.GenerateSuggestions()
		if err != nil {
			fmt.Printf("Error generating suggestions: %v\n", err)
			return
		}

		if len(suggestions) == 0 {
			fmt.Println("🤖 No suggestions available yet. Complete a few sessions to get personalized recommendations!")
			return
		}

		fmt.Println("🧠 AI Suggestions for Better Productivity:\n")
		for i, suggestion := range suggestions {
			confidence := int(suggestion.Confidence * 100)
			fmt.Printf("%d. %s (Confidence: %d%%)\n", i+1, suggestion.Message, confidence)
			
			if suggestion.WorkTime > 0 {
				fmt.Printf("   Suggested work time: %d minutes\n", suggestion.WorkTime)
			}
			if suggestion.BreakTime > 0 {
				fmt.Printf("   Suggested break time: %d minutes\n", suggestion.BreakTime)
			}
			if suggestion.Sessions > 0 {
				fmt.Printf("   Suggested sessions: %d\n", suggestion.Sessions)
			}
			fmt.Println()
		}

		// Show current stats
		stats, err := config.AnalyzePerformance()
		if err == nil {
			fmt.Printf("📊 Your Stats:\n")
			fmt.Printf("   Completion Rate: %.1f%%\n", stats.CompletionRate*100)
			fmt.Printf("   Average Work Time: %.1f minutes\n", stats.AverageWorkTime)
			fmt.Printf("   Productivity Score: %.1f/100\n", stats.ProductivityScore)
		}
	},
}

var calendarCmd = &cobra.Command{
	Use:   "calendar",
	Short: "View your focus session calendar",
	Run: func(cmd *cobra.Command, args []string) {
		months, _ := cmd.Flags().GetInt("months")
		
		calendar, err := config.GenerateCalendarView(months)
		if err != nil {
			fmt.Printf("Error generating calendar: %v\n", err)
			return
		}

		fmt.Print(calendar)
	},
}

var todayCmd = &cobra.Command{
	Use:   "today",
	Short: "Show today's statistics",
	Run: func(cmd *cobra.Command, args []string) {
		sessions, minutes, err := config.GetTodayStats()
		if err != nil {
			fmt.Printf("Error getting today's stats: %v\n", err)
			return
		}

		fmt.Println("📅 Today's Focus Sessions:")
		fmt.Printf("   Sessions completed: %d\n", sessions)
		fmt.Printf("   Total focus time: %d minutes (%.1f hours)\n", minutes, float64(minutes)/60)
		
		if sessions == 0 {
			fmt.Println("   🎯 Ready to start your first session today?")
		} else if minutes >= 120 {
			fmt.Println("   🔥 Great job! You're having a productive day!")
		} else {
			fmt.Println("   💪 Keep going! You're building momentum!")
		}
	},
}

func init() {
	calendarCmd.Flags().Int("months", 3, "Number of months to show")
	
	insightsCmd.AddCommand(suggestCmd)
	insightsCmd.AddCommand(calendarCmd)
	insightsCmd.AddCommand(todayCmd)
	rootCmd.AddCommand(insightsCmd)
}