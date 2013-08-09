package amyssh

import ()

func processUsers(cfg *Config) ([]string, map[string]*UsersConfig) {
	tagSet := make(map[string]struct{})
	userNameMap := make(map[string]*UsersConfig)
	for _, user := range cfg.Users {
		tagSet = SetFromList(tagSet, user.Tags)
		userNameMap[user.Name] = &user
	}

	return StringsFromSet(tagSet), userNameMap
}

func processHostTags(cfg *Config) []string {
	tagSet := SetFromList(nil, cfg.HostTags)
	return StringsFromSet(tagSet)
}

var globalConnection *Connection

func Perform(cfg *Config) error {
	userTags, userMap := processUsers(cfg)
	hostTags := processHostTags(cfg)
	var err error
	if globalConnection == nil {
		globalConnection, err = NewCon(cfg)
	}
	if err != nil {
		return err
	}
	keyMap, err := globalConnection.FetchKeys(hostTags, userTags)
	if err != nil {
		return err
	}

	return ProcessKeys(cfg, keyMap, userMap)
}
