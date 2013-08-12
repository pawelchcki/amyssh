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

func generateKeySet(userData *UsersConfig, keysMap map[string]StringSet) StringSet {
	keySet := NewSetFromList(userData.Keys)
	for _, userTag := range userData.Tags {
		keys := keysMap[userTag]
		if keys != nil {
			keySet.Union(keys)
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

type fileInfoCacheEntry struct {
	fileKeySet StringSet
	fileInfo   os.FileInfo
}

type fileInfoCacheType map[string]*fileInfoCacheEntry

var fileInfoCache fileInfoCacheType

func init() {
	fileInfoCache = make(fileInfoCacheType)
}

type duplicatedEntryError struct {
	entry string
}

func (e *duplicatedEntryError) Error() string {
	return "file had duplicated entry of: " + e.entry
}

func fileToKeySet(filePath string) (StringSet, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	keySet := make(StringSet)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if _, hasKey := keySet[line]; hasKey {
			return nil, &duplicatedEntryError{line}
		}
		keySet[line] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return keySet, nil
}

func (c fileInfoCacheType) isFileChanged(filePath string) (bool, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil // file doesn't exist we need update
		} else {
			return false, err
		}
	}
	if c[filePath] == nil {
		fileKeySet, err := fileToKeySet(filePath)
		if err != nil {
			if _, ok := err.(*duplicatedEntryError); ok {
				return true, nil
			} else {
				return false, err
			}
		}
		c[filePath] = &fileInfoCacheEntry{
			fileKeySet: fileKeySet,
			fileInfo:   fileInfo,
		}
	} else {
		if c[filePath].fileInfo.ModTime() != fileInfo.ModTime() ||
			c[filePath].fileInfo.Size() != fileInfo.Size() {
			return true, nil
		}
	}
	return false, nil
}
func (c fileInfoCacheType) isKeySetEqual(filePath string, keySet StringSet) (bool, error) {
	fileInfo := c[filePath]
	if fileInfo == nil {
		return false, nil
	}
	if len(fileInfo.fileKeySet) != len(keySet) {
		return false, nil
	}
	for k, _ := range keySet {
		_, hasKey := fileInfo.fileKeySet[k]
		if !hasKey {
			return false, nil
		}
	}
	return true, nil
}
func (c fileInfoCacheType) updateFileInfo(filePath string, keySet StringSet) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	if c[filePath] == nil {
		c[filePath] = &fileInfoCacheEntry{
			fileInfo:   fileInfo,
			fileKeySet: keySet,
		}
	}
	return nil
}

func fileNeedsUpdate(filePath string, keySet StringSet) (bool, error) {
	fileChanged, err := fileInfoCache.isFileChanged(filePath)
	if err != nil {
		return false, err
	}
	if fileChanged {
		return true, nil
	}
	keySetEqual, err := fileInfoCache.isKeySetEqual(filePath, keySet)
	if err != nil {
		return false, err
	}

	return !keySetEqual, nil
}

func processKey(cfg *Config, userName string, keysMap map[string]StringSet, userData *UsersConfig) error {
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

	err = fileInfoCache.updateFileInfo(authorizedKeysFilepath, keySet)
	if err != nil {
		return err
	}
	return nil
}

func ProcessKeys(cfg *Config, keysMap map[string]StringSet, userMap map[string]*UsersConfig) error {
	var lastError error
	for user, userData := range userMap {
		lastError = processKey(cfg, user, keysMap, userData)
	}
	return lastError
}
