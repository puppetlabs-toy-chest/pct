# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- [(GH-222)](https://github.com/uppetlabs/pdkgo/issues/222) Telemetry to the binary, which will report the operating system type and architecture when a command is run; the implementation allows for two binaries: one with telemetry configured and enabled, and one _without_ the telemetry included at all. <!-- For more information, see our [telemetry blog post](link to blog). -->
- [(GH-223)](https://github.com/uppetlabs/pdkgo/issues/223) Added hashed machine uuid generation and included in the telemetry; this will report a universally unique machine ID for each node running PCT and reporting telemetry.
- [(GH-224)](https://github.com/uppetlabs/pdkgo/issues/224) Added span tracing for each PCT command, ensuring telemetry can tell which commands are being called how often.
- [(GH-225)](https://github.com/uppetlabs/pdkgo/issues/225) Added error records to span tracing, ensuring that error information is included in telemetry for analysis and prioritization

## [0.4.0]

### Changed

- The Puppet Content templates shipped in 0.4.0 and the handling of templates in 0.4.0 is _not_ backward compatible with templates which do not have `id`, `author`, AND `version` defined in their metadata

### Added

- [(GH-183)](https://github.com/uppetlabs/pdkgo/issues/183) `pct new` handles namespaced templates
- [(GH-184)](https://github.com/uppetlabs/pdkgo/issues/184) `pct install` works against remote `tar.gz` files
- [(GH-185)](https://github.com/puppetlabs/pdkgo/issues/185) `pct build` validates pct-config.yml
- [(GH-167)](https://github.com/puppetlabs/pdkgo/issues/167) Implement `pct install` CLI command
- [(TEMPLATES-17)](https://github.com/puppetlabs/baker-round/issues/17) Ensure `puppet-content-template` includes the author key in the scaffolded config file
- [(TEMPLATES-18)](https://github.com/puppetlabs/baker-round/issues/18) Ensure all default templates have their author set to `puppetlabs`

## [0.3.0]

- [(GH-144)](https://github.com/puppetlabs/pdkgo/issues/144) Implement `pct build` CLI command

### Removed

- [(GH-172)](https://github.com/puppetlabs/pdkgo/issues/172) Removal of PDKShell commands

## [0.2.0]

### Added

- [(GH-83)](https://github.com/puppetlabs/pdkgo/issues/83) Allow for workspace configuration overrides
- [(GH-107](https://github.com/puppetlabs/pdkgo/issues/107) Initialize zerolog via cobra.OnInitialize method

### Fixed

- [(GH-15)](https://github.com/puppetlabs/pdkgo/issues/15) Unset necessary env vars in pdkshell
- [(GH-125)](https://github.com/puppetlabs/pdkgo/issues/125) Fail on errors, quote arguments
- [(GH-125)](https://github.com/puppetlabs/pdkgo/issues/125) Fix `$ver` bug in download script

## [0.1.0]

### Added

- [(GH-67)](https://github.com/puppetlabs/pdkgo/issues/67) Add installation scripts for PCT

### Fixed

- [(GH-64)](https://github.com/puppetlabs/pdkgo/issues/64) Strip pct from command name
- [(GH-65)](https://github.com/puppetlabs/pdkgo/issues/65) Allow deployment of empty files
- [(GH-14)](https://github.com/puppetlabs/pdkgo/issues/14) Return the exit code from the PDK when executed by the wrapper

## [0.1.0-pre]

### Added

- [(GH-2)](https://github.com/puppetlabs/pdkgo/issues/2) Created Puppet Content Templates package and modified pdk new to use PCT
- [(GH-7)](https://github.com/puppetlabs/pdkgo/issues/7) Added wrapper to all existing PDK commands

### Fixed

- [(GH-29)](https://github.com/puppetlabs/pdkgo/issues/29) Error if template not found

[Unreleased]: https://github.com/puppetlabs/pdkgo/compare/0.4.0..main
[0.4.0]: https://github.com/puppetlabs/pdkgo/releases/tag/0.3.0
[0.3.0]: https://github.com/puppetlabs/pdkgo/releases/tag/0.3.0
[0.2.0]: https://github.com/puppetlabs/pdkgo/releases/tag/0.2.0
[0.1.0]: https://github.com/puppetlabs/pdkgo/releases/tag/0.1.0
[0.1.0-pre]: https://github.com/puppetlabs/pdkgo/releases/tag/0.1.0-pre
