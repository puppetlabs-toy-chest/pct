package pct_config_processor_test

import (
	"path/filepath"
	"testing"

	"github.com/puppetlabs/pdkgo/internal/pkg/pct_config_processor"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

type ProcessConfigTest struct {
	name     string
	args     args
	expected expected
	mocks    mocks
}

type args struct {
	targetDir string
	sourceDir string
	force     bool
}

type expected struct {
	errorMsg       string
	namespacedPath string
}

type mocks struct {
	dirs  []string
	files map[string]string
}

func TestProcessConfig(t *testing.T) {
	configDir := "path/to/config"

	tests := []ProcessConfigTest{
		{
			name:     "Config file is present and is correctly constructed",
			args:     args{targetDir: "templates", sourceDir: configDir, force: false},
			expected: expected{errorMsg: "", namespacedPath: filepath.Join("templates", "test-user/test-template/0.1.0")},
			mocks: mocks{
				dirs: []string{"templates"},
				files: map[string]string{
					filepath.Join(configDir, "pct-config.yml"): `---
template:
  id: test-template
  author: test-user
  version: 0.1.0
`,
				},
			},
		},
		{
			name:     "Config file does not exist",
			args:     args{targetDir: "templates", sourceDir: configDir, force: false},
			expected: expected{errorMsg: "Invalid config: "},
		},
		{
			name:     "Config files exists but has invalid yaml",
			args:     args{targetDir: "templates", sourceDir: configDir, force: false},
			expected: expected{errorMsg: "Invalid config: "},
			mocks: mocks{
				dirs: []string{"templates"},
				files: map[string]string{
					filepath.Join(configDir, "pct-config.yml"): `---
		template: id: test-template author: test-user version: 0.1.0
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

			configProcessor := pct_config_processor.PctConfigProcessor{AFS: afs}

			returnedPath, err := configProcessor.ProcessConfig(tt.args.sourceDir, tt.args.targetDir, tt.args.force)

			if tt.expected.errorMsg != "" && err != nil {
				assert.Contains(t, err.Error(), tt.expected.errorMsg)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected.namespacedPath, returnedPath)
		})
	}

}
