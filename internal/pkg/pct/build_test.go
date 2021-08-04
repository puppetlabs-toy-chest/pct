package pct

import (
	"errors"
	"path/filepath"
	"testing"
)

func TestBuild(t *testing.T) {
	osUtil = osUtilHelpersImplMock{}
	tarUtil = tarHelpersImplMock{}
	ioUtil = ioUtilHelpersImplMock{}
	gzipUtil = gzipHelpersImplMock{}

	type args struct {
		templatePath string
		targetDir    string
	}

	tests := []struct {
		name                    string
		args                    args
		mockStatResponses       []mockStatResponse
		expectedFilePath        string
		wantErr                 bool
		mockTarErrResponse      error
		mockGzipErrResponse     error
		testTempDir             string
	}{
		{
			name: "Should return err if template path does not exist",
			args: args{
				templatePath: testDir,
				targetDir:    testDir,
			},
			mockStatResponses: []mockStatResponse{
				{
					// osUtil.Stat(templatePath)
					expectedName: testDir,
					mockError:    errors.New("Template path does not exist"),
				},
			},
			expectedFilePath: "",
			wantErr:          true,
		},
		{
			name: "Should return err if template path does not contain pct-config.yml",
			args: args{
				templatePath: testDir,
				targetDir:    testDir,
			},
			mockStatResponses: []mockStatResponse{
				{
					// osUtil.Stat(templatePath)
					expectedName: testDir,
					mockError:    nil,
				},
				{
					// osUtil.Stat(filepath.Join(templatePath, "pct-config.yml"))
					expectedName: filepath.Clean(filepath.Join(testDir, "pct-config.yml")),
					mockError:    errors.New("No pct-config.yml found"),
				},
			},
			expectedFilePath: "",
			wantErr:          true,
		},
		{
			name: "Should return err if content dir does not exist",
			args: args{
				templatePath: testDir,
				targetDir:    testDir,
			},
			mockStatResponses: []mockStatResponse{
				{
					// osUtil.Stat(templatePath)
					expectedName: testDir,
					mockError:    nil,
				},
				{
					// osUtil.Stat(filepath.Join(templatePath, "pct-config.yml"))
					expectedName: filepath.Clean(filepath.Join(testDir, "pct-config.yml")),
					mockError:    nil,
				},
				{
					// osUtil.Stat(filepath.Join(templatePath, "content"))
					expectedName: filepath.Clean(filepath.Join(testDir, "content")),
					mockError:    errors.New("No content dir found"),
				},
			},
			expectedFilePath: "",
			wantErr:          true,
		},
		{
			name: "Should not attempt to GZIP when TAR operation fails",
			args: args{
				templatePath: testDir,
				targetDir:    testDir,
			},
			mockStatResponses: []mockStatResponse{
				{
					// osUtil.Stat(templatePath)
					expectedName: testDir,
					mockError:    nil,
				},
				{
					// osUtil.Stat(filepath.Join(templatePath, "pct-config.yml"))
					expectedName: filepath.Clean(filepath.Join(testDir, "pct-config.yml")),
					mockError:    nil,
				},
				{
					// osUtil.Stat(filepath.Join(templatePath, "content"))
					expectedName: filepath.Clean(filepath.Join(testDir, "content")),
					mockError:    nil,
				},
			},
			expectedFilePath:   "",
			wantErr:            true,
			mockTarErrResponse: errors.New("Could not TAR the directory"),
		},
		{
			name: "Should TAR.GZ valid template to $MODULE_ROOT/pkg and return path",
			args: args{
				templatePath: testDir,
				targetDir:    testDir,
			},
			mockStatResponses: []mockStatResponse{
				{
					// osUtil.Stat(templatePath)
					expectedName: testDir,
					mockError:    nil,
				},
				{
					// osUtil.Stat(filepath.Join(templatePath, "pct-config.yml"))
					expectedName: filepath.Clean(filepath.Join(testDir, "pct-config.yml")),
					mockError:    nil,
				},
				{
					// osUtil.Stat(filepath.Join(templatePath, "content"))
					expectedName: filepath.Clean(filepath.Join(testDir, "content")),
					mockError:    nil,
				},
			},
			expectedFilePath:   filepath.Clean("/path/to/nowhere/pkg/template.tar.gz"),
			wantErr:            false,
			mockTarErrResponse: nil,
		},
		{
			name: "Should return error and empty path is GZIP operation fails",
			args: args{
				templatePath: testDir,
				targetDir:    testDir,
			},
			mockStatResponses: []mockStatResponse{
				{
					// osUtil.Stat(templatePath)
					expectedName: testDir,
					mockError:    nil,
				},
				{
					// osUtil.Stat(filepath.Join(templatePath, "pct-config.yml"))
					expectedName: filepath.Clean(filepath.Join(testDir, "pct-config.yml")),
					mockError:    nil,
				},
				{
					// osUtil.Stat(filepath.Join(templatePath, "content"))
					expectedName: filepath.Clean(filepath.Join(testDir, "content")),
					mockError:    nil,
				},
			},
			expectedFilePath:    "",
			wantErr:             true,
			mockTarErrResponse:  nil,
			mockGzipErrResponse: errors.New("Could not GZIP the TAR"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStatResponses = tt.mockStatResponses
			mockTarErrResponse = tt.mockTarErrResponse
			mockGzipErrResponse = tt.mockGzipErrResponse
			gotGzipArchiveFilePath, err := Build(tt.args.templatePath, tt.args.targetDir)
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
