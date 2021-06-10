package pct

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestMain(m *testing.M) {
	// hide logging output
	log.Logger = zerolog.New(ioutil.Discard).With().Timestamp().Logger()
	os.Exit(m.Run())
}

type osMock struct {
	GetwdFunc func() (dir string, err error)
}

func (osm osMock) Getwd() (dir string, err error) {
	return osm.GetwdFunc()
}

func TestDeploy(t *testing.T) {
	type args struct {
		info DeployInfo
	}
	tmp := t.TempDir()

	// this sets wrapper within pct.go
	osUtils = osMock{
		GetwdFunc: func() (dir string, err error) {
			return tmp, nil
		},
	}

	tmpFile := filepath.Base(tmp)

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "deploy a project and return the correct new files",
			args: args{
				info: DeployInfo{
					SelectedTemplate: "full-project",
					TemplateCache:    "testdata/examples",
					TargetOutputDir:  filepath.Join(tmp, "foobar"),
					TargetName:       "woo",
					PdkInfo: PDKInfo{
						Version:   "0.1.0",
						Commit:    "abc12345",
						BuildDate: "2021/06/27",
					},
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
				info: DeployInfo{
					SelectedTemplate: "full-project",
					TemplateCache:    "testdata/examples",
					TargetOutputDir:  "",
					TargetName:       "",
					PdkInfo: PDKInfo{
						Version:   "0.1.0",
						Commit:    "abc12345",
						BuildDate: "2021/06/27",
					},
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
				info: DeployInfo{
					SelectedTemplate: "full-project",
					TemplateCache:    "testdata/examples",
					TargetOutputDir:  "",
					TargetName:       "wibble",
					PdkInfo: PDKInfo{
						Version:   "0.1.0",
						Commit:    "abc12345",
						BuildDate: "2021/06/27",
					},
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
				info: DeployInfo{
					SelectedTemplate: "replace-thing",
					TemplateCache:    "testdata/examples",
					TargetOutputDir:  filepath.Join(tmp, "thing"),
					TargetName:       "woo",
					PdkInfo: PDKInfo{
						Version:   "0.1.0",
						Commit:    "abc12345",
						BuildDate: "2021/06/27",
					},
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
				info: DeployInfo{
					SelectedTemplate: "replace-thing",
					TemplateCache:    "testdata/examples",
					TargetOutputDir:  "",
					TargetName:       "",
					PdkInfo: PDKInfo{
						Version:   "0.1.0",
						Commit:    "abc12345",
						BuildDate: "2021/06/27",
					},
				},
			},
			want: []string{
				tmp,
				filepath.Join(tmp, tmpFile+".txt"),
			},
		},
		{
			name: "deploy a item without an outputDir",
			args: args{
				info: DeployInfo{
					SelectedTemplate: "replace-thing",
					TemplateCache:    "testdata/examples",
					TargetOutputDir:  "",
					TargetName:       "wibble",
					PdkInfo: PDKInfo{
						Version:   "0.1.0",
						Commit:    "abc12345",
						BuildDate: "2021/06/27",
					},
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
			if got := Deploy(tt.args.info); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Deploy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGet(t *testing.T) {
	type args struct {
		templateCache    string
		selectedTemplate string
	}
	tests := []struct {
		name    string
		args    args
		want    PuppetContentTemplate
		wantErr bool
	}{
		{
			name:    "returns error for non-existent template",
			args:    args{},
			wantErr: true,
		},
		{
			name: "returns tmpl for existent template",
			args: args{
				templateCache:    "testdata/examples",
				selectedTemplate: "full-project",
			},
			want: PuppetContentTemplate{
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
			got, err := Get(tt.args.templateCache, tt.args.selectedTemplate)
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

func Test_createTemplateFile(t *testing.T) {
	type args struct {
		info         DeployInfo
		configFile   string
		templateFile PuppetContentTemplateFileInfo
		tmpl         PuppetContentTemplate
	}

	tmp := t.TempDir()

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{
				info: DeployInfo{
					TargetName: "foobar",
					PdkInfo: PDKInfo{
						Version:   "0.1.0",
						Commit:    "abc12345",
						BuildDate: "2021/06/27",
					},
				},
				configFile: "testdata/examples/good-project/pct.yml",
				tmpl: PuppetContentTemplate{
					Type:    "project",
					Display: "Good Project",
					URL:     "https://github.com/puppetlabs/pct-good-project",
					Version: "0.1.0",
					Id:      "good-project",
				},
				templateFile: PuppetContentTemplateFileInfo{
					TemplatePath:   "testdata/examples/good-project/content/goodfile.txt.tmpl",
					TargetFilePath: filepath.Join(tmp, "foo.txt"),
					TargetDir:      tmp,
					TargetFile:     "",
					IsDirectory:    false,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createTemplateFile(tt.args.info, tt.args.configFile, tt.args.templateFile, tt.args.tmpl); (err != nil) != tt.wantErr {
				t.Errorf("createTemplateFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, err := os.Stat(tt.args.templateFile.TargetFilePath); err != nil {
				t.Errorf("createTemplateFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_processConfiguration(t *testing.T) {
	type args struct {
		info            DeployInfo
		configFile      string
		projectTemplate string
		tmpl            PuppetContentTemplate
	}
	cwd, _ := os.Getwd()
	hostName, _ := os.Hostname()
	u := getCurrentUser()
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "with a valid config, returns a correct map interface",
			args: args{
				info: DeployInfo{
					TargetName: "good-project",
					PdkInfo: PDKInfo{
						Version:   "0.1.0",
						Commit:    "abc12345",
						BuildDate: "2021/06/27",
					},
				},
				configFile:      "testdata/examples/good-project/pct.yml",
				projectTemplate: "",
				tmpl:            PuppetContentTemplate{},
			},
			want: map[string]interface{}{
				"user":     u,
				"cwd":      cwd,
				"hostname": hostName,
				"pct_name": "good-project",
				"pdk": map[string]interface{}{
					"build_date":  "2021/06/27",
					"commit_hash": "abc12345",
					"version":     "0.1.0",
				},
				"template": map[string]interface{}{
					"type":    "project",
					"display": "Good Project",
					"url":     "https://github.com/puppetlabs/pct-good-project",
					"version": "0.1.0",
					"id":      "good-project",
				},
				"puppet_module": map[string]interface{}{
					"author":  u,
					"license": "Apache-2.0",
					"version": "0.1.0",
					"summary": "A New Puppet Module",
				},
			},
		},
		{
			name: "with a valid config, and a workspace overide, returns a correct map interface",
			args: args{
				info: DeployInfo{
					TargetName:      "good-project",
					TargetOutputDir: "testdata/examples/good-project-override",
					PdkInfo: PDKInfo{
						Version:   "0.1.0",
						Commit:    "abc12345",
						BuildDate: "2021/06/27",
					},
				},
				configFile:      "testdata/examples/good-project/pct.yml",
				projectTemplate: "",
				tmpl:            PuppetContentTemplate{},
			},
			want: map[string]interface{}{
				"user":     u,
				"cwd":      cwd,
				"hostname": hostName,
				"pct_name": "good-project",
				"pdk": map[string]interface{}{
					"build_date":  "2021/06/27",
					"commit_hash": "abc12345",
					"version":     "0.1.0",
				},
				"template": map[string]interface{}{
					"type":    "project",
					"display": "Good Project",
					"url":     "https://github.com/puppetlabs/pct-good-project",
					"version": "0.1.0",
					"id":      "good-project",
				},
				"puppet_module": map[string]interface{}{
					"author":  u,
					"license": "Apache-2.0",
					"version": "0.2.0",
					"summary": "Output Override Summary",
				},
			},
		},
		{
			name: "with a non existant config, returns default config",
			args: args{
				info: DeployInfo{
					TargetName: "good-project",
					PdkInfo: PDKInfo{
						Version:   "0.1.0",
						Commit:    "abc12345",
						BuildDate: "2021/06/27",
					},
				},
				configFile:      "testdata/notthere/notthere/notthere.yml",
				projectTemplate: "",
				tmpl:            PuppetContentTemplate{},
			},
			want: map[string]interface{}{
				"pct_name": "good-project",
				"user":     u,
				"cwd":      cwd,
				"hostname": hostName,
				"pdk": map[string]interface{}{
					"build_date":  "2021/06/27",
					"commit_hash": "abc12345",
					"version":     "0.1.0",
				},
				"puppet_module": map[string]interface{}{
					"author": u,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processConfiguration(tt.args.info, tt.args.configFile, tt.args.projectTemplate, tt.args.tmpl)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %+v\nwant %+v\n", got, tt.want)
			}
		})
	}
}

func Test_readTemplateConfig(t *testing.T) {
	type args struct {
		configFile string
	}
	tests := []struct {
		name string
		args args
		want PuppetContentTemplateInfo
	}{
		{
			name: "returns tmpl struct from good config file",
			args: args{
				configFile: "testdata/examples/good-project/pct-config.yml",
			},
			want: PuppetContentTemplateInfo{
				Template: PuppetContentTemplate{
					Id:      "good-project",
					Display: "Good Project",
					Type:    "project",
					Version: "0.1.0",
					URL:     "https://github.com/puppetlabs/pct-good-project",
				},
				Defaults: map[string]interface{}{
					"puppet_module": map[string]interface{}{
						"license": "Apache-2.0",
						"version": "0.1.0",
						"summary": "A New Puppet Module",
					},
				},
			},
		},
		{
			name: "returns empty struct from non-existant config file",
			args: args{
				configFile: "testdata/examples/does-not-exist-project/pct.yml",
			},
			want: PuppetContentTemplateInfo{
				Template: PuppetContentTemplate{},
				Defaults: map[string]interface{}{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readTemplateConfig(tt.args.configFile); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readTemplateConfig() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_renderFile(t *testing.T) {
	type args struct {
		fileName string
		vars     interface{}
	}
	tests := []struct {
		name string
		args args
		want string
		err  bool
	}{
		{
			name: "takes a template file and returns correct text",
			args: args{
				fileName: "testdata/examples/good-project/content/goodfile.txt.tmpl",
				vars: map[string]interface{}{
					"example_data": "wakka",
				},
			},
			want: "This is wakka data",
			err:  false,
		},
		{
			name: "returns nil if file does not exist",
			args: args{
				fileName: "testdata/examples/non-existant-project/content/notthere.txt.tmpl",
				vars: map[string]interface{}{
					"example_data": "wakka",
				},
			},
			want: "",
			err:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := renderFile(tt.args.fileName, tt.args.vars)
			if tt.err && err == nil {
				t.Fail()
			} else if !tt.err && got != tt.want {
				t.Errorf("renderFile() = %v, want %v", got, tt.want)
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
			returnString := DisplayDefaults(tt.args.defaults, tt.format)

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
