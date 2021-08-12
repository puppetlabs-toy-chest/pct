package gzip_test

import (
	"path/filepath"
	"testing"

	"github.com/puppetlabs/pdkgo/internal/pkg/gzip"
	"github.com/spf13/afero"
)

func TestGunzip(t *testing.T) {
	tarOutputPath := t.TempDir()
	fs := afero.NewMemMapFs()

	type args struct {
		source string
		target string
	}
	tests := []struct {
		name                string
		args                args
		wantErr             bool
		expectedTarFilePath string
	}{
		{
			name: "Can unzip a GZ to the specified location",
			args: args{
				source: "testdata/examples/good-project.tar.gz",
				target: tarOutputPath,
			},
			wantErr:             false,
			expectedTarFilePath: filepath.Join(tarOutputPath, "good-project.tar"),
		},
	}
	for _, tt := range tests {

		// mock gzip
		afs := &afero.Afero{Fs: fs}
		f, _ := afs.Create(tt.args.source)
		tarballBytes := []byte{
			0x1F, 0x8B, 0x08, 0x08, 0xF7, 0x5E, 0x14, 0x4A, 0x00, 0x03, 0x67, 0x6F,
			0x6F, 0x64, 0x2D, 0x70, 0x72, 0x6F, 0x6A, 0x65, 0x63, 0x74, 0x2E, 0x74,
			0x61, 0x72, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00,
		}
		f.Write(tarballBytes) // nolint:errcheck

		t.Run(tt.name, func(t *testing.T) {
			g := &gzip.Gunzip{
				afs,
			}
			err := g.Gunzip(tt.args.source, tt.args.target)

			if (err != nil) != tt.wantErr {
				t.Errorf("Gunzip() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, err := afs.Stat(tt.expectedTarFilePath); err != nil {
				t.Errorf("Gunzip() did not output %v: %v", tt.expectedTarFilePath, err)
			}
		})
	}
}
