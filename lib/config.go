package amyssh

import "time"

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
	MinPollInterval:        100 * time.Millisecond,
	MaxPollInterval:        10000 * time.Millisecond,
	PerformanceThreshold:   100 * time.Millisecond, // Interval will be decreased if whole operation will take less than this
	BackoffThreshold:       200 * time.Millisecond, // Backoff threshold
	AuthorizedKeysFileName: "authorized_keys2",
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
}
