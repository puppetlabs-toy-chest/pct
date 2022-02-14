package config_processor

type ConfigProcessorI interface {
	ProcessConfig(sourceDir, targetDir string, force bool) (string, error)
	CheckConfig(configFile string) error
}
