package build_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/puppetlabs/pdkgo/cmd/build"
	"github.com/puppetlabs/pdkgo/pkg/mock"
)

func TestCreateBuildCommand(t *testing.T) {
	wd, _ := os.Getwd()
	defaultSourceDir := wd
	defaultTargetDir := filepath.Join(wd, "pkg")

	tests := []struct {
		name               string
		args               []string
		expectedSourceDir  string
		expectedTargetDir  string
		expectedErrorMatch string
	}{
		{
			name:              "executes without error when no flags passed",
			args:              []string{},
			expectedSourceDir: defaultSourceDir,
			expectedTargetDir: defaultTargetDir,
		},
		{
			name:               "executes with error for invalid flag",
			args:               []string{"--foo"},
			expectedErrorMatch: "unknown flag: --foo",
		},
		{
			name:              "uses sourcedir, targetdir when passed in",
			args:              []string{"--sourcedir", "/path/to/template", "--targetdir", "/path/to/output"},
			expectedSourceDir: "/path/to/template",
			expectedTargetDir: "/path/to/output",
		},
		{
			name:              "Sets correct defaults if sourcedir and targetdir undefined",
			args:              []string{},
			expectedSourceDir: defaultSourceDir,
			expectedTargetDir: defaultTargetDir,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := build.BuildCommand{
				ProjectType: "template",
				Builder: &mock.Builder{
					ProjectName:       "my-project",
					ExpectedSourceDir: tt.expectedSourceDir,
					ExpectedTargetDir: tt.expectedTargetDir,
				},
			}
			buildCmd := cmd.CreateCommand()

			b := bytes.NewBufferString("")
			buildCmd.SetOut(b)
			buildCmd.SetErr(b)

			buildCmd.SetArgs(tt.args)
			err := buildCmd.Execute()

			if err != nil {
				if tt.expectedErrorMatch == "" {
					t.Errorf("Unexpected error when none wanted: %v", err)
					return
				} else {
					out, _ := ioutil.ReadAll(b)
					assert.Regexp(t, tt.expectedErrorMatch, string(out))
				}
			} else if tt.expectedErrorMatch != "" {
				t.Errorf("Expected error '%s'but none raised", err)
			}
		})
	}
}
