package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "gcb-visualizer",
	Short: "GCB visualizer will visualize the google cloud build pipeline",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
}

// Execute will execute the cobra-based command
func Execute() error {
	return rootCmd.Execute()
}
