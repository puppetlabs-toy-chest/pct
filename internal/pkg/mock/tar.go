package mock

import (
	"fmt"
)

type UntarResponse struct {
	ReturnPath  string
	ErrResponse bool
}

type Tar struct {
	ReturnedPath  string
	ErrResponse   bool
	UntarResponse []UntarResponse
	unTarCalled   int
}

func (t *Tar) Tar(source, target string) (tarFilePath string, err error) {

	if t.ErrResponse {
		return "", fmt.Errorf("tar error")
	}

	return t.ReturnedPath, nil
}

func (t *Tar) Untar(source, target string) (tarFilePath string, err error) {
	if t.UntarResponse[t.unTarCalled].ErrResponse {
		t.unTarCalled++
		return "", fmt.Errorf("untar error")
	}

	path := t.UntarResponse[t.unTarCalled].ReturnPath
	t.unTarCalled++

	return path, nil
}
