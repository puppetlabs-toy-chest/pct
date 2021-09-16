/*
Package pct implements the Puppet Content template specification

Puppet Content Templates (PCT) codify a structure to produce content for any Puppet
Product. PCT can create any type of a Puppet Product project: Puppet control
repo, Puppet Module, Bolt project, etc. PCT can also create one or more independent
files, such as CI files or gitignores. This can be as simple as a name for a
Puppet Class or a set of CI files to add to a Puppet Module.
*/
package prm

import (
	"github.com/spf13/afero"
)

// PDKInfo contains the current version information of the compiled binary for
// use in template data
type PRMInfo struct {
	Version   string
	Commit    string
	BuildDate string
}

type PrmI interface {
	GetApplicationInfo() string
}

type Prm struct {
	AFS  *afero.Afero
	IOFS *afero.IOFS
}

func (p *Prm) GetApplicationInfo() string {
	return "0.1.0"
}
