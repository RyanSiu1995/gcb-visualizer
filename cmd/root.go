package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

func init() {
	rootCmd = &cobra.Command{
		Use:   "gcb-visualizer [cloud build YAML path]",
		Short: "GCB visualizer will visualize the google cloud build pipeline",
		Args:  cobra.ExactArgs(1),
		Run:   visualize,
	}
	rootCmd.Flags().StringVarP(&outputFormat, "output", "o", "", "Specify the output file, a temp png image by default (Support: dot, jpg and png)")
}

// Execute will execute the cobra-based command
func Execute() error {
	return rootCmd.Execute()
}
