package gzip_test

import (
	"path/filepath"
	"testing"

	"github.com/puppetlabs/pct/pkg/gzip"
	"github.com/spf13/afero"
)

func TestGzip(t *testing.T) {
	fs := afero.NewMemMapFs()
	afs := &afero.Afero{Fs: fs}

	gzipOutputPath, _ := afs.TempDir("", "")

	type args struct {
		source string
		target string
	}
	tests := []struct {
		name                 string
		args                 args
		wantErr              bool
		expectedGzipFilePath string
	}{
		{
			name: "Should create GZIP from TAR",
			args: args{
				source: "testdata/examples/good-project.tar",
				target: gzipOutputPath,
			},
			wantErr:              false,
			expectedGzipFilePath: filepath.Join(gzipOutputPath, "good-project.tar.gz"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// mock the tarball
			f, _ := afs.Create(tt.args.source)
			f.WriteString("tar contents") //nolint:errcheck

			// Initialize gzip with our mock filesystem
			g := &gzip.Gzip{
				&afero.Afero{Fs: fs},
			}
			_, err := g.Gzip(tt.args.source, tt.args.target)

			if err != nil && !tt.wantErr {
				t.Errorf("Gzip() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, err := afs.Stat(tt.expectedGzipFilePath); err != nil {
				t.Errorf("Gzip() did not create %v: %v", tt.expectedGzipFilePath, err)
			}
		})
	}
}
