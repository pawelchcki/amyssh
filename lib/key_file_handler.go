package amyssh

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
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
	if err != nil {
		return err
	}
	return nil
}

func convertUidGid(user *os_user.User) (uid, gid int, err error) {
	uid64, err := strconv.ParseInt(user.Uid, 10, 0)
	uid = int(uid64)
	if err != nil {
		return
	}
	gid64, err := strconv.ParseInt(user.Gid, 10, 0)
	gid = int(gid64)
	return
}

func generateKeySet(userData *UsersConfig, keysMap map[string][]string) StringSet {
	keySet := SetFromList(nil, userData.Keys)
	for _, userTag := range userData.Tags {
		keys := keysMap[userTag]
		if keys != nil {
			keySet = SetFromList(keySet, keys)
		}
	}
	return keySet
}

func writeTempKey(userName string, keySet StringSet) (*os.File, error) {
	f, err := ioutil.TempFile("", fmt.Sprintf("amyssh-%s.", userName))
	if err != nil {
		return nil, err
	}
	log.Printf("writing %d unique keys", len(keySet))
	defer func() {
		f.Close()
		if err != nil {
			os.Remove(f.Name())
		}
	}()

	for key, _ := range keySet {
		_, err := fmt.Fprintln(f, key)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	return f, nil
}

func fileNeedsUpdate(filePath string, keySet StringSet) (bool, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil // file doesn't exist we need update
		} else {
			return false, err
		}
	}
	f, err := os.Open(filePath)
	if err != nil {
		return false, err
	}

	defer f.Close()
	visitedKeys := make(StringSet)
	numOfKeys := len(keySet)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		_, hasKey := keySet[line]
		if hasKey {
			_, hasKey := visitedKeys[line]
			if hasKey {
				return true, nil
			}
			visitedKeys[line] = struct{}{}
		} else {
			return true, nil
		}
		if len(visitedKeys) > numOfKeys {
			return true, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}
	if len(visitedKeys) != numOfKeys {
		return true, nil
	}

	return false, nil
}

func backupAndSubstitute(keyFileName, tmpFileName string) error {
	_, err := os.Stat(keyFileName)
	if !os.IsNotExist(err) {
		backupFileName := fmt.Sprintf("%s-%s", keyFileName, time.Now().Format("060102150405.000"))
		log.Printf("saving backup file: %s", backupFileName)
		//TODO: copy instead of rename to be extra safe
		err = os.Rename(keyFileName, backupFileName)
		if err != nil {
			os.Remove(tmpFileName)
		}
	}
	err = os.Rename(tmpFileName, keyFileName)
	if err != nil {
		os.Remove(tmpFileName)
		return err
	}
	log.Printf("saved keys to: %s", keyFileName)
	return nil
}

func processKey(cfg *Config, userName string, keysMap map[string][]string, userData *UsersConfig) error {
	user, err := os_user.Lookup(userName)
	if err != nil {
		return err
	}

	keySet := generateKeySet(userData, keysMap)

	authorizedKeysFilepath := filepath.Join(user.HomeDir, ".ssh", cfg.AuthorizedKeysFileName)
	update, err := fileNeedsUpdate(authorizedKeysFilepath, keySet)
	if !update {
		return nil //file doesn't need changes skip
	}
	f, err := writeTempKey(userData.Name, keySet)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			f.Close()
			os.Remove(f.Name())
		}
	}()

	err = chown(f.Name(), user)
	if err != nil {
		return err
	}

	err = ensureSshDirExists(user)
	if err != nil {
		return err
	}

	err = backupAndSubstitute(authorizedKeysFilepath, f.Name())
	if err != nil {
		return err
	}

	return nil
}

func ProcessKeys(cfg *Config, keysMap map[string][]string, userMap map[string]*UsersConfig) error {
	var lastError error
	for user, userData := range userMap {
		lastError = processKey(cfg, user, keysMap, userData)
	}
	return lastError
}
