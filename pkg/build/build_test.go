package build_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/puppetlabs/pdkgo/internal/pkg/pct_config_processor"
	"github.com/puppetlabs/pdkgo/pkg/build"
	"github.com/puppetlabs/pdkgo/pkg/mock"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {

	type args struct {
		projectPath string
		targetDir   string
	}

	var mockSourceDir = "/path/to/my/cool-project"
	var mockConfigFilePath = filepath.Clean(filepath.Join(mockSourceDir, "my-config.yml"))

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
			name: "Should return err if project folder path does not exist",
			args: args{
				projectPath: mockSourceDir,
				targetDir:   mockSourceDir,
			},
			expectedFilePath: "",
			expectedErr:      "No project directory at /path/to/my/cool-project",
		},
		{
			name: "Should return err if project path does not contain my-config.yml",
			args: args{
				projectPath: mockSourceDir,
				targetDir:   mockSourceDir,
			},
			mockDirs: []string{
				mockSourceDir,
			},
			expectedFilePath: "",
			expectedErr:      "No 'my-config.yml' found in /path/to/my/cool-project",
		},
		{
			name: "Should return err if content dir does not exist",
			args: args{
				projectPath: mockSourceDir,
				targetDir:   mockSourceDir,
			},
			mockDirs: []string{
				mockSourceDir,
			},
			mockFiles: map[string]string{
				filepath.Clean(filepath.Join(mockSourceDir, "my-config.yml")): `---
template:
  id: builder
  author: puppetlabs
  version: 1.0.0
`,
			},
			expectedFilePath: "",
			expectedErr:      "No 'content' dir found in /path/to/my/cool-project",
		},
		{
			name: "Should not attempt to GZIP when TAR operation fails",
			args: args{
				projectPath: mockSourceDir,
				targetDir:   mockSourceDir,
			},
			mockDirs: []string{
				mockSourceDir,
				filepath.Join(mockSourceDir, "content"),
			},
			mockFiles: map[string]string{
				filepath.Clean(filepath.Join(mockSourceDir, "my-config.yml")): `---
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
				projectPath: mockSourceDir,
				targetDir:   mockSourceDir,
			},
			mockFiles: map[string]string{
				filepath.Clean(filepath.Join(mockSourceDir, "my-config.yml")): `---
template:
  id: builder
  author: puppetlabs
  version: 1.0.0
`,
			},
			mockDirs: []string{
				mockSourceDir,
				filepath.Join(mockSourceDir, "content"),
			},
			tarFile:          "/path/to/nowhere/pkg/nowhere.tar",
			expectedFilePath: "",
			expectedErr:      "gzip error",
			mockTarErr:       false,
			mockGzipErr:      true,
		},
		{
			name: "Should TAR.GZ valid project to $MODULE_ROOT/pkg and return path",
			args: args{
				projectPath: mockSourceDir,
				targetDir:   mockSourceDir,
			},
			mockDirs: []string{
				mockSourceDir,
				filepath.Join(mockSourceDir, "content"),
			},
			mockFiles: map[string]string{
				filepath.Clean(filepath.Join(mockSourceDir, "my-config.yml")): `---
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
			name: "Should complain that `id` is missing from my-config.yml",
			args: args{
				projectPath: mockSourceDir,
				targetDir:   mockSourceDir,
			},
			mockDirs: []string{
				mockSourceDir,
				filepath.Join(mockSourceDir, "content"),
			},
			mockFiles: map[string]string{
				filepath.Clean(filepath.Join(mockSourceDir, "my-config.yml")): `---
template:
  author: puppetlabs
  version: 1.0.0
`,
			},
			expectedErr: fmt.Sprintf("Invalid config: The following attributes are missing in %s:\n  * id\n", mockConfigFilePath),
			mockTarErr:  false,
		},
		{
			name: "Should complain that `author` is missing from my-config.yml",
			args: args{
				projectPath: mockSourceDir,
				targetDir:   mockSourceDir,
			},
			mockDirs: []string{
				mockSourceDir,
				filepath.Join(mockSourceDir, "content"),
			},
			mockFiles: map[string]string{
				filepath.Clean(filepath.Join(mockSourceDir, "my-config.yml")): `---
template:
  id: builder
  version: 1.0.0
`,
			},
			expectedErr: fmt.Sprintf("Invalid config: The following attributes are missing in %s:\n  * author\n", mockConfigFilePath),
			mockTarErr:  false,
		},
		{
			name: "Should complain that `version` is missing from my-config.yml",
			args: args{
				projectPath: mockSourceDir,
				targetDir:   mockSourceDir,
			},
			mockDirs: []string{
				mockSourceDir,
				filepath.Join(mockSourceDir, "content"),
			},
			mockFiles: map[string]string{
				filepath.Clean(filepath.Join(mockSourceDir, "my-config.yml")): `---
template:
  id: builder
  author: puppetlabs
`,
			},
			expectedErr: fmt.Sprintf("Invalid config: The following attributes are missing in %s:\n  * version\n", mockConfigFilePath),
			mockTarErr:  false,
		},
		{
			name: "Should complain all required are missing from my-config.yml",
			args: args{
				projectPath: mockSourceDir,
				targetDir:   mockSourceDir,
			},
			mockDirs: []string{
				mockSourceDir,
				filepath.Join(mockSourceDir, "content"),
			},
			mockFiles: map[string]string{
				filepath.Clean(filepath.Join(mockSourceDir, "my-config.yml")): `---
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

			p := &build.Builder{
				&mock.Tar{ReturnedPath: tt.tarFile, ErrResponse: tt.mockTarErr},
				&mock.Gzip{ReturnedPath: tt.gzipFile, ErrResponse: tt.mockGzipErr},
				afs,
				&pct_config_processor.PctConfigProcessor{AFS: afs},
				"my-config.yml",
			}

			gotGzipArchiveFilePath, err := p.Build(tt.args.projectPath, tt.args.targetDir)
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
