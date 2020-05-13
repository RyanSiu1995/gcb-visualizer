package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	if _, exist := os.LookupEnv("DEBUG"); exist {
		log.SetLevel(log.DebugLevel)
	}
	var visualize = &cobra.Command{
		Use:   "visualize [cloud build YAML path]",
		Short: "Visualize the given cloud build pipeline",
		Args:  cobra.ExactArgs(1),
		Run:   visualize,
	}
	visualize.Flags().StringVarP(&outputFormat, "output", "o", "", "Specify the output file, a temp png image by default (Support: dot, jpg and png)")
	rootCmd.AddCommand(visualize)
}
