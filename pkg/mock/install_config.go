package mock

import (
	"fmt"

	"github.com/puppetlabs/pct/pkg/config_processor"
)

type InstallConfig struct {
	ExpectedConfigFile string
	Metadata           config_processor.ConfigMetadata
	ErrResponse        error
}

func (ic *InstallConfig) GetConfigMetadata(configFile string) (metadata config_processor.ConfigMetadata, err error) {
	if ic.ErrResponse != nil {
		return metadata, ic.ErrResponse
	}

	if ic.ExpectedConfigFile != configFile {
		return ic.Metadata, fmt.Errorf("configFile (%v) did not match expected value (%v)", configFile, ic.ExpectedConfigFile)
	}

	return ic.Metadata, nil
}

func (ic *InstallConfig) CheckConfig(configFile string) error {
	if ic.ErrResponse != nil {
		return ic.ErrResponse
	}

	if ic.ExpectedConfigFile != configFile {
		return fmt.Errorf("configFile (%v) did not match expected value (%v)", configFile, ic.ExpectedConfigFile)
	}

	return ic.ErrResponse
}
