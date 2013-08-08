package amyssh

import "time"

var DefaultConfig = Config{
	DatabaseConfig{"localhost", 3306, "root", "", "amyssh"},
	[]UsersConfig{
		UsersConfig{"deployer", []string{"deploy", "admin"}, []string{}},
	},
	[]string{"default"},
	100 * time.Millisecond,
	10000 * time.Millisecond,
	100 * time.Millisecond, // Interval will be decreased if whole operation will take less than this
	200 * time.Millisecond, // Backoff threshold
	"authorized_keys2",
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
