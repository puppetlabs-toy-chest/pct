package pct

import (
	"fmt"
	"path/filepath"
)

var mockTarErrResponse error
var mockUntarErrResponse error

type tarHelpersImplMock struct{}
type untarHelpersImplMock struct{}

func (tarHelpersImplMock) Tar(source, target string) (tarFilePath string, err error) {
	if mockTarErrResponse != nil {
		return "", mockTarErrResponse
	}

	if source == testDir {
		return filepath.Join(testDir, "template.tar"), nil
	}
	return "", fmt.Errorf("Called with unexpected source: %v", source)
}

func (untarHelpersImplMock) Untar(source, target string) (tarFilePath string, err error) {
	if mockUntarErrResponse != nil {
		return "", mockUntarErrResponse
	}

	if source == filepath.Join(testDir, "good-project.tar") {
		return filepath.Join(testDir, "good-project"), nil
	}
	return "", fmt.Errorf("Called with unexpected source: %v", source)
}
