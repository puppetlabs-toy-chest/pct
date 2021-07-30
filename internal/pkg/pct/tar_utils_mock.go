package pct

import (
	"fmt"
	"path/filepath"
)

var mockTarErrResponse error

type tarHelpersImplMock struct{}

func (tarHelpersImplMock) Tar(source, target string) (tarFilePath string, err error) {
	if mockTarErrResponse != nil {
		return "", mockTarErrResponse
	}

	if source == testDir {
		return filepath.Join(testDir, "template.tar"), nil
	}
	return "", fmt.Errorf("Called with unexpected source: %v", source)
}
