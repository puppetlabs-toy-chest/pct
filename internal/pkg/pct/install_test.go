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
	name         string
	args         args
	expected     expected
	mocks        mocks
	mockReponses mockReponses
}

// what goes in
type args struct {
	templatePath string
	targetDir    string
	force        bool
}

// what comes out
type expected struct {
	errorMsg string
	filepath string
}

// filesystem mocks
type mocks struct {
	dirs  []string
	files map[string]string
}

// package mock responses
type mockReponses struct {
	untar  []mock.UntarResponse
	gunzip []mock.GunzipResponse
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

	tests := []InstallTest{
		{
			name: "if it the file does not exist",
			args: args{
				templatePath: filepath.Join(templatePath, "good-project.tar.gz"),
				targetDir:    templatePath,
			},
			expected: expected{
				errorMsg: fmt.Sprintf("No template package at %v", filepath.Join(templatePath, "good-project.tar.gz")),
			},
		},
		{
			name: "should extract a tar.gz to a template folder",
			args: args{
				templatePath: filepath.Join(templatePath, "good-project.tar.gz"),
				targetDir:    extractionPath,
			},
			expected: expected{
				filepath: filepath.Join(extractionPath, "puppetlabs", "good-project", "1.0.0"),
			},
			mockReponses: mockReponses{
				untar: []mock.UntarResponse{
					{
						ReturnPath:  filepath.Join(extractionPath, "good-project"),
						ErrResponse: false,
					},
					{
						ReturnPath:  filepath.Join(extractionPath, "puppetlabs", "good-project", "1.0.0"),
						ErrResponse: false,
					},
				},
				gunzip: []mock.GunzipResponse{
					{
						Fail:     false,
						FilePath: filepath.Join(extractionPath, "good-project.tar"),
					},
				},
			},
			mocks: mocks{
				dirs: []string{
					templatePath,
					extractionPath,
					filepath.Join(extractionPath, "good-project"),
				},
				files: map[string]string{
					filepath.Join(templatePath, "good-project.tar.gz"): string(tarballBytes),
					filepath.Join(extractionPath, "good-project", "pct-config.yml"): `---
template:
  id: good-project
  author: puppetlabs
  version: 1.0.0
`,
				},
			},
		},
		{
			name: "if it fails to gunzip",
			args: args{
				templatePath: filepath.Join(templatePath, "good-project.tar.gz"),
				targetDir:    templatePath,
			},
			mockReponses: mockReponses{
				gunzip: []mock.GunzipResponse{
					{
						ErrResponse: true,
					},
				},
			},
			mocks: mocks{
				dirs: []string{
					templatePath,
				},
				files: map[string]string{
					filepath.Join(templatePath, "good-project.tar.gz"): string(tarballBytes),
				},
			},
			expected: expected{
				errorMsg: "gunzip error",
			},
		},
		{
			name: "if it fails to untar",
			args: args{
				templatePath: filepath.Join(templatePath, "good-project.tar.gz"),
				targetDir:    templatePath,
			},
			mocks: mocks{
				dirs: []string{
					templatePath,
				},
				files: map[string]string{
					filepath.Join(templatePath, "good-project.tar.gz"): string(tarballBytes),
				},
			},
			mockReponses: mockReponses{
				gunzip: []mock.GunzipResponse{
					{
						Fail:     false,
						FilePath: filepath.Join(extractionPath, "good-project.tar"),
					},
				},
				untar: []mock.UntarResponse{
					{
						ErrResponse: true,
					},
				},
			},
			expected: expected{
				errorMsg: "untar error",
			},
		},
		{
			name: "If the tarball doesnt contain an valid config",
			args: args{
				templatePath: filepath.Join(templatePath, "good-project.tar.gz"),
				targetDir:    extractionPath,
			},
			expected: expected{
				errorMsg: "Invalid config: open " + filepath.FromSlash("path/to/extract/to/good-project/pct-config.yml") + ": file does not exist",
			},
			mockReponses: mockReponses{
				untar: []mock.UntarResponse{
					{
						ReturnPath:  filepath.Join(extractionPath, "good-project"),
						ErrResponse: false,
					},
					{
						ReturnPath:  filepath.Join(extractionPath, "puppetlabs", "good-project", "1.0.0"),
						ErrResponse: false,
					},
				},
				gunzip: []mock.GunzipResponse{
					{
						Fail:     false,
						FilePath: filepath.Join(extractionPath, "good-project.tar"),
					},
				},
			},
			mocks: mocks{
				dirs: []string{
					templatePath,
					extractionPath,
					filepath.Join(extractionPath, "good-project"),
				},
				files: map[string]string{
					filepath.Join(templatePath, "good-project.tar.gz"): string(tarballBytes),
				},
			},
		},
		{
			name: "when the template already exists",
			args: args{
				templatePath: filepath.Join(templatePath, "good-project.tar.gz"),
				targetDir:    extractionPath,
			},
			expected: expected{
				filepath: filepath.Join(extractionPath, "puppetlabs", "good-project", "1.0.0"),
			},
			mockReponses: mockReponses{
				untar: []mock.UntarResponse{
					{
						ReturnPath:  filepath.Join(extractionPath, "good-project"),
						ErrResponse: false,
					},
					{
						ReturnPath:  filepath.Join(extractionPath, "puppetlabs", "good-project", "1.0.0"),
						ErrResponse: false,
					},
				},
				gunzip: []mock.GunzipResponse{
					{
						Fail:     false,
						FilePath: filepath.Join(extractionPath, "good-project.tar"),
					},
				},
			},
			mocks: mocks{
				dirs: []string{
					templatePath,
					extractionPath,
					filepath.Join(extractionPath, "good-project"),
					filepath.Join(extractionPath, "puppetlabs", "good-project", "1.0.0"),
				},
				files: map[string]string{
					filepath.Join(templatePath, "good-project.tar.gz"): string(tarballBytes),
					filepath.Join(extractionPath, "good-project", "pct-config.yml"): `---
template:
  id: good-project
  author: puppetlabs
  version: 1.0.0
`,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fs := afero.NewMemMapFs()
			afs := &afero.Afero{Fs: fs}

			for _, path := range tt.mocks.dirs {
				afs.Mkdir(path, 0750) //nolint:gosec,errcheck // this result is not used in a secure application
			}

			for file, content := range tt.mocks.files {
				config, _ := afs.Create(file) //nolint:gosec,errcheck // this result is not used in a secure application
				config.Write([]byte(content)) //nolint:errcheck
			}

			installer := &pct.PctInstaller{
				&mock.Tar{UntarResponse: tt.mockReponses.untar},
				&mock.Gunzip{Fs: fs, GunzipResponse: tt.mockReponses.gunzip},
				afs,
				&afero.IOFS{Fs: fs},
			}

			returnedPath, err := installer.Install(tt.args.templatePath, tt.args.targetDir, tt.args.force)

			if tt.expected.errorMsg != "" && err != nil {
				assert.Contains(t, err.Error(), tt.expected.errorMsg)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected.filepath, returnedPath)
		})
	}
}
