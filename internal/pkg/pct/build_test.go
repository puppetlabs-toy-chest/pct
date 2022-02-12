package pct_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/puppetlabs/pdkgo/internal/pkg/pct"
	"github.com/puppetlabs/pdkgo/internal/pkg/pct_config_processor"
	"github.com/puppetlabs/pdkgo/pkg/mock"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {

	type args struct {
		templatePath string
		targetDir    string
	}

	var mockTemplateDir = "/path/to/my/cool-template"
	var mockConfigFilePath = filepath.Clean(filepath.Join(mockTemplateDir, "pct-config.yml"))

	tests := []struct {
		name                    string
		args                    args
		mockIsModuleRootErrResp error
		mockDirs                []string
		mockFiles               map[string]string
		expectedFilePath        string
		tarFile                 string
		gzipFile                string
		expectedErr             string
		mockTarErr              bool
		mockGzipErr             bool
		testTempDir             string
	}{
		{
			name: "Should return err if template path does not exist",
			args: args{
				templatePath: mockTemplateDir,
				targetDir:    mockTemplateDir,
			},
			expectedFilePath: "",
			expectedErr:      "No template directory at /path/to/my/cool-template",
		},
		{
			name: "Should return err if template path does not contain pct-config.yml",
			args: args{
				templatePath: mockTemplateDir,
				targetDir:    mockTemplateDir,
			},
			mockDirs: []string{
				mockTemplateDir,
			},
			expectedFilePath: "",
			expectedErr:      "No 'pct-config.yml' found in /path/to/my/cool-template",
		},
		{
			name: "Should return err if content dir does not exist",
			args: args{
				templatePath: mockTemplateDir,
				targetDir:    mockTemplateDir,
			},
			mockDirs: []string{
				mockTemplateDir,
			},
			mockFiles: map[string]string{
				filepath.Clean(filepath.Join(mockTemplateDir, "pct-config.yml")): `---
template:
  id: builder
  author: puppetlabs
  version: 1.0.0
`,
			},
			expectedFilePath: "",
			expectedErr:      "No 'content' dir found in /path/to/my/cool-template",
		},
		{
			name: "Should not attempt to GZIP when TAR operation fails",
			args: args{
				templatePath: mockTemplateDir,
				targetDir:    mockTemplateDir,
			},
			mockDirs: []string{
				mockTemplateDir,
				filepath.Join(mockTemplateDir, "content"),
			},
			mockFiles: map[string]string{
				filepath.Clean(filepath.Join(mockTemplateDir, "pct-config.yml")): `---
template:
  id: builder
  author: puppetlabs
  version: 1.0.0
`,
			},
			expectedFilePath: "",
			expectedErr:      "tar error",
			mockTarErr:       true,
		},
		{
			name: "Should return error and empty path if GZIP operation fails",
			args: args{
				templatePath: mockTemplateDir,
				targetDir:    mockTemplateDir,
			},
			mockFiles: map[string]string{
				filepath.Clean(filepath.Join(mockTemplateDir, "pct-config.yml")): `---
template:
  id: builder
  author: puppetlabs
  version: 1.0.0
`,
			},
			mockDirs: []string{
				mockTemplateDir,
				filepath.Join(mockTemplateDir, "content"),
			},
			tarFile:          "/path/to/nowhere/pkg/nowhere.tar",
			expectedFilePath: "",
			expectedErr:      "gzip error",
			mockTarErr:       false,
			mockGzipErr:      true,
		},
		{
			name: "Should TAR.GZ valid template to $MODULE_ROOT/pkg and return path",
			args: args{
				templatePath: mockTemplateDir,
				targetDir:    mockTemplateDir,
			},
			mockDirs: []string{
				mockTemplateDir,
				filepath.Join(mockTemplateDir, "content"),
			},
			mockFiles: map[string]string{
				filepath.Clean(filepath.Join(mockTemplateDir, "pct-config.yml")): `---
template:
  id: builder
  author: puppetlabs
  version: 1.0.0
`,
			},
			tarFile:          "/path/to/nowhere/pkg/nowhere.tar",
			gzipFile:         "/path/to/nowhere/pkg/nowhere.tar.gz",
			expectedFilePath: "/path/to/nowhere/pkg/nowhere.tar.gz",
			mockTarErr:       false,
		},
		{
			name: "Should complain that `id` is missing from pct-config.yml",
			args: args{
				templatePath: mockTemplateDir,
				targetDir:    mockTemplateDir,
			},
			mockDirs: []string{
				mockTemplateDir,
				filepath.Join(mockTemplateDir, "content"),
			},
			mockFiles: map[string]string{
				filepath.Clean(filepath.Join(mockTemplateDir, "pct-config.yml")): `---
template:
  author: puppetlabs
  version: 1.0.0
`,
			},
			expectedErr: fmt.Sprintf("Invalid config: The following attributes are missing in %s:\n  * id\n", mockConfigFilePath),
			mockTarErr:  false,
		},
		{
			name: "Should complain that `author` is missing from pct-config.yml",
			args: args{
				templatePath: mockTemplateDir,
				targetDir:    mockTemplateDir,
			},
			mockDirs: []string{
				mockTemplateDir,
				filepath.Join(mockTemplateDir, "content"),
			},
			mockFiles: map[string]string{
				filepath.Clean(filepath.Join(mockTemplateDir, "pct-config.yml")): `---
template:
  id: builder
  version: 1.0.0
`,
			},
			expectedErr: fmt.Sprintf("Invalid config: The following attributes are missing in %s:\n  * author\n", mockConfigFilePath),
			mockTarErr:  false,
		},
		{
			name: "Should complain that `version` is missing from pct-config.yml",
			args: args{
				templatePath: mockTemplateDir,
				targetDir:    mockTemplateDir,
			},
			mockDirs: []string{
				mockTemplateDir,
				filepath.Join(mockTemplateDir, "content"),
			},
			mockFiles: map[string]string{
				filepath.Clean(filepath.Join(mockTemplateDir, "pct-config.yml")): `---
template:
  id: builder
  author: puppetlabs
`,
			},
			expectedErr: fmt.Sprintf("Invalid config: The following attributes are missing in %s:\n  * version\n", mockConfigFilePath),
			mockTarErr:  false,
		},
		{
			name: "Should complain all required are missing from pct-config.yml",
			args: args{
				templatePath: mockTemplateDir,
				targetDir:    mockTemplateDir,
			},
			mockDirs: []string{
				mockTemplateDir,
				filepath.Join(mockTemplateDir, "content"),
			},
			mockFiles: map[string]string{
				filepath.Clean(filepath.Join(mockTemplateDir, "pct-config.yml")): `---
template:
  foo: bar
`,
			},
			expectedErr: fmt.Sprintf("Invalid config: The following attributes are missing in %s:\n  * id\n  * author\n  * version\n", mockConfigFilePath),
			mockTarErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fs := afero.NewMemMapFs()
			afs := &afero.Afero{Fs: fs}

			for _, path := range tt.mockDirs {
				afs.Mkdir(path, 0750) //nolint:gosec,errcheck // this result is not used in a secure application
			}

			for file, content := range tt.mockFiles {
				config, _ := afs.Create(file) //nolint:gosec,errcheck // this result is not used in a secure application
				config.Write([]byte(content)) //nolint:errcheck
			}

			p := &pct.Builder{
				&mock.Tar{ReturnedPath: tt.tarFile, ErrResponse: tt.mockTarErr},
				&mock.Gzip{ReturnedPath: tt.gzipFile, ErrResponse: tt.mockGzipErr},
				afs,
				&pct_config_processor.PctConfigProcessor{AFS: afs},
				"pct-config.yml",
			}

			gotGzipArchiveFilePath, err := p.Build(tt.args.templatePath, tt.args.targetDir)
			if (err != nil) && tt.expectedErr != "" {
				assert.Equal(t, tt.expectedErr, err.Error())
				return
			}
			if gotGzipArchiveFilePath != tt.expectedFilePath {
				t.Errorf("Build() = %v, want %v", gotGzipArchiveFilePath, tt.expectedFilePath)
			}
		})
	}
}
