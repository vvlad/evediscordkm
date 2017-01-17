package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/vvlad/evediscordkm"
)

func main() {
	configFile := flag.String("config", "", "configuration file path")
	flag.Parse()

	if _, err := os.Stat(*configFile); err != nil {
		fmt.Println("a config file is required and must be a valid file")
		os.Exit(1)
	}

	var config = evediscordkm.ReadFile(*configFile)
	var zKillBoard = evediscordkm.NewZKillboard(config)
	var reportProcessor = evediscordkm.NewReportProcessor(config)
	go zKillBoard.Run()
	go reportProcessor.Run()

	for report := range zKillBoard.Reports {
		reportProcessor.Reports <- report
	}

}
