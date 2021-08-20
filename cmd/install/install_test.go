package install_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/puppetlabs/pdkgo/cmd/install"
	"github.com/puppetlabs/pdkgo/internal/pkg/mock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestCreateinstallCommand(t *testing.T) {
	tests := []struct {
		name                    string
		args                    []string
		expectError             bool
		expectedTemplatePkgPath string
		expectedTargetDir       string
		viperTemplatePath       string
		expectedOutput          string
	}{
		{
			name:           "Should error when no args provided",
			args:           []string{},
			expectError:    true,
			expectedOutput: "Path to template package (tar.gz) should be first argument",
		},
		{
			name:           "Should error when > 1 arg provided",
			args:           []string{"first/arg", "second/undeed/arg"},
			expectError:    true,
			expectedOutput: "Incorrect number of arguments; path to template package (tar.gz) should be first argument",
		},
		{
			name:                    "Sets TemplatePkgPath to passed arg and InstallPath to default template dir",
			args:                    []string{"/path/to/my-cool-template.tar.gz"},
			expectError:             false,
			expectedTemplatePkgPath: "/path/to/my-cool-template.tar.gz",
			expectedTargetDir:       "/the/default/location/for/templates",
			viperTemplatePath:       "/the/default/location/for/templates",
		},
		{
			name:                    "Sets TemplatePkgPath and InstallPath to passed args",
			args:                    []string{"/path/to/my-cool-template.tar.gz", "--templatepath", "/a/new/place/for/templates"},
			expectError:             false,
			expectedTemplatePkgPath: "/path/to/my-cool-template.tar.gz",
			expectedTargetDir:       "/a/new/place/for/templates",
			viperTemplatePath:       "/the/default/location/for/templates",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.SetDefault("templatepath", tt.viperTemplatePath)
			cmd := install.InstallCommand{
				PctInstaller: &mock.PctInstaller{
					ExpectedTemplatePkg: tt.expectedTemplatePkgPath,
					ExpectedTargetDir:   tt.expectedTargetDir,
				},
			}
			installCmd := cmd.CreateCommand()

			b := bytes.NewBufferString("")
			installCmd.SetOut(b)
			installCmd.SetErr(b)

			installCmd.SetArgs(tt.args)
			err := installCmd.Execute()

			if (err != nil) != tt.expectError {
				t.Errorf("executeTestUnit() error = %v, wantErr %v", err, tt.expectError)
				return
			}

			if (err != nil) && tt.expectError {
				out, _ := ioutil.ReadAll(b)
				assert.Contains(t, string(out), tt.expectedOutput)
			}

		})
	}
}
