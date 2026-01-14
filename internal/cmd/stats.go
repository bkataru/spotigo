package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "View your listening statistics",
	Long: `View detailed statistics about your listening habits.

Available statistics:
  - Top artists and tracks (all time, yearly, monthly)
  - Genre distribution
  - Listening time patterns
  - Discovery timeline
  - Playlist analysis`,
	Run: func(cmd *cobra.Command, args []string) {
		runStats()
	},
}

var (
	statsPeriod string
	statsTop    int
)

func init() {
	statsCmd.Flags().StringVar(&statsPeriod, "period", "all", "time period: all, year, month, week")
	statsCmd.Flags().IntVar(&statsTop, "top", 10, "number of top items to show")

	// Subcommands for specific stats
	statsCmd.AddCommand(statsTopCmd)
	statsCmd.AddCommand(statsGenresCmd)
	statsCmd.AddCommand(statsTimelineCmd)
}

var statsTopCmd = &cobra.Command{
	Use:   "top [artists|tracks|albums]",
	Short: "Show your top artists, tracks, or albums",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		itemType := "artists"
		if len(args) > 0 {
			itemType = args[0]
		}
		fmt.Printf("Top %d %s (%s):\n", statsTop, itemType, statsPeriod)
		fmt.Println("  (no data available - run backup first)")
	},
}

var statsGenresCmd = &cobra.Command{
	Use:   "genres",
	Short: "Show genre distribution",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Genre Distribution:")
		fmt.Println("  (no data available - run backup first)")
	},
}

var statsTimelineCmd = &cobra.Command{
	Use:   "timeline",
	Short: "Show listening timeline",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Listening Timeline:")
		fmt.Println("  (no data available - run backup first)")
	},
}

func runStats() {
	fmt.Println("Listening Statistics")
	fmt.Println("====================")
	fmt.Println()
	fmt.Println("No backup data found.")
	fmt.Println("Run 'spotigo backup' first to generate statistics.")
}
