package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	util "github.com/RyanSiu1995/gcb-visualizer/internal/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var supportedFormat = []string{".dot", ".jpg", ".jpeg", ".png"}
var outputFormat = ""

const ver string = "1.0.1"

func visualize(cmd *cobra.Command, args []string) {
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
	if outputFormat == "" {
		err = util.Visualize(graph)
		if err != nil {
			fmt.Println("There is error during rendering the image... Use DEBUG flag to get more information...")
			log.Debug(err.Error())
			os.Exit(1)
		}
	} else {
		ext := filepath.Ext(outputFormat)
		if !util.Contains(supportedFormat, strings.ToLower(ext)) {
			fmt.Printf("%s is not the supported format. Please specify the correct extension\n", ext)
			os.Exit(1)
		}

		err = util.SaveGraph(graph, outputFormat)
		if err != nil {
			fmt.Println("There is error during rendering the image... Use DEBUG flag to get more information...")
			log.Debug(err.Error())
			os.Exit(1)
		}
	}
}
