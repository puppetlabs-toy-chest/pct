package mock

import (
	"os/exec"
	"reflect"
)

type Exec struct {
	ExpectedName string
	ExpectedArg  []string

	ResponseCmd *exec.Cmd
}

func (e *Exec) Command(name string, arg ...string) *exec.Cmd {
	if name == e.ExpectedName && reflect.DeepEqual(arg, e.ExpectedArg) {
		return e.ResponseCmd
	}
	return nil
}
