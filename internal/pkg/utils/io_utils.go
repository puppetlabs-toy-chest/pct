package utils

import "io/ioutil"

type IoUtilI interface {
	TempDir() (name string, err error)
}

type IoUtil struct{}

func (IoUtil) TempDir() (name string, err error) {
	return ioutil.TempDir("", "")
}
