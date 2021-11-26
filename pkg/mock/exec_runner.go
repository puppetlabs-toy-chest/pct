package mock

import (
	"fmt"
	"reflect"
)

type Exec struct {
	ExpectedName string
	ExpectedArg  []string

	ResponseMsg string
	ResponseError bool
}

func (e *Exec) Command(name string, arg ...string) error {
	if name == e.ExpectedName && reflect.DeepEqual(arg, e.ExpectedArg) {
		return nil
	}
	return fmt.Errorf("Unexpected Command %v %v", name, arg)
}

func (e *Exec) Output() ([]byte, error) {
	if e.ResponseError {
		return nil, fmt.Errorf(e.ResponseMsg)
 	}
	return []byte(e.ResponseMsg), nil
}
