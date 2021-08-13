package utils

import (
	"io"
	"os"
)

type OsUtilI interface {
	Hostname() (name string, err error)
	WriteString(w io.Writer, s string) (n int, err error)
	Getwd() (dir string, err error)
}

type OsUtil struct{}

func (*OsUtil) Hostname() (name string, err error) {
	return os.Hostname()
}

func (*OsUtil) WriteString(w io.Writer, s string) (n int, err error) {
	return io.WriteString(w, s)
}

func (*OsUtil) Getwd() (dir string, err error) {
	return os.Getwd()
}
