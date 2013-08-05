package amyssh

var DefaultConfig = Config{
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
type Config struct {
	Database DatabaseConfig
	Users    []UsersConfig
	HostTags []string
}
