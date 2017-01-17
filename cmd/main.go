package main

import (
	"flag"

	"github.com/vvlad/evediscordkm"
)

func main() {
	configFile := flag.String("config", "", "configuration file path")
	flag.Parse()
	var config = evediscordkm.ReadFile(*configFile)
	var zKillBoard = evediscordkm.NewZKillboard(config)
	var reportProcessor = evediscordkm.NewReportProcessor(config)
	go zKillBoard.Run()
	go reportProcessor.Run()

	for report := range zKillBoard.Reports {
		reportProcessor.Reports <- report
	}

}
