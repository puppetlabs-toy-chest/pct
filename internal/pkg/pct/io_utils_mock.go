package pct

import "path/filepath"

var testDir string = filepath.Clean("/path/to/nowhere")

type ioUtilHelpersImplMock struct{}

func (ioUtilHelpersImplMock) TempDir() (name string, err error) {
	return testDir, nil
}
