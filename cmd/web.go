package cmd

import (
	"github.com/Flack74/pom/web"
	"github.com/spf13/cobra"
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "🌐 Start web UI server",
	Long: `🌐 Web UI Bridge

Launch a web dashboard to access all Pomodoro features from any device:
  • Modern React interface with Galactic Flux theme
  • Real-time timer with WebSocket updates
  • All CLI features accessible via web
  • Cross-platform compatibility (Windows, Mac, Linux)

Examples:
  pom web                Start on default port 8080
  pom web -p 3000        Start on custom port`,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")
		server := web.NewServer()
		server.Start(port)
	},
}

func init() {
	webCmd.Flags().IntP("port", "p", 8080, "Port to run web server on")
	rootCmd.AddCommand(webCmd)
}