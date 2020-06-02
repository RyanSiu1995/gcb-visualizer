package main

import (
	"fmt"
	"os"

	"github.com/RyanSiu1995/gcb-visualizer/cmd"
	log "github.com/sirupsen/logrus"
)

func init() {
	if _, exist := os.LookupEnv("DEBUG"); exist {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	if err := cmd.Execute(); err != nil {
		log.Debug(err.Error())
		fmt.Println("Error encountered... Use DEBUG flag to see more error message...")
	}
}
