package install

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func nullFunction(cmd *cobra.Command, args []string) error {
	return nil
}

func TestCreateinstallCommand(t *testing.T) {
	// Viper test params
	viperTemplatePathValue := "/unit/test/template/path"
	viper.Set("templatepath", viperTemplatePathValue)

	// Template Path
	testTemplateFilePath := "/path/to/template.tar.gz"

	// Install Path
	testInstallPath := "/install/templates/here"

	tests := []struct {
		name               string
		args               []string
		returnCode         int
		out                string
		wantErr            bool
		expErrMsg          string
		expTemplatePkgPath string
		expInstallPath     string
	}{
		{
			name:      "errors when no template package provided",
			args:      []string{},
			wantErr:   true,
			expErrMsg: "Must specify a path to a template package",
		},
		{
			name:               "sets templatePkgPath and default installPath correctly",
			args:               []string{testTemplateFilePath},
			wantErr:            false,
			expTemplatePkgPath: testTemplateFilePath,
			expInstallPath:     viperTemplatePathValue,
		},
		{
			name:               "sets templatePkgPath and defined installPath correctly",
			args:               []string{testTemplateFilePath, "--templatepath", testInstallPath},
			wantErr:            false,
			expTemplatePkgPath: testTemplateFilePath,
			expInstallPath:     testInstallPath,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := CreateCommand()
			b := bytes.NewBufferString("")
			cmd.SetOut(b)
			cmd.SetErr(b)
			cmd.SetArgs(tt.args)
			cmd.RunE = nullFunction

			err := cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("executeTestUnit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.expTemplatePkgPath, templatePkgPath)
			assert.Equal(t, tt.expInstallPath, installPath)
		})

	}
}
