package main

import (
	amyssh "./amyssh"
	configurator "./configurator"
	"flag"
	// "fmt"
)

func main() {
	configurator.Initialize(amyssh.DefaultConfig)
	options := configurator.Options().(*amyssh.Config)

	flag.StringVar(configurator.ConfigFilePath(), "f", "/etc/amyssh.yml", "config file location")
	flag.UintVar(&options.Database.Port, "dbport", amyssh.DefaultConfig.Database.Port, "database port")
	flag.StringVar(&options.Database.Host, "dbhost", amyssh.DefaultConfig.Database.Host, "database host")
	flag.StringVar(&options.Database.User, "dbuser", amyssh.DefaultConfig.Database.User, "database user")
	flag.StringVar(&options.Database.Password, "dbpassword", amyssh.DefaultConfig.Database.Password, "database password")
	flag.Parse()
	cfg := configurator.Config().(amyssh.Config)

	// db := amyssh.NewCon(cfg)
	// hostTags := []string{"stage", "prod"}
	// userTags := []string{"admin", "deployer"}

	// keys, _ := db.FetchKeys(hostTags, userTags)
	// amyssh.Show(keys)

	amyssh.IntervalLoop(&cfg, amyssh.Perform)
}
