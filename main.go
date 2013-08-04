package main

import (
	config "./configurator"

	"flag"
	"fmt"
)

var defaultConfig = AmySSHConfig{
	DatabaseConfig{"localhost", 3306, "root", ""},
	[]UsersConfig{UsersConfig{"deployer", []string{"deploy", "admin"}, []string{}}},
	[]string{"default"},
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
	Keys []string
}
type AmySSHConfig struct {
	Database DatabaseConfig
	Users    []UsersConfig
	HostTags []string
}

func main() {
	config.Initialize(defaultConfig)

	options := config.Options().(*AmySSHConfig)

	flag.StringVar(config.ConfigFilePath(), "f", "/etc/amyssh.yml", "config file location")
	flag.UintVar(&options.Database.Port, "dbport", defaultConfig.Database.Port, "database port")
	flag.StringVar(&options.Database.Host, "dbhost", defaultConfig.Database.Host, "database host")
	flag.StringVar(&options.Database.User, "dbuser", defaultConfig.Database.User, "database user")
	flag.StringVar(&options.Database.Password, "dbpassword", defaultConfig.Database.Password, "database password")
	flag.Parse()
	fmt.Printf("%+v\n", config.Config().(AmySSHConfig))
}
