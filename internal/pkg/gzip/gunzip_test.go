package gzip

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGunzip(t *testing.T) {
	tarOutputPath := t.TempDir()
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
		t.Run(tt.name, func(t *testing.T) {
			if err := Gunzip(tt.args.source, tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("Gunzip() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, err := os.Stat(tt.expectedTarFilePath); err != nil {
				t.Errorf("Gunzip() did not output %v: %v", tt.expectedTarFilePath, err)
			}
		})
	}
}
