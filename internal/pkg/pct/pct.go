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

func Get(templateCache string, selectedTemplate string) (PuppetContentTemplate, error) {
	file := filepath.Join(templateCache, selectedTemplate, TemplateConfigFileName)
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return PuppetContentTemplate{}, fmt.Errorf("Couldn't find an installed template that matches '%s'", selectedTemplate)
	}
	i := readTemplateConfig(file)
	return i, nil
}

// List lists all templates in a given path and parses their configuration. Does
// not return any errors from parsing invalid templates, but returns them as
// debug log events
func List(templatePath string, templateName string) ([]PuppetContentTemplate, error) {
	matches, _ := filepath.Glob(templatePath + "/**/" + TemplateConfigFileName)
	var tmpls []PuppetContentTemplate
	for _, file := range matches {
		log.Debug().Msgf("Found: %+v", file)
		i := readTemplateConfig(file)
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
		if len(tmpls) == 1 {
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

	// pdk new foo-foo
	if info.TargetName == "" && info.TargetOutputDir == "" {
		cwd, _ := os.Getwd()
		info.TargetName = filepath.Base(cwd)
		info.TargetOutputDir = cwd
	}

	// pdk new foo-foo -n wakka
	if info.TargetName != "" && info.TargetOutputDir == "" {
		cwd, _ := os.Getwd()
		info.TargetOutputDir = filepath.Join(cwd, info.TargetName)
	}

	// pdk new foo-foo -o /foo/bar/baz
	if info.TargetName == "" && info.TargetOutputDir != "" {
		info.TargetName = filepath.Base(info.TargetOutputDir)
	}

	// pdk new foo-foo
	if info.TargetName == "" {
		cwd, _ := os.Getwd()
		info.TargetName = filepath.Base(cwd)
	}

	// pdk new foo-foo
	// pdk new foo-foo -n wakka
	// pdk new foo-foo -n wakka -o c:/foo
	// pdk new foo-foo -n wakka -o c:/foo/wakka
	switch tmpl.Type {
	case "project":
		if info.TargetOutputDir == "" {
			cwd, _ := os.Getwd()
			info.TargetOutputDir = cwd
		} else if strings.HasSuffix(info.TargetOutputDir, info.TargetName) {
			// user has specified outputpath with the info.Targetname in it
		} else {
			info.TargetOutputDir = filepath.Join(info.TargetOutputDir, info.TargetName)
		}
	case "item":
		if info.TargetOutputDir == "" {
			cwd, _ := os.Getwd()
			info.TargetOutputDir = cwd
		} else if strings.HasSuffix(info.TargetOutputDir, info.TargetName) {
			// user has specified outputpath with the info.Targetname in it
			info.TargetOutputDir, _ = filepath.Split(info.TargetOutputDir)
			log.Debug().Msgf("Changing target to :%s", info.TargetOutputDir)
			info.TargetOutputDir = filepath.Clean(info.TargetOutputDir)
			log.Debug().Msgf("Changing target to :%s", info.TargetOutputDir)
		}
		// } else {
		// 	// use what the user tells us
		// }

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
			err := createTemplateFile(info.TargetName, file, templateFile, tmpl, info.PdkInfo)
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

func createTemplateFile(targetName string, configFile string, templateFile PuppetContentTemplateFileInfo, tmpl PuppetContentTemplate, pdkInfo PDKInfo) error {
	log.Trace().Msgf("Creating: '%s'", templateFile.TargetFilePath)
	config := processConfiguration(
		targetName,
		configFile,
		templateFile.TemplatePath,
		tmpl,
		pdkInfo,
	)

	text := renderFile(templateFile.TemplatePath, config)
	if text == "" {
		return fmt.Errorf("Failed to create %s", templateFile.TargetFilePath)
	}

	log.Trace().Msgf("Writing: '%s' '%s'", templateFile.TargetFilePath, text)
	err := os.MkdirAll(templateFile.TargetDir, os.ModePerm)
	if err != nil {
		log.Error().Msgf("Error: %v", err)
		return err
	}

	file, err := os.Create(templateFile.TargetFilePath)
	if err != nil {
		log.Error().Msgf("Error: %v", err)
		return err
	}
	defer file.Close()

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

func processConfiguration(projectName string, configFile string, projectTemplate string, tmpl PuppetContentTemplate, pdkInfo PDKInfo) map[string]interface{} {
	v := viper.New()

	log.Trace().Msgf("PDKInfo: %+v", pdkInfo)
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
				- ~/.pdk/pdk.yml
				- user customizations for their preferences
	*/

	// Convention based variables
	v.SetDefault("pct_name", projectName)

	user := getCurrentUser()
	v.SetDefault("user", user)
	v.SetDefault("puppet_module.author", user)

	// Machine based variables
	cwd, _ := os.Getwd()
	hostName, _ := os.Hostname()
	v.SetDefault("cwd", cwd)
	v.SetDefault("hostname", hostName)

	// PDK binary specific variables
	v.SetDefault("pdk.version", pdkInfo.Version)
	v.SetDefault("pdk.commit_hash", pdkInfo.Commit)
	v.SetDefault("pdk.build_date", pdkInfo.BuildDate)

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
	v.SetConfigName(UserTemplateConfigName)
	v.SetConfigType("yml")
	v.AddConfigPath(userConfigPath)
	if err := v.MergeInConfig(); err == nil {
		log.Trace().Msgf("Merging config file: %v", v.ConfigFileUsed())
	} else {
		log.Debug().Msgf("Error reading config: %v", err)
	}

	config := make(map[string]interface{})
	err := v.Unmarshal(&config)
	if err != nil {
		log.Error().Msgf("unable to decode into struct, %v", err)
		return nil
	}

	return config
}

func readTemplateConfig(configFile string) PuppetContentTemplate {
	v := viper.New()
	userConfigFileBase := filepath.Base(configFile)
	v.AddConfigPath(filepath.Dir(configFile))
	v.SetConfigName(userConfigFileBase)
	v.SetConfigType("yml")
	if err := v.ReadInConfig(); err == nil {
		log.Trace().Msgf("Using template config file: %v", v.ConfigFileUsed())
	}
	var config PuppetContentTemplateInfo
	err := v.Unmarshal(&config)
	if err != nil {
		log.Error().Msgf("unable to decode into struct, %v", err)
	}
	return config.Template
}

func renderFile(fileName string, vars interface{}) string {
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
		return ""
	}

	return process(tmpl, vars)
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
