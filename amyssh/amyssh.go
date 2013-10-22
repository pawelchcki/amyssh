package main

import (
	"flag"
	"log"
	"os"

	"github.com/pchojnacki/amyssh"
)

func main() {

	// configurator.Initialize(amyssh.DefaultConfig)
	// options := configurator.Options().(*amyssh.Config)

	// flag.StringVar(configurator.ConfigFilePath(), "f", "/etc/amyssh.yml", "config file location")
	// flag.DurationVar(&options.MaxPollInterval, "maxinterval", options.MaxPollInterval, "maximum interval at which datasource will be polled")
	// flag.DurationVar(&options.MinPollInterval, "mininterval", options.MinPollInterval, "minimum interval at which datasource will be polled")
	// flag.StringVar(&options.LogFilePath, "l", options.LogFilePath, "specify log file location, stdout when empty")

	flag.Parse()

	// cfg := configurator.Config().(amyssh.Config)
	if amyssh.GetConfig().LogFilePath != "" {
		file, err := os.OpenFile(amyssh.GetConfig().LogFilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		log.SetOutput(file)
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Printf("configuration loaded: %+v", amyssh.GetConfig())
	amyssh.IntervalLoop(amyssh.GetConfig(), amyssh.Perform)
}
