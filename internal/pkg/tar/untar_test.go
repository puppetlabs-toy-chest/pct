package tar

import (
	"os"
	"path/filepath"
	"testing"
)

func TestUntar(t *testing.T) {

	tempDir := t.TempDir()

	type args struct {
		tarball string
		target  string
	}

	tests := []struct {
		name           string
		args           args
		wantErr        bool
		tarContentTree []string
	}{
		{
			name: "Can unTAR a template",
			args: args{
				tarball: "testdata/examples/good-project.tar",
				target:  tempDir,
			},
			wantErr: false,
			tarContentTree: []string{
				"content/empty.txt",
				"content/goodfile.txt.tmpl",
				"pct-config.yml",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := Untar(tt.args.tarball, tt.args.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("Untar() error = %v, wantErr %v", err, tt.wantErr)
			}

			for _, f := range tt.tarContentTree {
				filePath := filepath.Join(path, f)
				_, err := os.Stat(filePath)
				if err != nil {
					t.Errorf("Did not find expected file: %v", filePath)
				}
			}
		})
	}
}
