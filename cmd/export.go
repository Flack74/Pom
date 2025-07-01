package cmd

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/Flack74/pom/config"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "📤 Export your data",
	Long: `📤 Export/Import Your Data

Export your Pomodoro data for backup or analysis:
  • JSON format: Complete data export
  • CSV format: Session data for spreadsheets
  • Import from JSON backups

Examples:
  pom export json backup.json       Export all data to JSON
  pom export csv sessions.csv       Export sessions to CSV
  pom import backup.json            Import from JSON backup`,
}

var exportJSONCmd = &cobra.Command{
	Use:   "json [filename]",
	Short: "Export all data to JSON format",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := "pom-backup-" + time.Now().Format("2006-01-02") + ".json"
		if len(args) > 0 {
			filename = args[0]
		}

		// Ensure .json extension
		if filepath.Ext(filename) != ".json" {
			filename += ".json"
		}

		if err := config.ExportToJSON(filename); err != nil {
			fmt.Printf("Error exporting to JSON: %v\n", err)
			return
		}

		fmt.Printf("✅ Data exported to: %s\n", filename)
	},
}

var exportCSVCmd = &cobra.Command{
	Use:   "csv [filename]",
	Short: "Export session data to CSV format",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := "pom-sessions-" + time.Now().Format("2006-01-02") + ".csv"
		if len(args) > 0 {
			filename = args[0]
		}

		// Ensure .csv extension
		if filepath.Ext(filename) != ".csv" {
			filename += ".csv"
		}

		if err := config.ExportToCSV(filename); err != nil {
			fmt.Printf("Error exporting to CSV: %v\n", err)
			return
		}

		fmt.Printf("✅ Session data exported to: %s\n", filename)
	},
}

var importCmd = &cobra.Command{
	Use:   "import [filename]",
	Short: "Import data from JSON backup",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]

		if err := config.ImportFromJSON(filename); err != nil {
			fmt.Printf("Error importing data: %v\n", err)
			return
		}

		fmt.Printf("✅ Data imported from: %s\n", filename)
	},
}

func init() {
	exportCmd.AddCommand(exportJSONCmd)
	exportCmd.AddCommand(exportCSVCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(importCmd)
}