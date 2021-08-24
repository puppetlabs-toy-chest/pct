package pct_test

import (
	"fmt"
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
	templateContent  []byte
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

	templatePath := "path/to/somewhere"
	extractionPath := "path/to/extract/to"

	tarballBytes := []byte{
		0x1F, 0x8B, 0x08, 0x08, 0xF7, 0x5E, 0x14, 0x4A, 0x00, 0x03, 0x67, 0x6F,
		0x6F, 0x64, 0x2D, 0x70, 0x72, 0x6F, 0x6A, 0x65, 0x63, 0x74, 0x2E, 0x74,
		0x61, 0x72, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00,
	}

	fs := afero.NewMemMapFs()

	tests := []InstallTest{
		{
			name: "if it the file does not exist",
			args: args{
				templatePath: filepath.Join(templatePath, "good-project.tar.gz"),
				targetDir:    templatePath,
			},
			expectedError: fmt.Sprintf("No template package at %v", filepath.Join(templatePath, "good-project.tar.gz")),
		},
		{
			name: "should extract a tar.gz to a template folder",
			args: args{
				templatePath: filepath.Join(templatePath, "good-project.tar.gz"),
				targetDir:    extractionPath,
			},
			templateContent:  tarballBytes,
			unTarFile:        filepath.Join(extractionPath, "good-project.tar.gz"),
			untarFail:        false,
			expectedFilePath: filepath.Join(extractionPath, "good-project.tar.gz"),
			expectedError:    "false",
		},
		{
			name: "if it fails to gunzip",
			args: args{
				templatePath: filepath.Join(templatePath, "good-project.tar.gz"),
				targetDir:    templatePath,
			},
			gunzipErr:     true,
			expectedError: "gunzip error",
		},
		{
			name: "if it fails to untar",
			args: args{
				templatePath: filepath.Join(templatePath, "good-project.tar.gz"),
				targetDir:    templatePath,
			},
			untarFail:     true,
			expectedError: "untar error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			afs := &afero.Afero{Fs: fs}

			if len(tt.templateContent) > 0 {
				f, _ := afs.Create(tt.args.templatePath)
				f.Write(tarballBytes) // nolint:errcheck
			}

			installer := &pct.PctInstaller{
				&mock.Tar{ReturnedPath: tt.unTarFile, ErrResponse: tt.untarFail},
				&mock.Gunzip{Fs: fs, ErrResponse: tt.gunzipErr, Fail: tt.gunzipFail},
				afs,
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
