package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pom",
	Short: "pom is a cli pomodoro timer",
	Long:  "A command-line timer that follows the Pomodoro Technique — 25 min focus + 5 min break.",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops. An error while executing Pom '%s'\n", err)
		os.Exit(1)
	}
}
