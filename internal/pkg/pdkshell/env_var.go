package pdkshell

import (
	"strings"

	"github.com/rs/zerolog/log"
)

var osUtil osHelpers

func init() {
	osUtil = osHelpersImpl{}
}

// There are a number of ENV VARs that need to be unset on POSIX systems before execution:
// https://github.com/puppetlabs/pdk-vanagon/blob/0aa54c9129b137c2deabb0a417d59215df36fd91/resources/files/posix/pdk_env_wrapper
func getEnvVarsToUnset() []string {
	return []string{
		"GEM_HOME",
		"GEM_PATH",
		"DLN_LIBRARY_PATH",
		"RUBYLIB",
		"RUBYLIB_PREFIX",
		"RUBYOPT",
		"RUBYPATH",
		"RUBYSHELL",
		"LD_LIBRARY_PATH",
		"LD_PRELOAD",
	}
}

func validEnvVar(envVar string) (valid bool) {
	for _, illegalEnvVar := range getEnvVarsToUnset() {
		if strings.HasPrefix(envVar, illegalEnvVar) {
			log.Trace().Msgf("Dropping ENV VAR: %s", envVar)
			return false
		}
	}
	log.Trace().Msgf("Permitting ENV VAR: %s", envVar)
	return true
}

func getEnvironVars() []string {
	env := []string{}
	for _, envVar := range osUtil.Environ() {
		if validEnvVar(envVar) {
			env = append(env, envVar)
		}
	}
	return env
}
