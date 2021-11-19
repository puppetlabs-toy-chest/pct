package pct_install

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/puppetlabs/pdkgo/internal/pkg/pct"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

type PctInstall struct {
	AFS *afero.Afero
}

func (p *PctInstall) ProcessConfig(sourceDir, targetDir string, force bool) (string, error) {
	// Read config to determine template properties
	info, err := p.readConfig(filepath.Join(sourceDir, "pct-config.yml"))
	if err != nil {
		return "", fmt.Errorf("Invalid config: %v", err.Error())
	}

	// Create namespaced directory and move contents of temp folder to it
	namespacedPath, err := p.setupTemplateNamespace(targetDir, info, sourceDir, force)
	if err != nil {
		return "", fmt.Errorf("Unable to install in namespace: %v", err.Error())
	}
	return namespacedPath, nil
}

func (p *PctInstall) readConfig(configFile string) (info pct.PuppetContentTemplateInfo, err error) {
	fileBytes, err := p.AFS.ReadFile(configFile)
	if err != nil {
		return info, err
	}

	// use viper to parse the config as it knows how to work with mapstructure squash
	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewBuffer(fileBytes))
	if err != nil {
		return info, err
	}

	err = viper.Unmarshal(&info)
	if err != nil {
		return info, err
	}

	return info, err
}

func (p *PctInstall) setupTemplateNamespace(targetDir string, info pct.PuppetContentTemplateInfo, untarPath string, force bool) (string, error) {
	// author/id/version
	templatePath := filepath.Join(targetDir, info.Template.Author, info.Template.Id)

	err := p.AFS.MkdirAll(templatePath, 0750)
	if err != nil {
		return "", err
	}

	namespacePath := filepath.Join(targetDir, info.Template.Author, info.Template.Id, info.Template.Version)

	// finally move to the full path
	err = p.AFS.Rename(untarPath, namespacePath)
	// unable to check for this specific error as windows may instead return `access denied`
	// if err != nil && os.IsExist(err) {
	if err != nil {
		// if a template already exists
		if !force {
			// error unless forced
			return "", fmt.Errorf("Template already installed (%s)", namespacePath)
		} else {
			// remove the exiting template
			err = p.AFS.RemoveAll(namespacePath)
			if err != nil {
				return "", fmt.Errorf("Unable to overwrite existing template: %v", err)
			}
			// perform the move again
			err = p.AFS.Rename(untarPath, namespacePath)
			if err != nil {
				return "", fmt.Errorf("Unable to force install: %v", err)
			}
		}
	}

	return namespacePath, err
}
