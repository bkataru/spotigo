package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Semantic search across your music library",
	Long: `Search your music library using natural language queries.

Examples:
  spotigo search "upbeat songs for working out"
  spotigo search "melancholic piano music"
  spotigo search "songs I listened to last summer"
  spotigo search "artists similar to Radiohead"

The search uses AI embeddings to understand meaning, not just keywords.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]
		runSearch(query)
	},
}

var (
	searchLimit  int
	searchType   string
	searchFormat string
)

func init() {
	searchCmd.Flags().IntVar(&searchLimit, "limit", 10, "maximum number of results")
	searchCmd.Flags().StringVar(&searchType, "type", "all", "search type: all, tracks, artists, albums, playlists")
	searchCmd.Flags().StringVar(&searchFormat, "format", "table", "output format: table, json, csv")
}

func runSearch(query string) {
	fmt.Printf("Searching for: \"%s\"\n", query)
	fmt.Printf("  Type: %s\n", searchType)
	fmt.Printf("  Limit: %d\n", searchLimit)
	fmt.Println()

	// RAG semantic search is not yet implemented
	fmt.Println("Semantic search is not yet implemented.")
	fmt.Println("This feature will use chromem-go for vector embeddings.")
	fmt.Println()
	fmt.Println("For now, run 'spotigo backup' to save your library locally.")
}
