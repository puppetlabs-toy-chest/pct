package pct_test

import (
	"path/filepath"
	"testing"

	"github.com/puppetlabs/pdkgo/internal/pkg/mock"
	"github.com/puppetlabs/pdkgo/internal/pkg/pct"
	"github.com/spf13/afero"
)

func TestBuild(t *testing.T) {

	type args struct {
		templatePath string
		targetDir    string
	}

	var mockTemplateDir = "/path/to/my/cool-template"

	tests := []struct {
		name                    string
		args                    args
		mockIsModuleRootErrResp error
		mockDirs                []string
		mockFiles               []string
		expectedFilePath        string
		tarFile                 string
		gzipFile                string
		wantErr                 bool
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
			wantErr:          true,
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
			wantErr:          true,
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
			mockFiles: []string{
				filepath.Clean(filepath.Join(mockTemplateDir, "pct-config.yml")),
			},
			expectedFilePath: "",
			wantErr:          true,
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
			mockFiles: []string{
				filepath.Join(mockTemplateDir, "pct-config.yml"),
			},
			expectedFilePath: "",
			wantErr:          true,
			mockTarErr:       true,
		},
		{
			name: "Should return error and empty path if GZIP operation fails",
			args: args{
				templatePath: mockTemplateDir,
				targetDir:    mockTemplateDir,
			},
			mockFiles: []string{
				filepath.Join(mockTemplateDir, "pct-config.yml"),
			},
			mockDirs: []string{
				mockTemplateDir,
				filepath.Join(mockTemplateDir, "content"),
			},
			tarFile:          "/path/to/nowhere/pkg/nowhere.tar",
			expectedFilePath: "",
			wantErr:          true,
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
			mockFiles: []string{
				filepath.Join(mockTemplateDir, "pct-config.yml"),
			},
			tarFile:          "/path/to/nowhere/pkg/nowhere.tar",
			gzipFile:         "/path/to/nowhere/pkg/nowhere.tar.gz",
			expectedFilePath: "/path/to/nowhere/pkg/nowhere.tar.gz",
			wantErr:          false,
			mockTarErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fs := afero.NewMemMapFs()
			afs := &afero.Afero{Fs: fs}

			for _, path := range tt.mockDirs {
				afs.Mkdir(path, 0750) //nolint:errcheck
			}

			for _, path := range tt.mockFiles {
				afs.Create(path) //nolint:gosec,errcheck // this result is not used in a secure application
			}

			p := &pct.Builder{
				&mock.Tar{ReturnedPath: tt.tarFile, ErrResponse: tt.mockTarErr},
				&mock.Gzip{ReturnedPath: tt.gzipFile, ErrResponse: tt.mockGzipErr},
				afs,
			}

			gotGzipArchiveFilePath, err := p.Build(tt.args.templatePath, tt.args.targetDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotGzipArchiveFilePath != tt.expectedFilePath {
				t.Errorf("Build() = %v, want %v", gotGzipArchiveFilePath, tt.expectedFilePath)
			}
		})
	}
}
