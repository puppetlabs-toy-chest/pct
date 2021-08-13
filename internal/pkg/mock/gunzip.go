package mock

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

type Gunzip struct {
	Fs          afero.Fs
	ErrResponse bool
	Fail        bool
}

func (g *Gunzip) Gunzip(source, target string) (err error) {
	if g.ErrResponse {
		return fmt.Errorf("gunzip error")
	}

	// complete mocking of this function is not possible as the target it not always known

	// this code extracts a tar.gz, producing a tar witin target
	// using the mock fs, ensure that this exists
	// unless we want to test that NOT EXIST condition
	if !g.Fail {
		afs := &afero.Afero{Fs: g.Fs}
		tar := strings.TrimSuffix(filepath.Join(target, filepath.Base(source)), ".gz")
		afs.Create(tar) // nolint:errcheck  // #nosec // this result is not used in a secure application
	}

	return nil
}
