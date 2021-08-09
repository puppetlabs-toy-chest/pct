package pct

import (
	"fmt"
	"path/filepath"
)

var mockGzipErrResponse error
var mockGunzipErrResponse error

type gzipHelpersImplMock struct{}
type gunzipHelpersImplMock struct{}

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
