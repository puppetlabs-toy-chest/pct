package pct_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/puppetlabs/pdkgo/internal/pkg/mock"
	"github.com/puppetlabs/pdkgo/internal/pkg/pct"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestMain(m *testing.M) {
	// hide logging output
	log.Logger = zerolog.New(ioutil.Discard).With().Timestamp().Logger()
	os.Exit(m.Run())
}

func TestDeploy(t *testing.T) {
	type args struct {
		info            pct.DeployInfo
		templateConfig  string
		templateContent map[string]string
	}

	tmp := t.TempDir()

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "deploy a project and return the correct new files",
			args: args{
				info: pct.DeployInfo{
					SelectedTemplate: "full-project",
					TemplateCache:    "testdata/examples",
					TargetOutputDir:  filepath.Join(tmp, "foobar"),
					TargetName:       "woo",
					PdkInfo: pct.PDKInfo{
						Version:   "0.1.0",
						Commit:    "abc12345",
						BuildDate: "2021/06/27",
					},
				},
				templateConfig: `---
template:
  id: full-project
  type: project

`,
				templateContent: map[string]string{
					"metadata.json": "fixed string content",
				},
			},
			want: []string{
				filepath.Join(tmp, "foobar", "woo"),
				filepath.Join(tmp, "foobar", "woo", "metadata.json"),
			},
		},
		{
			name: "deploy a project without a name or outputDir",
			args: args{
				info: pct.DeployInfo{
					SelectedTemplate: "full-project",
					TemplateCache:    "testdata/examples",
					TargetOutputDir:  "",
					TargetName:       "",
					PdkInfo: pct.PDKInfo{
						Version:   "0.1.0",
						Commit:    "abc12345",
						BuildDate: "2021/06/27",
					},
				},
				templateConfig: `---
template:
  id: full-project
  type: project

`,
				templateContent: map[string]string{
					"metadata.json": "fixed string content",
				},
			},
			want: []string{
				tmp,
				filepath.Join(tmp, "metadata.json"),
			},
		},
		{
			name: "deploy a project without an outputDir",
			args: args{
				info: pct.DeployInfo{
					SelectedTemplate: "full-project",
					TemplateCache:    "testdata/examples",
					TargetOutputDir:  "",
					TargetName:       "wibble",
					PdkInfo: pct.PDKInfo{
						Version:   "0.1.0",
						Commit:    "abc12345",
						BuildDate: "2021/06/27",
					},
				},
				templateConfig: `---
template:
  id: full-project
  type: project

`,
				templateContent: map[string]string{
					"metadata.json": "fixed string content",
				},
			},
			want: []string{
				filepath.Join(tmp, "wibble"),
				filepath.Join(tmp, "wibble", "metadata.json"),
			},
		},
		{
			name: "deploy a item and return the correctly named new files",
			args: args{
				info: pct.DeployInfo{
					SelectedTemplate: "replace-thing",
					TemplateCache:    "testdata/examples",
					TargetOutputDir:  filepath.Join(tmp, "thing"),
					TargetName:       "woo",
					PdkInfo: pct.PDKInfo{
						Version:   "0.1.0",
						Commit:    "abc12345",
						BuildDate: "2021/06/27",
					},
				},
				templateConfig: `---
template:
  id: replace-thing
  type: item

`,
				templateContent: map[string]string{
					"{{pct_name}}.txt.tmpl": `This is example text

Summary: {{.example_replace.summary}}`,
				},
			},
			want: []string{
				filepath.Join(tmp, "thing"),
				filepath.Join(tmp, "thing", "woo.txt"),
			},
		},
		{
			name: "deploy a item without a name or outputDir",
			args: args{
				info: pct.DeployInfo{
					SelectedTemplate: "replace-thing",
					TemplateCache:    "testdata/examples",
					TargetOutputDir:  "",
					TargetName:       "",
					PdkInfo: pct.PDKInfo{
						Version:   "0.1.0",
						Commit:    "abc12345",
						BuildDate: "2021/06/27",
					},
				},
				templateConfig: `---
template:
  id: replace-thing
  type: item

`,
				templateContent: map[string]string{
					"{{pct_name}}.txt.tmpl": `This is example text

Summary: {{.example_replace.summary}}`,
				},
			},
			want: []string{
				tmp,
				filepath.Join(tmp, filepath.Base(tmp)+".txt"),
			},
		},
		{
			name: "deploy a item without an outputDir",
			args: args{
				info: pct.DeployInfo{
					SelectedTemplate: "replace-thing",
					TemplateCache:    "testdata/examples",
					TargetOutputDir:  "",
					TargetName:       "wibble",
					PdkInfo: pct.PDKInfo{
						Version:   "0.1.0",
						Commit:    "abc12345",
						BuildDate: "2021/06/27",
					},
				},
				templateConfig: `---
template:
  id: replace-thing
  type: item

`,
				templateContent: map[string]string{
					"{{pct_name}}.txt.tmpl": `This is example text

Summary: {{.example_replace.summary}}`,
				},
			},
			want: []string{
				tmp,
				filepath.Join(tmp, "wibble.txt"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			afs := &afero.Afero{Fs: fs}
			iofs := &afero.IOFS{Fs: fs}

			// Create the template
			templateDir := filepath.Join(tt.args.info.TemplateCache, tt.args.info.SelectedTemplate)
			contentDir := filepath.Join(templateDir, "content")
			afs.MkdirAll(contentDir, 0750) //nolint:errcheck
			// Create template config
			config, _ := afs.Create(filepath.Join(templateDir, "pct-config.yml"))
			config.Write([]byte(tt.args.templateConfig)) //nolint:errcheck
			// Create the contents
			for file, content := range tt.args.templateContent {
				nf, _ := afs.Create(filepath.Join(contentDir, file))
				nf.Write([]byte(content)) //nolint:errcheck
			}

			p := &pct.Pct{
				&mock.OsUtil{WD: tmp},
				&mock.UtilsHelper{TestDir: tmp},
				afs,
				iofs,
			}

			if got := p.Deploy(tt.args.info); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Deploy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGet(t *testing.T) {
	type args struct {
		templateCache    string
		selectedTemplate string
		setup            bool
		templateConfig   string
		templateContent  map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    pct.PuppetContentTemplate
		wantErr bool
	}{
		{
			name: "returns error for non-existent template",
			args: args{
				templateCache:    "testdata/examples",
				selectedTemplate: "i-dont-exist",
				setup:            false,
			},
			wantErr: true,
		},
		{
			name: "returns tmpl for existent template",
			args: args{
				templateCache:    "testdata/examples",
				selectedTemplate: "full-project",
				setup:            true,
				templateConfig: `---
template:
  id: full-project
  type: project
  display: Full Project
  version: 0.1.0
  url: https://github.com/puppetlabs/pct-full-project
`,
				templateContent: map[string]string{
					"metadata.json": "fixed string content",
				},
			},
			want: pct.PuppetContentTemplate{
				Id:      "full-project",
				Type:    "project",
				Display: "Full Project",
				Version: "0.1.0",
				URL:     "https://github.com/puppetlabs/pct-full-project",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			afs := &afero.Afero{Fs: fs}
			iofs := &afero.IOFS{Fs: fs}

			if tt.args.setup {
				// Create the template
				templateDir := filepath.Join(tt.args.templateCache, tt.args.selectedTemplate)
				contentDir := filepath.Join(templateDir, "content")
				afs.MkdirAll(contentDir, 0750) //nolint:errcheck
				// Create template config
				config, _ := afs.Create(filepath.Join(templateDir, "pct-config.yml"))
				config.Write([]byte(tt.args.templateConfig)) //nolint:errcheck
				// Create the contents
				for file, content := range tt.args.templateContent {
					nf, _ := afs.Create(filepath.Join(contentDir, file))
					nf.Write([]byte(content)) //nolint:errcheck
				}
			}

			p := &pct.Pct{
				&mock.OsUtil{},
				&mock.UtilsHelper{},
				afs,
				iofs,
			}

			got, err := p.Get(tt.args.templateCache, tt.args.selectedTemplate)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDisplayDefaults(t *testing.T) {
	type args struct {
		defaults map[string]interface{}
	}
	tests := []struct {
		name   string
		args   args
		format string
		want   string
	}{
		{
			name: "table example",
			args: args{
				defaults: map[string]interface{}{
					"foo": map[string]interface{}{
						"bar": "wibble",
						"wobble": []string{
							"one", "two", "three",
						},
					},
				},
			},
			format: "table",
			want: `foo:
  bar: wibble
  wobble:
  - one
  - two
  - three
`,
		},
		{
			name: "json example",
			args: args{
				defaults: map[string]interface{}{
					"foo": map[string]interface{}{
						"bar": "wibble",
						"wobble": []string{
							"one", "two", "three",
						},
					},
				},
			},
			format: "json",
			want: `{
  "foo": {
    "bar": "wibble",
    "wobble": [
      "one",
      "two",
      "three"
    ]
  }
}`,
		},
		{
			name: "empty table example",
			args: args{
				defaults: map[string]interface{}{},
			},
			format: "table",
			want:   "This template has no configuration options.",
		},
		{
			name: "empty json example",
			args: args{
				defaults: map[string]interface{}{},
			},
			format: "json",
			want:   "{\n  \n}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			p := &pct.Pct{
				&mock.OsUtil{},
				&mock.UtilsHelper{},
				&afero.Afero{Fs: fs},
				&afero.IOFS{Fs: fs},
			}

			returnString := p.DisplayDefaults(tt.args.defaults, tt.format)

			if len(tt.args.defaults) == 0 {
				assert.Equal(t, returnString, tt.want)

				return
			}

			var output map[string]interface{}
			var expected map[string]interface{}

			if tt.format == "table" {
				err := yaml.Unmarshal([]byte(returnString), &output)
				if err != nil {
					assert.Fail(t, "returned data is not YAML")
				}

				err = yaml.Unmarshal([]byte(tt.want), &expected)
				if err != nil {
					assert.Fail(t, "expected data is not YAML")
				}
			} else if tt.format == "json" {
				err := json.Unmarshal([]byte(returnString), &output)
				if err != nil {
					assert.Fail(t, "returned data is not JSON")
				}

				err = json.Unmarshal([]byte(tt.want), &expected)
				if err != nil {
					assert.Fail(t, "expected data is not JSON")
				}
			}

			assert.Equal(t, expected, output)
		})
	}
}

func TestFormatTemplates(t *testing.T) {
	type args struct {
		tmpls      []pct.PuppetContentTemplate
		jsonOutput string
	}
	tests := []struct {
		name    string
		p       *pct.Pct
		args    args
		matches []string
		wantErr bool
	}{
		{
			name: "When no templates are passed",
			args: args{
				tmpls:      []pct.PuppetContentTemplate{},
				jsonOutput: "table",
			},
			matches: []string{},
		},
		{
			name: "When only one template is passed",
			args: args{
				tmpls: []pct.PuppetContentTemplate{
					{
						Id:      "foo",
						Author:  "bar",
						Type:    "Item",
						Display: "Foo Item",
						Version: "0.1.0",
						URL:     "https://github.com/puppetlabs/pct-good-project",
					},
				},
				jsonOutput: "table",
			},
			matches: []string{
				`DisplayName:\s+Foo Item`,
				`Author:\s+bar`,
				`Name:\s+foo`,
				`TemplateType:\s+Item`,
				`TemplateURL:\s+https://github.com/puppetlabs/pct-good-project`,
				`TemplateVersion:\s+0\.1\.0`,
			},
		},
		{
			name: "When more than one template is passed",
			args: args{
				tmpls: []pct.PuppetContentTemplate{
					{
						Id:      "foo",
						Author:  "baz",
						Type:    "Item",
						Display: "Foo Item",
						Version: "0.1.0",
						URL:     "https://github.com/puppetlabs/pct-good-project",
					},
					{
						Id:      "bar",
						Author:  "baz",
						Type:    "Item",
						Display: "Bar Item",
						Version: "0.1.0",
						URL:     "https://github.com/puppetlabs/pct-good-project",
					},
				},
				jsonOutput: "table",
			},
			matches: []string{
				`DISPLAYNAME \| AUTHOR \| NAME \| TYPE`,
				`Foo Item\s+\|\sbaz\s+\|\sfoo\s+\|\sItem`,
				`Bar Item\s+\|\sbaz\s+\|\sbar\s+\|\sItem`,
			},
		},
		{
			name: "When format is specified as json",
			args: args{
				tmpls: []pct.PuppetContentTemplate{
					{
						Id:      "foo",
						Author:  "baz",
						Type:    "Item",
						Display: "Foo Item",
						Version: "0.1.0",
						URL:     "https://github.com/puppetlabs/pct-good-project",
					},
					{
						Id:      "bar",
						Author:  "baz",
						Type:    "Item",
						Display: "Bar Item",
						Version: "0.1.0",
						URL:     "https://github.com/puppetlabs/pct-good-project",
					},
				},
				jsonOutput: "json",
			},
			matches: []string{
				`\"Id\": \"foo\"`,
				`\"Id\": \"bar\"`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := tt.p.FormatTemplates(tt.args.tmpls, tt.args.jsonOutput)
			if err != nil {
				t.Errorf("Pct.FormatTemplates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, m := range tt.matches {
				assert.Regexp(t, m, output)
			}
		})
	}
}

func TestList(t *testing.T) {
	type stubbedConfig struct {
		relativeConfigPath string
		configContent      string
	}
	type args struct {
		templatePath   string
		templateName   string
		stubbedConfigs []stubbedConfig
	}
	tests := []struct {
		name string
		args args
		want []pct.PuppetContentTemplate
	}{
		{
			name: "when no templates are found",
			args: args{
				templatePath: "stubbed/templates/none",
			},
		},
		{
			name: "when an invalid template is found",
			args: args{
				templatePath: "stubbed/templates/invalid",
				stubbedConfigs: []stubbedConfig{
					{
						relativeConfigPath: "some_author/bad-template/0.1.0",
						configContent:      "I am WILDLY INVALID",
					},
				},
			},
		},
		{
			name: "when valid templates are found",
			args: args{
				templatePath: "stubbed/templates/valid",
				stubbedConfigs: []stubbedConfig{
					{
						relativeConfigPath: "some_author/first/0.1.0",
						configContent: `---
template:
  author: some_author
  id: first
  type: project
  display: First Template
  version: 0.1.0
  url: https://github.com/some_author/pct-first-template
`,
					},
					{
						relativeConfigPath: "some_author/second/0.1.0",
						configContent: `---
template:
  author: some_author
  id: second
  type: project
  display: Second Template
  version: 0.1.0
  url: https://github.com/some_author/pct-second-template
`,
					},
				},
			},
			want: []pct.PuppetContentTemplate{
				{
					Author:  "some_author",
					Id:      "first",
					Type:    "project",
					Display: "First Template",
					Version: "0.1.0",
					URL:     "https://github.com/some_author/pct-first-template",
				},
				{
					Author:  "some_author",
					Id:      "second",
					Type:    "project",
					Display: "Second Template",
					Version: "0.1.0",
					URL:     "https://github.com/some_author/pct-second-template",
				},
			},
		},
		{
			name: "when templates are found with the same author/id and different versions",
			args: args{
				templatePath: "stubbed/templates/multiversion",
				stubbedConfigs: []stubbedConfig{
					{
						relativeConfigPath: "some_author/first/0.1.0",
						configContent: `---
template:
  author: some_author
  id: first
  type: project
  display: First Template
  version: 0.1.0
  url: https://github.com/some_author/pct-first-template
`,
					},
					{
						relativeConfigPath: "some_author/first/0.2.0",
						configContent: `---
template:
  author: some_author
  id: first
  type: project
  display: First Template
  version: 0.2.0
  url: https://github.com/some_author/pct-first-template
`,
					},
				},
			},
			want: []pct.PuppetContentTemplate{
				{
					Author:  "some_author",
					Id:      "first",
					Type:    "project",
					Display: "First Template",
					Version: "0.2.0",
					URL:     "https://github.com/some_author/pct-first-template",
				},
			},
		},
		{
			name: "when templateName is specified",
			args: args{
				templatePath: "stubbed/templates/named",
				templateName: "first",
				stubbedConfigs: []stubbedConfig{
					{
						relativeConfigPath: "some_author/first/0.1.0",
						configContent: `---
template:
  author: some_author
  id: first
  type: project
  display: First Template
  version: 0.1.0
  url: https://github.com/some_author/pct-first-template
`,
					},
					{
						relativeConfigPath: "some_author/second/0.1.0",
						configContent: `---
template:
  author: some_author
  id: second
  type: project
  display: Second Template
  version: 0.1.0
  url: https://github.com/some_author/pct-second-template
`,
					},
				},
			},
			want: []pct.PuppetContentTemplate{
				{
					Author:  "some_author",
					Id:      "first",
					Type:    "project",
					Display: "First Template",
					Version: "0.1.0",
					URL:     "https://github.com/some_author/pct-first-template",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			afs := &afero.Afero{Fs: fs}
			iofs := &afero.IOFS{Fs: fs}

			// Create the template
			for _, st := range tt.args.stubbedConfigs {
				templateDir := filepath.Join(tt.args.templatePath, st.relativeConfigPath)
				afs.MkdirAll(templateDir, 0750) //nolint:errcheck
				// Create template config
				config, _ := afs.Create(filepath.Join(templateDir, "pct-config.yml"))
				config.Write([]byte(st.configContent)) //nolint:errcheck
			}

			p := &pct.Pct{
				&mock.OsUtil{},
				&mock.UtilsHelper{},
				afs,
				iofs,
			}

			got := p.List(tt.args.templatePath, tt.args.templateName)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pct.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func Test_createTemplateFile(t *testing.T) {
// 	type args struct {
// 		info         pct.DeployInfo
// 		configFile   string
// 		templateFile pct.PuppetContentTemplateFileInfo
// 		tmpl         pct.PuppetContentTemplate
// 	}

// 	tmp := t.TempDir()

// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "",
// 			args: args{
// 				info: pct.DeployInfo{
// 					TargetName: "foobar",
// 					PdkInfo: pct.PDKInfo{
// 						Version:   "0.1.0",
// 						Commit:    "abc12345",
// 						BuildDate: "2021/06/27",
// 					},
// 				},
// 				configFile: "testdata/examples/good-project/pct.yml",
// 				tmpl: pct.PuppetContentTemplate{
// 					Type:    "project",
// 					Display: "Good Project",
// 					URL:     "https://github.com/puppetlabs/pct-good-project",
// 					Version: "0.1.0",
// 					Id:      "good-project",
// 				},
// 				templateFile: pct.PuppetContentTemplateFileInfo{
// 					TemplatePath:   "testdata/examples/good-project/content/goodfile.txt.tmpl",
// 					TargetFilePath: filepath.Join(tmp, "foo.txt"),
// 					TargetDir:      tmp,
// 					TargetFile:     "",
// 					IsDirectory:    false,
// 				},
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if err := pct.createTemplateFile(tt.args.info, tt.args.configFile, tt.args.templateFile, tt.args.tmpl); (err != nil) != tt.wantErr {
// 				t.Errorf("createTemplateFile() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 			if _, err := os.Stat(tt.args.templateFile.TargetFilePath); err != nil {
// 				t.Errorf("createTemplateFile() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func Test_processConfiguration(t *testing.T) {
// 	type args struct {
// 		info            pct.DeployInfo
// 		configFile      string
// 		projectTemplate string
// 		tmpl            pct.PuppetContentTemplate
// 	}
// 	cwd, _ := os.Getwd()
// 	hostName, _ := os.Hostname()
// 	u := pct.getCurrentUser()
// 	tests := []struct {
// 		name string
// 		args args
// 		want map[string]interface{}
// 	}{
// 		{
// 			name: "with a valid config, returns a correct map interface",
// 			args: args{
// 				info: pct.DeployInfo{
// 					TargetName: "good-project",
// 					PdkInfo: pct.PDKInfo{
// 						Version:   "0.1.0",
// 						Commit:    "abc12345",
// 						BuildDate: "2021/06/27",
// 					},
// 				},
// 				configFile:      "testdata/examples/good-project/pct.yml",
// 				projectTemplate: "",
// 				tmpl:            pct.PuppetContentTemplate{},
// 			},
// 			want: map[string]interface{}{
// 				"user":     u,
// 				"cwd":      cwd,
// 				"hostname": hostName,
// 				"pct_name": "good-project",
// 				"pdk": map[string]interface{}{
// 					"build_date":  "2021/06/27",
// 					"commit_hash": "abc12345",
// 					"version":     "0.1.0",
// 				},
// 				"template": map[string]interface{}{
// 					"type":    "project",
// 					"display": "Good Project",
// 					"url":     "https://github.com/puppetlabs/pct-good-project",
// 					"version": "0.1.0",
// 					"id":      "good-project",
// 				},
// 				"puppet_module": map[string]interface{}{
// 					"author":  u,
// 					"license": "Apache-2.0",
// 					"version": "0.1.0",
// 					"summary": "A New Puppet Module",
// 				},
// 			},
// 		},
// 		{
// 			name: "with a valid config, and a workspace overide, returns a correct map interface",
// 			args: args{
// 				info: pct.DeployInfo{
// 					TargetName:      "good-project",
// 					TargetOutputDir: "testdata/examples/good-project-override",
// 					PdkInfo: pct.PDKInfo{
// 						Version:   "0.1.0",
// 						Commit:    "abc12345",
// 						BuildDate: "2021/06/27",
// 					},
// 				},
// 				configFile:      "testdata/examples/good-project/pct.yml",
// 				projectTemplate: "",
// 				tmpl:            pct.PuppetContentTemplate{},
// 			},
// 			want: map[string]interface{}{
// 				"user":     u,
// 				"cwd":      cwd,
// 				"hostname": hostName,
// 				"pct_name": "good-project",
// 				"pdk": map[string]interface{}{
// 					"build_date":  "2021/06/27",
// 					"commit_hash": "abc12345",
// 					"version":     "0.1.0",
// 				},
// 				"template": map[string]interface{}{
// 					"type":    "project",
// 					"display": "Good Project",
// 					"url":     "https://github.com/puppetlabs/pct-good-project",
// 					"version": "0.1.0",
// 					"id":      "good-project",
// 				},
// 				"puppet_module": map[string]interface{}{
// 					"author":  u,
// 					"license": "Apache-2.0",
// 					"version": "0.2.0",
// 					"summary": "Output Override Summary",
// 				},
// 			},
// 		},
// 		{
// 			name: "with a non existant config, returns default config",
// 			args: args{
// 				info: pct.DeployInfo{
// 					TargetName: "good-project",
// 					PdkInfo: pct.PDKInfo{
// 						Version:   "0.1.0",
// 						Commit:    "abc12345",
// 						BuildDate: "2021/06/27",
// 					},
// 				},
// 				configFile:      "testdata/notthere/notthere/notthere.yml",
// 				projectTemplate: "",
// 				tmpl:            pct.PuppetContentTemplate{},
// 			},
// 			want: map[string]interface{}{
// 				"pct_name": "good-project",
// 				"user":     u,
// 				"cwd":      cwd,
// 				"hostname": hostName,
// 				"pdk": map[string]interface{}{
// 					"build_date":  "2021/06/27",
// 					"commit_hash": "abc12345",
// 					"version":     "0.1.0",
// 				},
// 				"puppet_module": map[string]interface{}{
// 					"author": u,
// 				},
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := pct.processConfiguration(tt.args.info, tt.args.configFile, tt.args.projectTemplate, tt.args.tmpl)
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("got = %+v\nwant %+v\n", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_readTemplateConfig(t *testing.T) {
// 	type args struct {
// 		configFile string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want pct.PuppetContentTemplateInfo
// 	}{
// 		{
// 			name: "returns tmpl struct from good config file",
// 			args: args{
// 				configFile: "testdata/examples/good-project/pct-config.yml",
// 			},
// 			want: pct.PuppetContentTemplateInfo{
// 				Template: pct.PuppetContentTemplate{
// 					Id:      "good-project",
// 					Display: "Good Project",
// 					Type:    "project",
// 					Version: "0.1.0",
// 					URL:     "https://github.com/puppetlabs/pct-good-project",
// 				},
// 				Defaults: map[string]interface{}{
// 					"puppet_module": map[string]interface{}{
// 						"license": "Apache-2.0",
// 						"version": "0.1.0",
// 						"summary": "A New Puppet Module",
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name: "returns empty struct from non-existant config file",
// 			args: args{
// 				configFile: "testdata/examples/does-not-exist-project/pct.yml",
// 			},
// 			want: pct.PuppetContentTemplateInfo{
// 				Template: PuppetContentTemplate{},
// 				Defaults: map[string]interface{}{},
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := pct.readTemplateConfig(tt.args.configFile); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("readTemplateConfig() = %+v, want %+v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_renderFile(t *testing.T) {
// 	type args struct {
// 		fileName string
// 		vars     interface{}
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want string
// 		err  bool
// 	}{
// 		{
// 			name: "takes a template file and returns correct text",
// 			args: args{
// 				fileName: "testdata/examples/good-project/content/goodfile.txt.tmpl",
// 				vars: map[string]interface{}{
// 					"example_data": "wakka",
// 				},
// 			},
// 			want: "This is wakka data",
// 			err:  false,
// 		},
// 		{
// 			name: "returns nil if file does not exist",
// 			args: args{
// 				fileName: "testdata/examples/non-existant-project/content/notthere.txt.tmpl",
// 				vars: map[string]interface{}{
// 					"example_data": "wakka",
// 				},
// 			},
// 			want: "",
// 			err:  true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := pct.renderFile(tt.args.fileName, tt.args.vars)
// 			if tt.err && err == nil {
// 				t.Fail()
// 			} else if !tt.err && got != tt.want {
// 				t.Errorf("renderFile() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
