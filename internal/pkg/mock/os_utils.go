package mock

import (
	"io"
)

type OsUtil struct {
	WD string
}

func (*OsUtil) Hostname() (name string, err error) {
	return "fake.host", nil
}

func (*OsUtil) WriteString(w io.Writer, s string) (n int, err error) {
	return 0, nil
}

func (o *OsUtil) Getwd() (dir string, err error) {
	return o.WD, nil
}
