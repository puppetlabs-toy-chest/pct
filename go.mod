module github.com/puppetlabs/pdkgo

go 1.16

replace github.com/puppetlabs/pdkgo/docs/md => ./docs/md

require (
	github.com/alecthomas/chroma v0.9.4 // indirect
	github.com/charmbracelet/glamour v0.3.0
	github.com/denisbrodbeck/machineid v1.0.1
	github.com/gernest/front v0.0.0-20210301115436-8a0b0a782d0a
	github.com/go-logr/logr v1.2.2 // indirect
	github.com/hashicorp/go-version v1.3.0
	github.com/json-iterator/go v1.1.12
	github.com/microcosm-cc/bluemonday v1.0.16 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/muesli/reflow v0.3.0 // indirect
	github.com/muesli/termenv v0.9.0 // indirect
	github.com/olekukonko/tablewriter v0.0.5
	github.com/puppetlabs/pdkgo/docs/md v0.0.0-00010101000000-000000000000
	github.com/rs/zerolog v1.26.0
	github.com/spf13/afero v1.6.0
	github.com/spf13/cobra v1.3.0
	github.com/spf13/viper v1.10.0
	github.com/stretchr/testify v1.7.0
	github.com/yuin/goldmark v1.4.4 // indirect
	go.opentelemetry.io/otel v1.3.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.3.0
	go.opentelemetry.io/otel/sdk v1.3.0
	go.opentelemetry.io/otel/trace v1.3.0
	google.golang.org/grpc v1.43.0
	gopkg.in/yaml.v2 v2.4.0
)
