package pct

import "path/filepath"

type utilsHelperImplMock struct{}

var mockIsModuleRootErrResp error

func (utilsHelperImplMock) IsModuleRoot() (string, error) {
	if mockIsModuleRootErrResp != nil {
		return "", mockIsModuleRootErrResp
	}
	return filepath.Clean(testDir), nil
}
