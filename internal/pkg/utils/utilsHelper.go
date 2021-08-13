package utils

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

type UtilsHelperI interface {
	IsModuleRoot() (string, error)
	Dir() (string, error)
}

type UtilsHelper struct{}

// Check if we're currently in the module root dir.
// Return the sanitized file path if we are in a module root, otherwise an empty string.
func (u *UtilsHelper) IsModuleRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	_, err = os.Stat(filepath.Join(wd, "metadata.json"))
	if err != nil {
		return "", err
	}

	return filepath.Clean(wd), nil
}

func (u *UtilsHelper) Dir() (string, error) {
	return homedir.Dir()
}
