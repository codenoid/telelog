package main

import (
	"flag"
	"fmt"

	"github.com/codenoid/telelog"
	"github.com/hpcloud/tail"
)

var file string
var level string

func init() {
	fmt.Println("Telelog, make sure you already set these env : ")
	fmt.Println("TELELOG_BOT_TOKEN")
	fmt.Println("TELELOG_APP_NAME")
	fmt.Println("TELELOG_DEBUG_MODE (optional)")
	fmt.Println("TELELOG_RECIPIENT_LIST")

	flag.StringVar(&file, "file", "", "path to file to tail")
	flag.StringVar(&level, "level", "debug", "level (error|warn|info|debug")

	flag.Parse()

}

func main() {

	logger := telelog.LoggerNew()

	logger.SetAppName(file)
	if level == "debug" {
		logger.SetDebug(true)
	}

	t, err := tail.TailFile(file, tail.Config{Follow: true, MustExist: true, ReOpen: true})
	if err != nil {
		panic(err)
	}

	for line := range t.Lines {
		switch level {
		case "error":
			logger.Error(line.Text)
		case "warn", "warning":
			logger.Warn(line.Text)
		case "info":
			logger.Info(line.Text)
		default:
			logger.Debug(line.Text)
		}
	}

}
