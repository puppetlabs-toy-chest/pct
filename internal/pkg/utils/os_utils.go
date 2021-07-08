package utils

import (
	"os"
)

type OsUtil interface {
	Stat(string) (os.FileInfo, error)
	IsNotExist(err error) bool
}

type OsUtilHelpersImpl struct{}

func (OsUtilHelpersImpl) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (OsUtilHelpersImpl) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}
