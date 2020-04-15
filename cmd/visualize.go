package cmd

import (
	"fmt"
	util "github.com/RyanSiu1995/cloudbuild-visualizer/internal/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var supportedFormat = []string{".dot", ".jpg", ".jpeg", ".png"}
var output = ""

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
				fmt.Printf("File %s does not exist!\n", filePath)
				os.Exit(1)
			}

			cloudbuild, err := util.ParseYaml(filePath)
			if err != nil {
				fmt.Printf("Cannot parse the YAML file %s\n", filePath)
				os.Exit(1)
			}

			graph := util.BuildStepsToDAG(cloudbuild.Steps)
			if output == "" {
				err = util.Visualize(graph)
				if err != nil {
					fmt.Println("There is error during rendering the image... Use DEBUG flag to get more information...")
					log.Debug(err.Error())
					os.Exit(1)
				}
			} else {
				ext := filepath.Ext(output)
				if !util.Contains(supportedFormat, strings.ToLower(ext)) {
					fmt.Printf("%s is not the supported format. Please specify the correct extension\n", ext)
					os.Exit(1)
				}

				err = util.SaveGraph(graph, output)
				if err != nil {
					fmt.Println("There is error during rendering the image... Use DEBUG flag to get more information...")
					log.Debug(err.Error())
					os.Exit(1)
				}
			}
		},
	}
	visualize.Flags().StringVarP(&output, "output", "o", "", "Specify the output file, a temp png image by default (Support: dot, jpg and png)")
	rootCmd.AddCommand(visualize)
}
