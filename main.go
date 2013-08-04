package main

import (
	config "./configurator"

	"flag"
	"fmt"
)

func defaultConfig() *AmySSHConfig {
	return &AmySSHConfig{
		DatabaseConfig{"localhost", 3306, "root", ""},
		[]UsersConfig{UsersConfig{"deployer", []string{"deploy", "all"}}},
		[]string{"default"},
	}
}

type DatabaseConfig struct {
	Host     string
	Port     uint
	User     string
	Password string
}
type UsersConfig struct {
	Name string
	Tags []string
}
type AmySSHConfig struct {
	Database DatabaseConfig
	Users    []UsersConfig
	HostTags []string
}

func main() {
	defaultCfg := *defaultConfig()
	config.Initialize(defaultCfg)

	options := config.Options().(*AmySSHConfig)

	flag.StringVar(config.ConfigFilePath(), "f", "/etc/amyssh.yml", "config file location")
	flag.UintVar(&options.Database.Port, "dbport", defaultCfg.Database.Port, "database port")
	flag.StringVar(&options.Database.Host, "dbhost", defaultCfg.Database.Host, "database host")
	flag.StringVar(&options.Database.User, "dbuser", defaultCfg.Database.User, "database user")
	flag.StringVar(&options.Database.Password, "dbpassword", defaultCfg.Database.Password, "database password")
	flag.Parse()
	fmt.Printf("%+v\n", config.Config().(AmySSHConfig))
}
