/*
Package pct implements the Puppet Content template specification

Puppet Content Templates (PCT) codify a structure to produce content for any Puppet
Product. PCT can create any type of a Puppet Product project: Puppet control
repo, Puppet Module, Bolt project, etc. PCT can also create one or more independent
files, such as CI files or gitignores. This can be as simple as a name for a
Puppet Class or a set of CI files to add to a Puppet Module.
*/
package pct

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"text/template"

	jsoniter "github.com/json-iterator/go"
	"github.com/mitchellh/go-homedir"
	"github.com/olekukonko/tablewriter"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const (
	TemplateConfigName         = "pct-config"
	TemplateConfigFileName     = "pct-config.yml"
	UserTemplateConfigName     = "pct"
	UserTemplateConfigFileName = "pct.yml"
)

// PuppetContentTemplateInfo is the housing struct for marshaling YAML data
type PuppetContentTemplateInfo struct {
	Template PuppetContentTemplate `mapstructure:"template"`
	Defaults map[string]interface{}
}

// PuppetContentTemplate houses the actual information about each template
type PuppetContentTemplate struct {
	Id      string `mapstructure:"id"`
	Type    string `mapstructure:"type"`
	Display string `mapstructure:"display"`
	Version string `mapstructure:"version"`
	URL     string `mapstructure:"url"`
}

// PuppetContentTemplateFileInfo represents the resolved target path information
// for a given template
type PuppetContentTemplateFileInfo struct {
	TemplatePath   string
	TargetFilePath string
	TargetDir      string
	TargetFile     string
	IsDirectory    bool
}

// PDKInfo contains the current version information of the compiled binary for
// use in template data
type PDKInfo struct {
	Version   string
	Commit    string
	BuildDate string
}

// DeployInfo represents the final set of information needed to deploy a template
type DeployInfo struct {
	SelectedTemplate string
	TemplateCache    string
	TargetOutputDir  string
	TargetName       string
	PdkInfo          PDKInfo
}

type osWrapper interface {
	Getwd() (string, error)
}
type osFunc struct {
}

func (osf osFunc) Getwd() (dir string, err error) {
	return os.Getwd()
}

var osUtils osWrapper = osFunc{}

func Get(templateCache string, selectedTemplate string) (PuppetContentTemplate, error) {
	info, err := GetInfo(templateCache, selectedTemplate)
	return info.Template, err
}

func GetInfo(templateCache string, selectedTemplate string) (PuppetContentTemplateInfo, error) {
	file := filepath.Join(templateCache, selectedTemplate, TemplateConfigFileName)
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return PuppetContentTemplateInfo{}, fmt.Errorf("Couldn't find an installed template that matches '%s'", selectedTemplate)
	}
	i := readTemplateConfig(file)
	return i, nil
}

// List lists all templates in a given path and parses their configuration. Does
// not return any errors from parsing invalid templates, but returns them as
// debug log events
func List(templatePath string, templateName string) ([]PuppetContentTemplate, error) {
	log.Debug().Msgf("Searching %+v for templates", templatePath)
	matches, _ := filepath.Glob(templatePath + "/**/" + TemplateConfigFileName)

	var tmpls []PuppetContentTemplate
	for _, file := range matches {
		log.Debug().Msgf("Found: %+v", file)
		i := readTemplateConfig(file).Template
		tmpls = append(tmpls, i)
	}

	if templateName != "" {
		log.Debug().Msgf("Filtering for: %s", templateName)
		tmpls = filterFiles(tmpls, func(f PuppetContentTemplate) bool { return f.Id == templateName })
	}

	return tmpls, nil
}

// FormatTemplates formats one or more templates to display on the console in
// table format or json format.
func FormatTemplates(tmpls []PuppetContentTemplate, jsonOutput string) error {
	switch jsonOutput {
	case "table":
		fmt.Println("")

		count := len(tmpls)
		if count < 1 {
			log.Warn().Msgf("Could not locate any templates at %+v", viper.GetString("templatepath"))
		} else if count == 1 {
			fmt.Printf("DisplayName:     %v\n", tmpls[0].Display)
			fmt.Printf("Name:            %v\n", tmpls[0].Id)
			fmt.Printf("TemplateType:    %v\n", tmpls[0].Type)
			fmt.Printf("TemplateURL:     %v\n", tmpls[0].URL)
			fmt.Printf("TemplateVersion: %v\n", tmpls[0].Version)
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"DisplayName", "Name", "Type"})
			table.SetBorder(false)
			for _, v := range tmpls {
				table.Append([]string{v.Display, v.Id, v.Type})
			}
			table.Render()
		}
	case "json":
		j := jsoniter.ConfigFastest
		prettyJSON, err := j.MarshalIndent(&tmpls, "", "  ")
		if err != nil {
			log.Error().Msgf("Error converting to json: %v", err)
		}
		fmt.Printf("%s\n", string(prettyJSON))
	}
	return nil
}

// DisplayDefaults returns the Default values from a Templates configuration file
func DisplayDefaults(defaults map[string]interface{}, format string) string {
	var err error
	var prettyDefaults []byte

	switch format {
	case "table":
		if len(defaults) == 0 {
			return "This template has no configuration options."
		}

		prettyDefaults, err = yaml.Marshal(defaults)
		if err != nil {
			log.Error().Msgf("Error converting to yaml: %v", err)
		}
	case "json":
		j := jsoniter.ConfigFastest
		prettyDefaults, err = j.MarshalIndent(defaults, "", "  ")
		if err != nil {
			log.Error().Msgf("Error converting to json: %v", err)
		}
	}

	return string(prettyDefaults)
}

// FormatDeployment formats the files returned by the Deploy method to display
// on the console in table format or json format.
func FormatDeployment(deployed []string, jsonOutput string) error {
	switch jsonOutput {
	case "table":
		for _, d := range deployed {
			log.Info().Msgf("Deployed: %v", d)
		}
	case "json":
		j := jsoniter.ConfigFastest
		prettyJSON, _ := j.MarshalIndent(deployed, "", "  ")
		fmt.Printf("%s\n", prettyJSON)
	}
	return nil
}

// Deploy deploys a selected template to a target path with a target name using
// data from both the configuration inside the template and provided by the
// User in their user config file
func Deploy(info DeployInfo) []string {

	log.Trace().Msgf("PDKInfo: %+v", info.PdkInfo)

	file := filepath.Join(info.TemplateCache, info.SelectedTemplate, TemplateConfigFileName)
	log.Debug().Msgf("Template: %s", file)
	tmpl := readTemplateConfig(file)
	log.Trace().Msgf("Parsed: %+v", tmpl)

	if info.TargetName == "" && info.TargetOutputDir == "" { // pdk new foo-foo
		cwd, _ := osUtils.Getwd()
		info.TargetName = filepath.Base(cwd)
		info.TargetOutputDir = cwd
	} else if info.TargetName != "" && info.TargetOutputDir == "" { // pdk new foo-foo -n wakka
		cwd, _ := osUtils.Getwd()
		if tmpl.Template.Type == "project" {
			info.TargetOutputDir = filepath.Join(cwd, info.TargetName)
		} else {
			info.TargetOutputDir = cwd
		}
	} else if info.TargetName == "" && info.TargetOutputDir != "" { // pdk new foo-foo -o /foo/bar/baz
		info.TargetName = filepath.Base(info.TargetOutputDir)
	} else if info.TargetName != "" && info.TargetOutputDir != "" { // pdk new foo-foo -n wakka -o /foo/bar/baz
		if tmpl.Template.Type == "project" {
			info.TargetOutputDir = filepath.Join(info.TargetOutputDir, info.TargetName)
		}
	}

	contentDir := filepath.Join(info.TemplateCache, info.SelectedTemplate, "content")
	log.Debug().Msgf("Target Name: %s", info.TargetName)
	log.Debug().Msgf("Target Output: %s", info.TargetOutputDir)

	replacer := strings.NewReplacer(
		contentDir, info.TargetOutputDir,
		"{{pct_name}}", info.TargetName,
		".tmpl", "",
	)

	var templateFiles []PuppetContentTemplateFileInfo
	err := filepath.WalkDir(contentDir, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		log.Trace().Msgf("Processing: %s", path)
		targetFile := replacer.Replace(path)
		log.Debug().Msgf("Resolved '%s' to '%s'", path, targetFile)

		dir, file := filepath.Split(targetFile)
		i := PuppetContentTemplateFileInfo{
			TemplatePath:   path,
			TargetFilePath: targetFile,
			TargetDir:      dir,
			TargetFile:     file,
			IsDirectory:    info.IsDir(),
		}
		log.Trace().Msgf("Processed: %+v", i)

		templateFiles = append(templateFiles, i)
		return nil
	})
	if err != nil {
		log.Error().AnErr("content", err)
	}

	var deployed []string
	for _, templateFile := range templateFiles {
		log.Debug().Msgf("Deploying: %s", templateFile.TargetFilePath)
		if templateFile.IsDirectory {
			err := createTemplateDirectory(templateFile.TargetFilePath)
			if err == nil {
				deployed = append(deployed, templateFile.TargetFilePath)
			}
		} else {
			err := createTemplateFile(info, file, templateFile, tmpl.Template)
			if err != nil {
				log.Error().Msgf("%s", err)
				continue
			}
			deployed = append(deployed, templateFile.TargetFilePath)
		}
	}

	return deployed
}

func createTemplateDirectory(targetDir string) error {
	log.Trace().Msgf("Creating: '%s'", targetDir)
	err := os.MkdirAll(targetDir, os.ModePerm)

	if err != nil {
		log.Error().Msgf("Error: %v", err)
		return err
	}

	return nil
}

func createTemplateFile(info DeployInfo, configFile string, templateFile PuppetContentTemplateFileInfo, tmpl PuppetContentTemplate) error {
	log.Trace().Msgf("Creating: '%s'", templateFile.TargetFilePath)
	config := processConfiguration(
		info,
		configFile,
		templateFile.TemplatePath,
		tmpl,
	)

	text, err := renderFile(templateFile.TemplatePath, config)
	if err != nil {
		return fmt.Errorf("Failed to create %s", templateFile.TargetFilePath)
	}

	log.Trace().Msgf("Writing: '%s' '%s'", templateFile.TargetFilePath, text)
	err = os.MkdirAll(templateFile.TargetDir, os.ModePerm)
	if err != nil {
		log.Error().Msgf("Error: %v", err)
		return err
	}

	file, err := os.Create(templateFile.TargetFilePath)
	if err != nil {
		log.Error().Msgf("Error: %v", err)
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Error().Msgf("Error closing file: %s\n", err)
		}
	}()

	_, err = io.WriteString(file, text)
	if err != nil {
		log.Error().Msgf("Error: %v", err)
		return err
	}

	err = file.Sync()
	if err != nil {
		log.Error().Msgf("Error: %v", err)
		return err
	}

	return nil
}

func processConfiguration(info DeployInfo, configFile string, projectTemplate string, tmpl PuppetContentTemplate) map[string]interface{} {
	v := viper.New()

	log.Trace().Msgf("PDKInfo: %+v", info.PdkInfo)
	/*
		Inheritance (each level overwritten by next):
			convention based variables
				- pdk specific variables based on transformed user input
			machine variables
				- information that comes from the current machine
				- user name, hostname, etc
			template variables
				- information from the template itself
				- designed to be runnable defaults for everything inside template
			user overrides
				- ~/.pdk/pct.yml
				- user customizations for their preferences
			Workspace overrides
			  - ${cwd}/pct.yml
				- ${outputDir}/pct.yml
	*/

	// Convention based variables
	v.SetDefault("pct_name", info.TargetName)

	user := getCurrentUser()
	v.SetDefault("user", user)
	v.SetDefault("puppet_module.author", user)

	// Machine based variables
	cwd, _ := os.Getwd()
	hostName, _ := os.Hostname()
	v.SetDefault("cwd", cwd)
	v.SetDefault("hostname", hostName)

	// PDK binary specific variables
	v.SetDefault("pdk.version", info.PdkInfo.Version)
	v.SetDefault("pdk.commit_hash", info.PdkInfo.Commit)
	v.SetDefault("pdk.build_date", info.PdkInfo.BuildDate)

	// Template specific variables
	log.Trace().Msgf("Adding %v", filepath.Dir(configFile))
	v.SetConfigName(TemplateConfigName)
	v.SetConfigType("yml")
	v.AddConfigPath(filepath.Dir(configFile))
	if err := v.ReadInConfig(); err == nil {
		log.Trace().Msgf("Merging config file: %v", v.ConfigFileUsed())
	} else {
		log.Error().Msgf("Error reading config: %v", err)
	}

	// User specified variable overrides
	home, _ := homedir.Dir()
	userConfigPath := filepath.Join(home, ".pdk")
	log.Trace().Msgf("Adding %v", userConfigPath)
	vUser := viper.New()
	vUser.SetConfigName(UserTemplateConfigName)
	vUser.SetConfigType("yml")
	vUser.AddConfigPath(userConfigPath)
	if err := vUser.ReadInConfig(); err == nil {
		log.Trace().Msgf("Merging config file: %v", v.ConfigFileUsed())
	} else {
		log.Debug().Msgf("Error reading config: %v", err)
	}

	// workspace overrides
	vWorkspace := viper.New()
	vWorkspace.SetConfigName(UserTemplateConfigName)
	vWorkspace.SetConfigType("yml")
	vWorkspace.AddConfigPath(info.TargetOutputDir)
	if err := vWorkspace.ReadInConfig(); err == nil {
		log.Trace().Msgf("Merging config file: %v", v.ConfigFileUsed())
	} else {
		log.Debug().Msgf("Error reading config: %v", err)
	}

	userMergeErr := v.MergeConfigMap(vUser.AllSettings())
	if userMergeErr != nil {
		log.Warn().Msgf("Unable to merge user configuration values: %s", userMergeErr.Error())
	}
	workspaceMergeErr := v.MergeConfigMap(vWorkspace.AllSettings())
	if userMergeErr != nil {
		log.Warn().Msgf("Unable to merge workspace configuration values: %s", workspaceMergeErr.Error())
	}

	config := make(map[string]interface{})
	err := v.Unmarshal(&config)
	if err != nil {
		log.Error().Msgf("unable to decode into struct, %v", err)
		return nil
	}

	return config
}

func readTemplateConfig(configFile string) PuppetContentTemplateInfo {
	v := viper.New()
	userConfigFileBase := filepath.Base(configFile)
	v.AddConfigPath(filepath.Dir(configFile))
	v.SetConfigName(userConfigFileBase)
	v.SetConfigType("yml")
	if err := v.ReadInConfig(); err == nil {
		log.Trace().Msgf("Using template config file: %v", v.ConfigFileUsed())
	}
	var config PuppetContentTemplateInfo
	// unmarshall the known structure
	err := v.Unmarshal(&config)
	if err != nil {
		log.Error().Msgf("unable to decode into struct, %v", err)
	}

	// unmarhsall everything
	all := make(map[string]interface{})
	err = v.Unmarshal(&all)
	if err != nil {
		log.Error().Msgf("unable to decode into struct, %v", err)
	}
	// remove the known structure, leaving the unknown...
	delete(all, "template")
	// store the unknown as part of the big config
	config.Defaults = all

	return config
}

func renderFile(fileName string, vars interface{}) (string, error) {
	tmpl, err := template.
		New(filepath.Base(fileName)).
		Funcs(
			template.FuncMap{
				"toClassName": func(itemName string) string {
					return strings.Title(strings.ToLower(itemName))
				},
			},
		).
		ParseFiles(fileName)

	if err != nil {
		log.Error().Msgf("Error parsing config: %v", err)
		return "", err
	}

	return process(tmpl, vars), nil
}

func process(t *template.Template, vars interface{}) string {
	var tmplBytes bytes.Buffer

	err := t.Execute(&tmplBytes, vars)
	if err != nil {
		log.Error().Msgf("Error parsing config: %v", err)
		return ""
	}
	return tmplBytes.String()
}

func filterFiles(ss []PuppetContentTemplate, test func(PuppetContentTemplate) bool) (ret []PuppetContentTemplate) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func getCurrentUser() string {
	user, _ := user.Current()
	if strings.Contains(user.Username, "\\") {
		v := strings.Split(user.Username, "\\")
		return v[1]
	}
	return user.Username
}
