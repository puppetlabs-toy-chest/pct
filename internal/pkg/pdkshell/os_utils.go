package pdkshell

import "os"

type osHelpers interface {
	Environ() []string
}

type osHelpersImpl struct{}

func (osHelpersImpl) Environ() []string {
	return os.Environ()
}
