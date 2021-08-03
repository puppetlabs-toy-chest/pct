package pct

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

var testDir string = filepath.Clean("/path/to/nowhere")

var fileInfoMock fileInfo
var mockIsModuleRootErrResp error
var mockIsNotExistResponse bool
var mockGzipErrResponse error
var mockGunzipErrResponse error
var mockStatResponses []mockStatResponse
var mockTarErrResponse error
var mockUntarErrResponse error

type fileInfo = fs.FileInfo
type gzipHelpersImplMock struct{}
type gunzipHelpersImplMock struct{}
type ioUtilHelpersImplMock struct{}
type mockStatResponse struct {
	expectedName string
	mockError    error
}
type osUtilHelpersImplMock struct{}
type tarHelpersImplMock struct{}
type untarHelpersImplMock struct{}
type utilsHelperImplMock struct{}

func (utilsHelperImplMock) IsModuleRoot() (string, error) {
	if mockIsModuleRootErrResp != nil {
		return "", mockIsModuleRootErrResp
	}
	return filepath.Clean(testDir), nil
}

func (osUtilHelpersImplMock) IsNotExist(err error) bool {
	return mockIsNotExistResponse
}

func (gzipHelpersImplMock) Gzip(source, target string) (gzipFilePath string, err error) {
	if mockGzipErrResponse != nil {
		return "", mockGzipErrResponse
	}

	if source == filepath.Join(testDir, "template.tar") {
		return filepath.Join(target, "template.tar.gz"), nil
	}
	return "", fmt.Errorf("Called with unexpected source: %v", source)
}

func (gunzipHelpersImplMock) Gunzip(source, target string) (err error) {
	if mockGunzipErrResponse != nil {
		return mockGunzipErrResponse
	}

	if source == filepath.Join(testDir, "good-project.tar.gz") {
		return nil
	}
	return fmt.Errorf("Called with unexpected source: %v", source)
}

func (osUtilHelpersImplMock) Stat(name string) (os.FileInfo, error) {
	mockIsNotExistResponse = false
	for _, mockStatResponse := range mockStatResponses {
		if name == mockStatResponse.expectedName {
			if mockStatResponse.mockError != nil {
				mockIsNotExistResponse = true
			}
			return fileInfoMock, mockStatResponse.mockError
		}
	}
	return nil, nil
}

func (tarHelpersImplMock) Tar(source, target string) (tarFilePath string, err error) {
	if mockTarErrResponse != nil {
		return "", mockTarErrResponse
	}

	if source == testDir {
		return filepath.Join(testDir, "template.tar"), nil
	}
	return "", fmt.Errorf("Called with unexpected source: %v", source)
}
func (ioUtilHelpersImplMock) TempDir() (name string, err error) {
	return testDir, nil
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
