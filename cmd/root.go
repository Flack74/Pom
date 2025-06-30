package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	version   = "1.0.0"
	buildDate = "unknown"
	gitCommit = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "pom",
	Short: "🍅 A beautiful and feature-rich Pomodoro timer for your terminal",
	Long: `🍅 Pom - Your Friendly Terminal Pomodoro Timer

A beautiful and feature-rich Pomodoro timer that helps you stay focused and productive.
Built with love using the time-tested Pomodoro Technique®.

Key Features:
  • Beautiful progress bar with real-time updates
  • Multiple color themes (default, minimal, vibrant)
  • Daily goals and streak tracking
  • Task planning and time tracking
  • Comprehensive statistics
  • Desktop notifications
  • Motivational messages

Quick Start:
  pom start              Start a default session (25min work + 5min break)
  pom start -w 30 -b 5   Custom work/break duration
  pom stats              View your progress
  pom theme set vibrant  Change the look and feel

Need help? Try 'pom [command] --help' for detailed information.`,
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is provided, show help
		cmd.Help()
	},
}

func init() {
	// Set custom version template
	rootCmd.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "version %s" .Version}}
Build Date: {{printf "%s" .buildDate}}
Git Commit: {{printf "%s" .gitCommit}}
Go Version: {{printf "%s" .goVersion}}
OS/Arch:    {{printf "%s/%s" .os .arch}}
`)
}

func Execute() {
	// Add build information to root command
	rootCmd.SetVersionTemplate(rootCmd.VersionTemplate())
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  "Print detailed version information about the pom binary",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version:    %s\n", version)
			fmt.Printf("Build Date: %s\n", buildDate)
			fmt.Printf("Git Commit: %s\n", gitCommit)
			fmt.Printf("Go Version: %s\n", runtime.Version())
			fmt.Printf("OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
		},
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops. An error while executing Pom '%s'\n", err)
		os.Exit(1)
	}
}
