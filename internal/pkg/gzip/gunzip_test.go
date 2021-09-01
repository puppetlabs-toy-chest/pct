package gzip_test

import (
	"path/filepath"
	"testing"

	"github.com/puppetlabs/pdkgo/internal/pkg/gzip"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
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
		mockFiles           map[string][]byte
		wantErr             string
		expectedTarFilePath string
	}{
		{
			name: "Can unzip a GZ to the specified location",
			args: args{
				source: "testdata/examples/good-project.tar.gz",
				target: tarOutputPath,
			},
			mockFiles: map[string][]byte{
				"testdata/examples/good-project.tar.gz": {
					0x1F, 0x8B, 0x08, 0x08, 0xF7, 0x5E, 0x14, 0x4A, 0x00, 0x03, 0x67, 0x6F,
					0x6F, 0x64, 0x2D, 0x70, 0x72, 0x6F, 0x6A, 0x65, 0x63, 0x74, 0x2E, 0x74,
					0x61, 0x72, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00,
				},
			},
			expectedTarFilePath: filepath.Join(tarOutputPath, "good-project.tar"),
		},
		{
			name: "Fails when the tar.gz does not exist",
			args: args{
				source: "testdata/examples/noexist-project.tar.gz",
				target: tarOutputPath,
			},
			wantErr: "open " + filepath.FromSlash("testdata/examples/noexist-project.tar.gz") + ": file does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, true)

			afs := &afero.Afero{Fs: fs}
			for file, content := range tt.mockFiles {
				config, _ := afs.Create(file) //nolint:gosec,errcheck // this result is not used in a secure application
				config.Write([]byte(content)) //nolint:errcheck
			}

			g := &gzip.Gunzip{AFS: afs}
			file, err := g.Gunzip(tt.args.source, tt.args.target)

			if tt.wantErr != "" && assert.Error(t, err) {
				assert.Equal(t, tt.wantErr, err.Error())
			}

			assert.Equal(t, tt.expectedTarFilePath, file)

			_, err = afs.Stat(tt.expectedTarFilePath)
			assert.NoError(t, err)

		})
	}

}
