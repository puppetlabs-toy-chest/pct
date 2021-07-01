package tar

import (
	"path/filepath"
	"testing"
)

func TestTar(t *testing.T) {

	tempDir := t.TempDir()

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
			name: "Can TAR a template directory",
			args: args{
				source: "testdata/examples/full-project",
				target: tempDir,
			},
			wantErr:             false,
			expectedTarFilePath: filepath.Join(tempDir, "full-project.tar"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tarFilePath, err := Tar(tt.args.source, tt.args.target)

			if err != nil != tt.wantErr {
				t.Errorf("Tar() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tarFilePath != tt.expectedTarFilePath {
				t.Errorf("Tar() error = %v", err)
			}
		})
	}
}
