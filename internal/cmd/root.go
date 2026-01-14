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
		if err := cmd.Help(); err != nil {
			fmt.Printf("Error showing help: %v\n", err)
		}
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
	if err := viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose")); err != nil {
		fmt.Printf("Error binding verbose flag: %v\n", err)
	}

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
	model := tui.InitialModel()
	p := bubbletea.NewProgram(model, bubbletea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		return
	}

	// Type assert to get the final model state
	m, ok := finalModel.(tui.Model)
	if !ok {
		fmt.Println("Error: unexpected model type")
		return
	}

	// Check if user quit without selecting
	if m.IsQuitting() {
		fmt.Println("Goodbye!")
		return
	}

	// Get the selected command
	command := m.GetSelectedCommand()
	if command == "" {
		return
	}

	fmt.Printf("\nSelected: %s\n\n", m.GetSelectedDisplay())

	// Execute the corresponding command
	switch command {
	case "backup":
		backupCmd.Run(backupCmd, []string{})
	case "chat":
		chatCmd.Run(chatCmd, []string{})
	case "search":
		fmt.Println("Enter your search query:")
		fmt.Println("  spotigo search \"your query here\"")
		fmt.Println()
		fmt.Println("Or run 'spotigo search index' first to build the search index.")
	case "stats":
		statsCmd.Run(statsCmd, []string{})
	case "auth-status":
		authStatusCmd.Run(authStatusCmd, []string{})
	case "models-status":
		modelsStatusCmd.Run(modelsStatusCmd, []string{})
	case "exit":
		fmt.Println("Goodbye!")
	}
}
