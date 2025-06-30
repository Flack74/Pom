package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Flack74/pom/config"

	"github.com/spf13/cobra"
)

var goalsCmd = &cobra.Command{
	Use:   "goals",
	Short: "🎯 Manage daily Pomodoro goals",
	Long: `🎯 Set and Track Your Daily Goals

Stay motivated by setting and tracking daily Pomodoro goals:
  • Set target number of sessions per day
  • Set target focus minutes per day
  • Track your progress
  • Build and maintain streaks
  • Get motivational feedback

Examples:
  pom goals set 8 240     Set goal: 8 sessions, 240 minutes
  pom goals show          View current goals and progress
  pom goals reset         Reset today's progress`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.ShowProgress(); err != nil {
			fmt.Printf("Error showing progress: %v\n", err)
		}
	},
}

var setGoalCmd = &cobra.Command{
	Use:   "set [sessions] [minutes]",
	Short: "Set daily goals",
	Long:  "Set daily goals for number of sessions and total focus minutes",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		sessions, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Invalid number of sessions: %v\n", err)
			return
		}

		minutes, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("Invalid number of minutes: %v\n", err)
			return
		}

		goal := config.Goal{
			DailySessionTarget: sessions,
			DailyMinutes:       minutes,
			StartDate:          time.Now(),
		}

		if err := config.SaveGoal(goal); err != nil {
			fmt.Printf("Error saving goal: %v\n", err)
			return
		}

		fmt.Printf("Daily goals set: %d sessions, %d minutes\n", sessions, minutes)
	},
}

var showGoalCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current goals and progress",
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.ShowProgress(); err != nil {
			fmt.Printf("Error showing progress: %v\n", err)
		}
	},
}

func init() {
	goalsCmd.AddCommand(setGoalCmd)
	goalsCmd.AddCommand(showGoalCmd)
	rootCmd.AddCommand(goalsCmd)
}
