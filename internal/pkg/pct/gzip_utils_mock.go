package pct

import (
	"fmt"
	"path/filepath"
)

var mockGzipErrResponse error

type gzipHelpersImplMock struct{}

func (gzipHelpersImplMock) Gzip(source, target string) (gzipFilePath string, err error) {
	if mockGzipErrResponse != nil {
		return "", mockGzipErrResponse
	}

	if source == filepath.Join(testDir, "template.tar") {
		return filepath.Join(target, "template.tar.gz"), nil
	}
	return "", fmt.Errorf("Called with unexpected source: %v", source)
}
