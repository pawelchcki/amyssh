package amyssh

import (
	"fmt"
	// "math/rand"
	// "time"
	// "os/user"
)

var _ = fmt.Println

func tagsFromSet(tagSet map[string]struct{}) []string {
	tags := make([]string, 0, len(tagSet))
	for k, _ := range tagSet {
		tags = append(tags, k)
	}
	return tags
}

func processUsers(cfg *Config) ([]string, map[string]*UsersConfig) {
	tagSet := make(map[string]struct{})
	userNameMap := make(map[string]*UsersConfig)
	for _, user := range cfg.Users {
		for _, tag := range user.Tags {
			tagSet[tag] = struct{}{}
		}
		userNameMap[user.Name] = &user
	}

	return tagsFromSet(tagSet), userNameMap
}

func processHostTags(cfg *Config) []string {
	tagSet := make(map[string]struct{})

	for _, tag := range cfg.HostTags {
		tagSet[tag] = struct{}{}
	}

	return tagsFromSet(tagSet)
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
