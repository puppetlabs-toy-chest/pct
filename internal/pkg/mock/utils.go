package mock

import (
	"path/filepath"
)

type UtilsHelper struct {
	TestDir             string
	IsModuleRootErrResp error
	Home                string
	ReaderError         bool
}

func (u *UtilsHelper) IsModuleRoot() (string, error) {
	if u.IsModuleRootErrResp != nil {
		return "", u.IsModuleRootErrResp
	}
	return filepath.Clean(u.TestDir), nil
}

func (u *UtilsHelper) Dir() (string, error) {
	return u.Home, nil
}
