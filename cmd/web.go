package cmd

import (
	"fmt"
	"os"
	"github.com/Flack74/pom/web"
	"github.com/spf13/cobra"
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "🌐 Start web UI server",
	Long: `🌐 Web UI Bridge

Launch a web dashboard to access all Pomodoro features from any device:
  • Embedded Galactic Flux theme interface
  • Real-time timer with progress visualization
  • All CLI features accessible via web
  • Cross-platform compatibility (Windows, Mac, Linux)
  • Background daemon mode available

Examples:
  pom web                Start on default port 8080
  pom web -p 3000        Start on custom port
  pom web -d             Run in background (daemon mode)`,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")
		daemon, _ := cmd.Flags().GetBool("daemon")
		
		if daemon {
			fmt.Printf("🌐 Web UI starting in background on port %d\n", port)
			fmt.Printf("🔗 Access at: http://localhost:%d\n", port)
			fmt.Println("ℹ️  Use 'pkill pom' to stop the server")
			go func() {
				server := web.NewServer()
				if err := server.Start(port); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to start web server: %v\n", err)
					os.Exit(1)
				}
			}()
			// Keep process alive
			select {}
		} else {
			server := web.NewServer()
			if err := server.Start(port); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to start web server: %v\n", err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	webCmd.Flags().IntP("port", "p", 8080, "Port to run web server on")
	webCmd.Flags().BoolP("daemon", "d", false, "Run in background (daemon mode)")
	rootCmd.AddCommand(webCmd)
}