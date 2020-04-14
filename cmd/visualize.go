package cmd

import (
	"fmt"
	util "github.com/RyanSiu1995/cloudbuild-visualizer/internal/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	if _, exist := os.LookupEnv("DEBUG"); exist {
		log.SetLevel(log.DebugLevel)
	}
}

func init() {
	var visualize = &cobra.Command{
		Use:   "visualize [cloud build YAML path]",
		Short: "Visualize the given cloud build pipeline",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			filePath := args[0]
			if _, err := os.Stat(filePath); err != nil {
				fmt.Printf("File %s does not exist!", filePath)
			}

			cloudbuild, err := util.ParseYaml(filePath)
			if err != nil {
				fmt.Printf("Cannot parse the YAML file %s", filePath)
			}

			graph := util.BuildStepsToDAG(cloudbuild.Steps)
			err = util.Visualize(graph)
			if err != nil {
				fmt.Println("There is error during rendering the image... Use DEBUG flag to get more information...")
				log.Debug(err.Error())
			}
		},
	}
	rootCmd.AddCommand(visualize)
}
