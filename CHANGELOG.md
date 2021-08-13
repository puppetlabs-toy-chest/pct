# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

[Unreleased]: https://github.com/puppetlabs/pdkgo/compare/0.2.0..main
[0.2.0]: https://github.com/puppetlabs/pdkgo/releases/tag/0.2.0
[0.1.0]: https://github.com/puppetlabs/pdkgo/releases/tag/0.1.0
[0.1.0-pre]: https://github.com/puppetlabs/pdkgo/releases/tag/0.1.0-pre
