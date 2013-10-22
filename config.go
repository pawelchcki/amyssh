package amyssh

import (
	"time"

	"github.com/pchojnacki/configurator"
)

var DefaultConfig = Config{
	Database: DatabaseConfig{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: "",
		DbName:   "amyssh"},
	Users: []UsersConfig{
		UsersConfig{
			Name: "deployer",
			Tags: []string{"deploy", "admin"},
			Keys: []string{},
		},
	},
	HostTags:               []string{"default"},
	MinPollInterval:        500 * time.Millisecond,
	MaxPollInterval:        10 * time.Second,
	PerformanceThreshold:   5 * time.Millisecond,  // Interval will be decreased if whole operation will take less than this
	BackoffThreshold:       20 * time.Millisecond, // Backoff threshold
	AuthorizedKeysFileName: "authorized_keys2",
	LogFilePath:            "",
}

type DatabaseConfig struct {
	Host     string
	Port     uint
	User     string
	Password string
	DbName   string
}

type UsersConfig struct {
	Name string
	Tags []string
	Keys []string
}

type Config struct {
	Database        DatabaseConfig
	Users           []UsersConfig
	HostTags        []string
	MinPollInterval time.Duration
	MaxPollInterval time.Duration

	PerformanceThreshold time.Duration
	BackoffThreshold     time.Duration

	AuthorizedKeysFileName string
	LogFilePath            string
}

var globalCachedConfig *Config
var globalConfigurator *configurator.ConfigHolder

func init() {
	globalConfigurator = configurator.NewConfig(Config{})
	globalConfigurator.Add("default", DefaultConfig)

	// globalConfigurator.LoadFile('main', "/etc/amyssh.yml")
	// globalConfig
}

func GetConfig() *Config {
	if globalCachedConfig == nil {
		globalCachedConfig = globalConfigurator.Merge([]string{"default", "main_file"}).(*Config)
	}
	return globalCachedConfig
}
