package amyssh

import (
	"fmt"
	// "math/rand"
	"io/ioutil"
	"os"
	os_user "os/user"
	"path/filepath"
	"strconv"
	"time"
)

func ensureSshDirExists(user *os_user.User) error {
	err := os.Chdir(user.HomeDir)
	if err != nil {
		return err
	}
	sshDir := ".ssh"
	fi, err := os.Stat(sshDir)
	if err != nil && os.IsNotExist(err) {
		err = os.Mkdir(sshDir, 0700)
		if err != nil {
			return err
		}
		err = chown(sshDir, user)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return fmt.Errorf(".ssh is not a directory")
	}
	return nil
}

func chown(path string, user *os_user.User) error {
	uid, gid, err := convertUidGid(user)
	if err != nil {
		return err
	}
	err = os.Chown(path, uid, gid)
	if err != nil && false { //todo fatality
		return err
	}
	return nil
}

func convertUidGid(user *os_user.User) (uid, gid int, err error) {
	uid0, err := strconv.ParseInt(user.Uid, 10, 0)
	uid = int(uid0)
	if err != nil {
		return
	}
	gid0, err := strconv.ParseInt(user.Uid, 10, 0)
	gid = int(gid0)
	return
}

func generateKeySet(userData *UsersConfig, keysMap map[string][]string) map[string]struct{} {
	keySet := SetFromList(nil, userData.Keys)
	for _, userTag := range userData.Tags {
		keys := keysMap[userTag]
		if keys != nil {
			keySet = SetFromList(keySet, keys)
		}
	}
	return keySet
}

func writeTempKey(userName string, keySet map[string]struct{}) (*os.File, error) {
	f, err := ioutil.TempFile("", fmt.Sprintf("amyssh-%s.", userName))
	for key, _ := range keySet {
		_, err := fmt.Fprintln(f, key)
		if err != nil {
			f.Close()
			os.Remove(f.Name())
			return nil, err
		}
	}
	err = f.Close()
	if err != nil {
		return nil, err
	}

	return f, nil
}

func fileDataChanged(filePath string, keySet map[string]struct{}) bool {
	f, err := os.Open(filePath)
	if err != nil {
		return true
	}
	defer f.Close()

	line := ""
	visitedKeys := make(map[string]struct{})
	numOfKeys := len(keySet)

	for _, err := fmt.Fscanln(f, &line); err == nil; _, err = fmt.Fscanln(f, &line) {
		_, hasKey := keySet[line]
		if hasKey {
			_, hasKey := visitedKeys[line]
			if hasKey {
				return true
			}
			visitedKeys[line] = struct{}{}
		} else {
			return true
		}
		if len(visitedKeys) > numOfKeys {
			return true
		}
	}
	if len(visitedKeys) != numOfKeys {
		return true
	}

	return false
}

func backupAndSubstitute(keyFileName, tmpFileName string) error {
	_, err := os.Stat(keyFileName)
	if !os.IsNotExist(err) {
		//TODO: copy instead of rename to be extra safe
		err := os.Rename(keyFileName, fmt.Sprintf("%s-%s", keyFileName, time.Now().Format("060102150405.000")))
		if err != nil {
			os.Remove(tmpFileName)
		}
	}
	err = os.Rename(tmpFileName, keyFileName)
	if err != nil {
		os.Remove(tmpFileName)
		return err
	}
	return nil
}

func (cfg *Config) processKey(userName string, keysMap map[string][]string, userData *UsersConfig) error {
	user, err := os_user.Lookup(userName)
	if err != nil {
		return err
	}

	keySet := generateKeySet(userData, keysMap)
	authorizedKeysFilepath := filepath.Join(user.HomeDir, ".ssh", cfg.AuthorizedKeysFileName)
	if !fileDataChanged(authorizedKeysFilepath, keySet) {
		return nil //file doesn't need changes skip
	}
	f, err := writeTempKey(userData.Name, keySet)
	if err != nil {
		return err
	}

	err = chown(f.Name(), user)
	if err != nil {
		return err
	}

	err = ensureSshDirExists(user)
	if err != nil {
		os.Remove(f.Name())
		return err
	}

	err = backupAndSubstitute(authorizedKeysFilepath, f.Name())
	if err != nil {
		os.Remove(f.Name())
		return err
	}

	return nil
}

func ProcessKeys(cfg *Config, keysMap map[string][]string, userMap map[string]*UsersConfig) error {
	var lastError error
	for user, userData := range userMap {
		lastError = cfg.processKey(user, keysMap, userData)
	}
	return lastError
}
