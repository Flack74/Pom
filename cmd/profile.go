package cmd

import (
	"fmt"
	"strconv"

	"github.com/Flack74/pom/config"
	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "👥 Manage work profiles",
	Long: `👥 Multi-Profile Support

Create and manage different profiles for various work contexts:
  • Work profile: Longer sessions for deep focus
  • Study profile: Balanced sessions for learning
  • Quick profile: Short bursts for small tasks
  • Custom profiles: Tailored to your needs

Examples:
  pom profile list                    List all profiles
  pom profile use work               Switch to work profile
  pom profile create "coding" 45 10 3  Create custom profile
  pom profile delete "old-profile"   Remove a profile`,
}

var listProfilesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available profiles",
	Run: func(cmd *cobra.Command, args []string) {
		profiles, err := config.LoadProfiles()
		if err != nil {
			fmt.Printf("Error loading profiles: %v\n", err)
			return
		}

		cfg, _ := config.LoadConfig()
		
		fmt.Println("\n👥 Available Profiles:")
		for _, profile := range profiles.Profiles {
			current := ""
			if profile.Name == cfg.CurrentProfile {
				current = " (current)"
			}
			fmt.Printf("  %s%s\n", profile.Name, current)
			fmt.Printf("    %s\n", profile.Description)
			fmt.Printf("    Work: %dm, Break: %dm, Sessions: %d\n\n", 
				profile.WorkMinutes, profile.BreakMinutes, profile.NumSessions)
		}
	},
}

var useProfileCmd = &cobra.Command{
	Use:   "use [profile-name]",
	Short: "Switch to a different profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		profileName := args[0]
		
		// Check if profile exists
		_, err := config.GetProfile(profileName)
		if err != nil {
			fmt.Printf("Profile '%s' not found\n", profileName)
			return
		}

		// Update config
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		cfg.CurrentProfile = profileName
		if err := config.SaveConfig(cfg); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}

		fmt.Printf("✅ Switched to profile: %s\n", profileName)
	},
}

var createProfileCmd = &cobra.Command{
	Use:   "create [name] [work-minutes] [break-minutes] [sessions]",
	Short: "Create a new profile",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		workMin, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("Invalid work minutes: %s\n", args[1])
			return
		}
		breakMin, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Printf("Invalid break minutes: %s\n", args[2])
			return
		}
		sessions, err := strconv.Atoi(args[3])
		if err != nil {
			fmt.Printf("Invalid sessions: %s\n", args[3])
			return
		}

		description, _ := cmd.Flags().GetString("description")
		if description == "" {
			description = fmt.Sprintf("Custom profile: %dm work, %dm break", workMin, breakMin)
		}

		profile := config.Profile{
			Name:         name,
			WorkMinutes:  workMin,
			BreakMinutes: breakMin,
			NumSessions:  sessions,
			Description:  description,
		}

		if err := config.AddProfile(profile); err != nil {
			fmt.Printf("Error creating profile: %v\n", err)
			return
		}

		fmt.Printf("✅ Created profile: %s\n", name)
	},
}

func init() {
	createProfileCmd.Flags().String("description", "", "Profile description")
	
	profileCmd.AddCommand(listProfilesCmd)
	profileCmd.AddCommand(useProfileCmd)
	profileCmd.AddCommand(createProfileCmd)
	rootCmd.AddCommand(profileCmd)
}