package build

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func nullFunction(cmd *cobra.Command, args []string) error {
	return nil
}

func TestCreatebuildCommand(t *testing.T) {
	wd, _ := os.Getwd()
	defaultSourceDir := wd
	defaultTargetDir := filepath.Join(wd, "pkg")

	tests := []struct {
		name       string
		args       []string
		returnCode int
		out        string
		wantCmd    *cobra.Command
		wantErr    bool
		f          func(cmd *cobra.Command, args []string) error
		expSrcDir  string
		expTargDir string
	}{
		{
			name:       "executes without error for valid flag",
			args:       []string{"build"},
			f:          nullFunction,
			out:        "",
			wantErr:    false,
			expSrcDir:  defaultSourceDir,
			expTargDir: defaultTargetDir,
		},
		{
			name:       "executes with error for invalid flag",
			args:       []string{"--foo"},
			f:          nullFunction,
			out:        "unknown flag: --foo",
			wantErr:    true,
			expSrcDir:  "",
			expTargDir: "",
		},
		{
			name:       "uses sourcedir, targetdir when passed in",
			args:       []string{"build", "--sourcedir", "/path/to/template", "--targetdir", "/path/to/output"},
			f:          nullFunction,
			out:        "",
			wantErr:    false,
			expSrcDir:  "/path/to/template",
			expTargDir: "/path/to/output",
		},
		{
			name:       "Sets correct defaults if sourcedir and targetdir undefined",
			args:       []string{"build"},
			f:          nullFunction,
			out:        "",
			wantErr:    false,
			expSrcDir:  defaultSourceDir,
			expTargDir: defaultTargetDir,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := CreateCommand()
			b := bytes.NewBufferString("")
			cmd.SetOut(b)
			cmd.SetErr(b)
			cmd.SetArgs(tt.args)
			cmd.RunE = tt.f

			err := cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("executeTestUnit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			out, err := ioutil.ReadAll(b)
			if err != nil {
				t.Errorf("Failed to read stdout: %v", err)
				return
			}

			assert.Equal(t, tt.expSrcDir, sourceDir)
			assert.Equal(t, tt.expTargDir, targetDir)

			output := string(out)
			r := regexp.MustCompile(tt.out)
			if !r.MatchString(output) {
				t.Errorf("output did not match regexp /%s/\n> output\n%s\n", r, output)
				return
			}
		})
	}

}
