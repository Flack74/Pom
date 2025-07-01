package cmd

import (
	"fmt"

	"github.com/Flack74/pom/config"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "🔄 Cloud synchronization",
	Long: `🔄 Cloud Sync (Optional Backup)

Sync your Pomodoro data across devices using cloud storage:
  • GitHub: Use a private repository for sync
  • Dropbox: Use Dropbox for file sync (requires rclone)

Setup:
  export POM_GITHUB_REPO="https://github.com/user/pom-data"
  export POM_GITHUB_TOKEN="your_token"
  
  OR
  
  export POM_DROPBOX_TOKEN="your_token"

Examples:
  pom sync setup github         Configure GitHub sync
  pom sync setup dropbox        Configure Dropbox sync
  pom sync push                 Upload data to cloud
  pom sync pull                 Download data from cloud
  pom sync status               Check sync configuration`,
}

var syncSetupCmd = &cobra.Command{
	Use:   "setup [provider]",
	Short: "Configure cloud sync provider",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		provider := args[0]
		
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		switch provider {
		case "github":
			cfg.CloudSync = true
			cfg.CloudProvider = "github"
			fmt.Println("🔧 GitHub sync configured!")
			fmt.Println("Set these environment variables:")
			fmt.Println("  export POM_GITHUB_REPO=\"https://github.com/user/pom-data\"")
			fmt.Println("  export POM_GITHUB_TOKEN=\"your_personal_access_token\"")
		case "dropbox":
			cfg.CloudSync = true
			cfg.CloudProvider = "dropbox"
			fmt.Println("🔧 Dropbox sync configured!")
			fmt.Println("Install rclone and set:")
			fmt.Println("  export POM_DROPBOX_TOKEN=\"your_access_token\"")
		default:
			fmt.Printf("Unknown provider: %s. Use 'github' or 'dropbox'\n", provider)
			return
		}

		if err := config.SaveConfig(cfg); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}

		fmt.Printf("✅ Cloud sync enabled with %s\n", provider)
	},
}

var syncPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Upload data to cloud",
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.SyncData(true); err != nil {
			fmt.Printf("Error uploading data: %v\n", err)
			return
		}
		fmt.Println("✅ Data uploaded to cloud")
	},
}

var syncPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Download data from cloud",
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.SyncData(false); err != nil {
			fmt.Printf("Error downloading data: %v\n", err)
			return
		}
		fmt.Println("✅ Data downloaded from cloud")
	},
}

var syncStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check sync configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		fmt.Println("🔄 Cloud Sync Status:")
		fmt.Printf("   Enabled: %t\n", cfg.CloudSync)
		fmt.Printf("   Provider: %s\n", cfg.CloudProvider)
		
		if cfg.CloudSync {
			provider := config.GetSyncProvider(cfg)
			if provider != nil {
				fmt.Printf("   Available: %t\n", provider.IsAvailable())
			}
		}
	},
}

func init() {
	syncCmd.AddCommand(syncSetupCmd)
	syncCmd.AddCommand(syncPushCmd)
	syncCmd.AddCommand(syncPullCmd)
	syncCmd.AddCommand(syncStatusCmd)
	rootCmd.AddCommand(syncCmd)
}