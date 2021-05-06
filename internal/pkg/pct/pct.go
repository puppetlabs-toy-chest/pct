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
	TemplateConfigName     = "pct"
	TemplateConfigFileName = "pct.yml"
)

type PDKInfo struct {
	Version   string
	Commit    string
	BuildDate string
}

type PuppetContentTemplateFileInfo struct {
	TemplatePath   string
	TargetFilePath string
	TargetDir      string
	TargetFile     string
	IsDirectory    bool
}

type PuppetContentTemplate struct {
	Template PuppetContentTemplateInfo `mapstructure:"template"`
}
type PuppetContentTemplateInfo struct {
	Name    string `mapstructure:"name"`
	Type    string `mapstructure:"type"`
	Display string `mapstructure:"display"`
	Version string `mapstructure:"version"`
	URL     string `mapstructure:"url"`
}

func List(templatePath string, templateName string) ([]PuppetContentTemplateInfo, error) {
	matches, _ := filepath.Glob(templatePath + "/**/" + TemplateConfigFileName)
	var tmpls []PuppetContentTemplateInfo
	for _, file := range matches {
		log.Debug().Msgf("Found: %+v", file)
		i := readTemplateConfig(file)
		tmpls = append(tmpls, i)
	}

	if templateName != "" {
		log.Debug().Msgf("Filtering for: %s", templateName)
		tmpls = filterFiles(tmpls, func(f PuppetContentTemplateInfo) bool { return f.Name == templateName })
	}

	return tmpls, nil
}

func FormatTemplates(tmpls []PuppetContentTemplateInfo, jsonOutput bool) error {
	if jsonOutput {
		j := jsoniter.ConfigFastest
		prettyJSON, err := j.MarshalIndent(&tmpls, "", "  ")
		if err != nil {
			log.Error().Msgf("Error converting to json: %v", err)
		}
		fmt.Printf("%s\n", string(prettyJSON))
	} else {
		fmt.Println("")
		if len(tmpls) == 1 {
			fmt.Printf("DisplayName:     %v\n", tmpls[0].Display)
			fmt.Printf("Name:            %v\n", tmpls[0].Name)
			fmt.Printf("TemplateType:    %v\n", tmpls[0].Type)
			fmt.Printf("TemplateURL:     %v\n", tmpls[0].URL)
			fmt.Printf("TemplateVersion: %v\n", tmpls[0].Version)
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"DisplayName", "Name", "Type"})
			table.SetBorder(false)
			for _, v := range tmpls {
				table.Append([]string{v.Display, v.Name, v.Type})
			}
			table.Render()
		}
	}

	return nil
}

func FormatDeployment(deployed []string, jsonOutput bool) error {
	if jsonOutput {
		j := jsoniter.ConfigFastest
		prettyJSON, _ := j.MarshalIndent(deployed, "", "  ")
		fmt.Printf("%s\n", prettyJSON)
	} else {
		for _, d := range deployed {
			log.Info().Msgf("Deployed: %v", d)
		}
	}

	return nil
}

func Deploy(selectedTemplate string, localTemplateCache string, targetOutput string, targetName string, pdkInfo PDKInfo) []string {

	log.Trace().Msgf("PDKInfo: %+v", pdkInfo)

	file := filepath.Join(localTemplateCache, selectedTemplate, TemplateConfigFileName)
	log.Debug().Msgf("Template: %s", file)
	tmpl := readTemplateConfig(file)
	log.Trace().Msgf("Parsed: %+v", tmpl)

	// pdk new foo-foo
	if targetName == "" && targetOutput == "" {
		cwd, _ := os.Getwd()
		targetName = filepath.Base(cwd)
		targetOutput = cwd
	}

	// pdk new foo-foo -n wakka
	if targetName != "" && targetOutput == "" {
		cwd, _ := os.Getwd()
		targetOutput = filepath.Join(cwd, targetName)
	}

	// pdk new foo-foo -o /foo/bar/baz
	if targetName == "" && targetOutput != "" {
		targetName = filepath.Base(targetOutput)
	}

	// pdk new foo-foo
	if targetName == "" {
		cwd, _ := os.Getwd()
		targetName = filepath.Base(cwd)
	}

	// pdk new foo-foo
	// pdk new foo-foo -n wakka
	// pdk new foo-foo -n wakka -o c:/foo
	// pdk new foo-foo -n wakka -o c:/foo/wakka
	switch tmpl.Type {
	case "project":
		if targetOutput == "" {
			cwd, _ := os.Getwd()
			targetOutput = cwd
		} else if strings.HasSuffix(targetOutput, targetName) {
			// user has specified outputpath with the targetname in it
		} else {
			targetOutput = filepath.Join(targetOutput, targetName)
		}
	case "item":
		if targetOutput == "" {
			cwd, _ := os.Getwd()
			targetOutput = cwd
		} else if strings.HasSuffix(targetOutput, targetName) {
			// user has specified outputpath with the targetname in it
			targetOutput, _ = filepath.Split(targetOutput)
			log.Debug().Msgf("Changing target to :%s", targetOutput)
			targetOutput = filepath.Clean(targetOutput)
			log.Debug().Msgf("Changing target to :%s", targetOutput)
		}
		// } else {
		// 	// use what the user tells us
		// }

	}

	contentDir := filepath.Join(localTemplateCache, selectedTemplate, "content")
	log.Debug().Msgf("Target Name: %s", targetName)
	log.Debug().Msgf("Target Output: %s", targetOutput)

	var templateFiles []PuppetContentTemplateFileInfo
	err := filepath.WalkDir(contentDir, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		log.Trace().Msgf("Processing: %s", path)

		replacer := strings.NewReplacer(
			contentDir, targetOutput,
			"__REPLACE__", targetName,
			".tmpl", "",
		)
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
			err := createTemplateFile(targetName, file, templateFile, tmpl, pdkInfo)
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

func createTemplateFile(targetName string, configFile string, templateFile PuppetContentTemplateFileInfo, tmpl PuppetContentTemplateInfo, pdkInfo PDKInfo) error {
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

func processConfiguration(projectName string, configFile string, projectTemplate string, tmpl PuppetContentTemplateInfo, pdkInfo PDKInfo) map[string]interface{} {
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
	switch tmpl.Type {
	case "project":
		v.SetDefault("project_name", projectName)
	case "item":
		v.SetDefault("item_name", projectName)
	}
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
	v.SetConfigName("pdk")
	v.SetConfigType("yml")
	v.AddConfigPath(userConfigPath)
	if err := v.MergeInConfig(); err == nil {
		log.Trace().Msgf("Merging config file: %v", v.ConfigFileUsed())
	} else {
		log.Error().Msgf("Error reading config: %v", err)
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
	var config PuppetContentTemplate
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

func filterFiles(ss []PuppetContentTemplateInfo, test func(PuppetContentTemplateInfo) bool) (ret []PuppetContentTemplateInfo) {
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
