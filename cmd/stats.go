package cmd

import (
	"fmt"
	"os"

	"pom/logs"

	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "View Pomodoro session statistics",
	Long:  "Display statistics about your completed Pomodoro sessions, including total sessions and focus time.",
	Run: func(cmd *cobra.Command, args []string) {
		totalSessions, totalFocusMinutes, avgSessionsPerDay, err := logs.GetSessionStats()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting stats: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("%s📊 Pomodoro Statistics%s\n\n", colorPurple, colorReset)
		fmt.Printf("%s🎯 Total Sessions: %d%s\n", colorBlue, totalSessions, colorReset)
		fmt.Printf("%s⏰ Total Focus Time: %.0f minutes (%.1f hours)%s\n",
			colorGreen, totalFocusMinutes, totalFocusMinutes/60, colorReset)
		fmt.Printf("%s📈 Average Sessions Per Day: %.1f%s\n",
			colorYellow, avgSessionsPerDay, colorReset)
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
