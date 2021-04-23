// +build !windows

package pdkshell

import "github.com/rs/zerolog/log"

// getPDKInstallDirectory returns the directory PDK is located in
func getPDKInstallDirectory(shortName bool) (string, error) {
	pdkInstallDir := "/opt/puppetlabs/pdk"
	log.Trace().Msgf("Install dir : %s", pdkInstallDir)
	return pdkInstallDir, nil
}
