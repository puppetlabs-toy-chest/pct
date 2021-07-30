package utils

import "io/ioutil"

type IoUtil interface {
	TempDir() (name string, err error)
}

type IoUtilHelpersImpl struct{}

func (IoUtilHelpersImpl) TempDir() (name string, err error) {
	return ioutil.TempDir("", "")
}
