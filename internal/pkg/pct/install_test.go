package pct

import (
	"path/filepath"
	"testing"
)

type InstallTest struct {
	name                  string
	args                  args
	wantErr               bool
	expectedFilePath      string
	mockUntarErrResponse  error
	mockGunzipErrResponse error
	mockStatResponses     []mockStatResponse
}
type args struct {
	templatePath string
	targetDir    string
}

func TestInstall(t *testing.T) {

	untarUtil = untarHelpersImplMock{}
	gunzipUtil = gunzipHelpersImplMock{}
	osUtil = osUtilHelpersImplMock{}
	ioUtil = ioUtilHelpersImplMock{}

	tests := []InstallTest{
		{
			name: "should extract a tar.gx to a template folder",
			args: args{
				templatePath: filepath.Join(testDir, "good-project.tar.gz"),
				targetDir:    testDir,
			},
			mockStatResponses: []mockStatResponse{
				{
					// osUtil.Stat(filepath.Join(templatePath, "pct-config.yml"))
					expectedName: filepath.Clean(filepath.Join(testDir, "good-project.tar.gz")),
					mockError:    nil,
				},
			},
			expectedFilePath: filepath.Join(testDir, "good-project"),
			wantErr:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUntarErrResponse = tt.mockUntarErrResponse
			mockGunzipErrResponse = tt.mockGunzipErrResponse
			mockStatResponses = tt.mockStatResponses

			got, err := Install(tt.args.templatePath, tt.args.targetDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("Install() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expectedFilePath {
				t.Errorf("Install() = %v, want %v", got, tt.expectedFilePath)
			}
		})
	}
}
