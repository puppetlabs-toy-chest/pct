package gzip

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGzip(t *testing.T) {
	gzipOutputPath := t.TempDir()

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
			_, err := Gzip(tt.args.source, tt.args.target)

			if err != nil && !tt.wantErr {
				t.Errorf("Gzip() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, err := os.Stat(tt.expectedGzipFilePath); err != nil {
				t.Errorf("Gzip() did not create %v: %v", tt.expectedGzipFilePath, err)
			}
		})
	}
}
