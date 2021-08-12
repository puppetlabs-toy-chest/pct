package pct_test

import (
	"path/filepath"
	"testing"

	"github.com/puppetlabs/pdkgo/internal/pkg/mock"
	"github.com/puppetlabs/pdkgo/internal/pkg/pct"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

type InstallTest struct {
	name             string
	args             args
	expectedError    string
	expectedFilePath string
	unTarFile        string
	untarFail        bool
	gunzipErr        bool
	gunzipFail       bool
}
type args struct {
	templatePath string
	targetDir    string
}

func TestInstall(t *testing.T) {

	testDir := "/path/to/somewhere"
	fs := afero.NewMemMapFs()

	tests := []InstallTest{
		{
			name: "should extract a tar.gz to a template folder",
			args: args{
				templatePath: filepath.Join(testDir, "good-project.tar.gz"),
				targetDir:    testDir,
			},
			unTarFile:        filepath.Join(testDir, "good-project"),
			untarFail:        false,
			expectedFilePath: filepath.Join(testDir, "good-project"),
			expectedError:    "false",
		},
		{
			name: "if it fails to gunzip",
			args: args{
				templatePath: filepath.Join(testDir, "good-project.tar.gz"),
				targetDir:    testDir,
			},
			gunzipErr:     true,
			expectedError: "gunzip error",
		},
		{
			name: "if it the file isnt where it should be",
			args: args{
				templatePath: filepath.Join(testDir, "good-project.tar.gz"),
				targetDir:    testDir,
			},
			gunzipErr:     false,
			gunzipFail:    true,
			expectedError: "file does not exist",
		},
		{
			name: "if it fails to untar",
			args: args{
				templatePath: filepath.Join(testDir, "good-project.tar.gz"),
				targetDir:    testDir,
			},
			untarFail:     true,
			expectedError: "untar error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			installer := &pct.PctInstaller{
				&mock.Tar{ReturnedPath: tt.unTarFile, ErrResponse: tt.untarFail},
				&mock.Gunzip{Fs: fs, ErrResponse: tt.gunzipErr, Fail: tt.gunzipFail},
				&afero.Afero{Fs: fs},
				&afero.IOFS{Fs: fs},
			}
			got, err := installer.Install(tt.args.templatePath, tt.args.targetDir)
			if err != nil {
				assert.Contains(t, err.Error(), tt.expectedError)
			}
			assert.Equal(t, tt.expectedFilePath, got)
		})
	}
}
