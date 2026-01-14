// Package cmd provides the CLI commands for Spotigo
package cmd

import (
	"fmt"

	bubbletea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"spotigo/internal/config"
	"spotigo/internal/tui"
)

var (
	cfgFile string
	tuiMode bool
	cfg     *config.Config
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "spotigo",
	Short: "AI-powered local music intelligence platform",
	Long: `Spotigo 2.0 - Your personal AI music assistant

A fully offline, privacy-first tool that backs up your Spotify library
and provides AI-powered insights about your music taste.

Features:
  - Complete Spotify library backup (playlists, saved tracks, artists)
  - AI chat for music discussions and recommendations
  - Semantic search across your library
  - Listening statistics and insights
  - TUI and web interfaces

All data stays local. All AI runs locally via Ollama.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		if tuiMode {
			runTUI()
			return
		}
		// Default: show help
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.spotigo.yaml)")
	rootCmd.PersistentFlags().BoolVar(&tuiMode, "tui", false, "launch in TUI mode")
	rootCmd.PersistentFlags().Bool("verbose", false, "enable verbose output")

	// Bind flags to viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	// Add subcommands
	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(chatCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(statsCmd)
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(modelsCmd)
}

func initConfig() error {
	var err error
	cfg, err = config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	return nil
}

// GetConfig returns the loaded configuration
func GetConfig() *config.Config {
	return cfg
}

func runTUI() {
	p := bubbletea.NewProgram(tui.InitialModel(), bubbletea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		return
	}

	// Handle the selected choice based on cursor position
	choice := tui.GetSelectedChoice(tui.InitialModel())
	if choice == "" {
		return
	}

	fmt.Printf("Selected: %s\n", choice)

	// Execute the corresponding command
	switch choice {
	case "🎵 Backup Library":
		fmt.Println("Run: spotigo backup")
	case "💬 AI Chat":
		fmt.Println("Run: spotigo chat")
	case "🔍 Search Music":
		fmt.Println("Run: spotigo search \"your query\"")
	case "📊 Statistics":
		fmt.Println("Run: spotigo stats")
	case "🔑 Auth Status":
		fmt.Println("Run: spotigo auth status")
	case "🤖 Models Status":
		fmt.Println("Run: spotigo models status")
	case "❌ Exit":
		fmt.Println("Goodbye!")
	}
}
