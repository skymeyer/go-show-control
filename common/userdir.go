package common

import (
	"os"
	"path/filepath"
)

func MustUserConfigDir(name string) string {
	dir, err := GetUserConfigDir(name)
	if err != nil {
		panic("user config dir")
	}
	return dir
}

func MustUserCacheDir(name string) string {
	dir, err := GetUserCacheDir(name)
	if err != nil {
		panic("user cache dir")
	}
	return dir
}

func GetUserConfigDir(name string) (string, error) {
	return ensureUserDir(name, os.UserConfigDir)
}

func GetUserCacheDir(name string) (string, error) {
	return ensureUserDir(name, os.UserCacheDir)
}

type userDirFn func() (string, error)

func ensureUserDir(name string, userDir userDirFn) (string, error) {
	base, err := userDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(base, name)
	if err := os.MkdirAll(dir, os.FileMode(0700)); err != nil {
		return "", err
	}
	return dir, nil
}
