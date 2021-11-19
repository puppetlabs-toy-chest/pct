package config

type ConfigNamespaceInfo struct {
	Template struct {
		Id      string `mapstructure:"id"`
		Author  string `mapstructure:"author"`
		Version string `mapstructure:"version"`
	} `mapstructure:"template"`
	Defaults map[string]interface{}
}
