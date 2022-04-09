package config_processor

type ConfigProcessorI interface {
	GetConfigMetadata(configFile string) (metadata ConfigMetadata, err error)
	CheckConfig(configFile string) error
}

type ConfigMetadata struct {
	Id      string
	Author  string
	Version string
}
