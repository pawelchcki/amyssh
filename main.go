package main

import (
	configurator "./configurator"
	amyssh "./lib"
	"flag"
	"log"
)

func main() {
	configurator.Initialize(amyssh.DefaultConfig)
	options := configurator.Options().(*amyssh.Config)

	flag.StringVar(configurator.ConfigFilePath(), "f", "/etc/amyssh.yml", "config file location")
	flag.DurationVar(&options.MaxPollInterval, "maxinterval", options.MaxPollInterval, "maximum interval at which datasource will be polled")
	flag.DurationVar(&options.MinPollInterval, "mininterval", options.MinPollInterval, "minimum interval at which datasource will be polled")

	flag.Parse()

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Printf("loading\n")
	cfg := configurator.Config().(amyssh.Config)

	amyssh.IntervalLoop(&cfg, amyssh.Perform)
}
