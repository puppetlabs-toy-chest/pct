package pct_config_processor

import (
	"bytes"
	"fmt"

	"github.com/puppetlabs/pctkg/config_processor"
	"github.com/puppetlabs/pctnternal/pkg/pct"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

type PctConfigProcessor struct {
	AFS *afero.Afero
}

func (p *PctConfigProcessor) GetConfigMetadata(configFile string) (metadata config_processor.ConfigMetadata, err error) {
	configInfo, err := p.ReadConfig(configFile)
	if err != nil {
		return metadata, err
	}

	err = p.CheckConfig(configFile)
	if err != nil {
		return metadata, err
	}

	metadata = config_processor.ConfigMetadata{
		Author:  configInfo.Template.Author,
		Id:      configInfo.Template.Id,
		Version: configInfo.Template.Version,
	}
	return metadata, nil
}

func (p *PctConfigProcessor) CheckConfig(configFile string) error {
	info, err := p.ReadConfig(configFile)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("The following attributes are missing in %s:\n", configFile)
	orig := msg
	// These parts are essential for build and deployment.

	if info.Template.Id == "" {
		msg = msg + "  * id\n"
	}
	if info.Template.Author == "" {
		msg = msg + "  * author\n"
	}
	if info.Template.Version == "" {
		msg = msg + "  * version\n"
	}
	if msg != orig {
		return fmt.Errorf(msg)
	}

	return nil
}

func (p *PctConfigProcessor) ReadConfig(configFile string) (info pct.PuppetContentTemplateInfo, err error) {
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
