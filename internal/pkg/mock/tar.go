package mock

import (
	"fmt"
)

type Tar struct {
	ReturnedPath string
	ErrResponse  bool
}

func (t *Tar) Tar(source, target string) (tarFilePath string, err error) {

	if t.ErrResponse {
		return "", fmt.Errorf("tar error")
	}

	return t.ReturnedPath, nil
}

func (t *Tar) Untar(source, target string) (tarFilePath string, err error) {
	if t.ErrResponse {
		return "", fmt.Errorf("untar error")
	}

	return t.ReturnedPath, nil
}
