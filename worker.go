package amyssh

type processedData struct {
	uniqueUserTags []string
	userTagMap     map[string]*UsersConfig
	uniqueHostTags []string
}

func processUsers(cfg *Config) ([]string, map[string]*UsersConfig) {
	tagSet := make(StringSet)
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
var globalData *processedData

func Perform(cfg *Config) error {
	if globalData == nil {
		globalData = &processedData{}
		globalData.uniqueUserTags, globalData.userTagMap = processUsers(cfg)
		globalData.uniqueHostTags = processHostTags(cfg)
	}
	var err error
	if globalConnection == nil {
		globalConnection, err = NewCon(cfg)
	}
	if err != nil {
		return err
	}
	keyMap, err := globalConnection.FetchKeys(globalData.uniqueHostTags,
		globalData.uniqueUserTags)
	if err != nil {
		return err
	}

	return ProcessKeys(cfg, keyMap, globalData.userTagMap)
}
