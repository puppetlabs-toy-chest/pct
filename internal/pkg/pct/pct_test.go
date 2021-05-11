package pct

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDeploy(t *testing.T) {
	type args struct {
		info DeployInfo
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Deploy(tt.args.info); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Deploy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createTemplateFile(t *testing.T) {
	type args struct {
		targetName   string
		configFile   string
		templateFile PuppetContentTemplateFileInfo
		tmpl         PuppetContentTemplate
		pdkInfo      PDKInfo
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
				targetName: "foobar",
				configFile: "testdata/examples/good-project/pct.yml",
				pdkInfo: PDKInfo{
					Version:   "0.1.0",
					Commit:    "abc12345",
					BuildDate: "2021/06/27",
				},
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
			if err := createTemplateFile(tt.args.targetName, tt.args.configFile, tt.args.templateFile, tt.args.tmpl, tt.args.pdkInfo); (err != nil) != tt.wantErr {
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
		projectName     string
		configFile      string
		projectTemplate string
		tmpl            PuppetContentTemplate
		pdkInfo         PDKInfo
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
				projectName:     "good-project",
				configFile:      "testdata/examples/good-project/pct.yml",
				projectTemplate: "",
				tmpl:            PuppetContentTemplate{},
				pdkInfo: PDKInfo{
					Version:   "0.1.0",
					Commit:    "abc12345",
					BuildDate: "2021/06/27",
				},
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
			name: "with a non existant config, returns default config",
			args: args{
				projectName:     "good-project",
				configFile:      "testdata/notthere/notthere/notthere.yml",
				projectTemplate: "",
				tmpl:            PuppetContentTemplate{},
				pdkInfo: PDKInfo{
					Version:   "0.1.0",
					Commit:    "abc12345",
					BuildDate: "2021/06/27",
				},
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
			got := processConfiguration(tt.args.projectName, tt.args.configFile, tt.args.projectTemplate, tt.args.tmpl, tt.args.pdkInfo)
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
		want PuppetContentTemplate
	}{
		{
			name: "returns tmpl struct from good config file",
			args: args{
				configFile: "testdata/examples/good-project/pct-config.yml",
			},
			want: PuppetContentTemplate{
				Id:      "good-project",
				Display: "Good Project",
				Type:    "project",
				Version: "0.1.0",
				URL:     "https://github.com/puppetlabs/pct-good-project",
			},
		},
		{
			name: "returns empty struct from non-existant config file",
			args: args{
				configFile: "testdata/examples/does-not-exist-project/pct.yml",
			},
			want: PuppetContentTemplate{},
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := renderFile(tt.args.fileName, tt.args.vars); got != tt.want {
				t.Errorf("renderFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
