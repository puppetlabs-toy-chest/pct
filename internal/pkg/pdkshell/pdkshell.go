package pdkshell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/rs/zerolog/log"
)

type PDKInfo struct {
	RubyVersion      string
	InstallDirectory string
	RubyExecutable   string
	PDKExecutable    string
	CertDirectory    string
	CertPemFile      string
}

// Execute runs a given pdk command
// It first detects where the PDK Ruby installation is on the local system
// Then it builds the correct command line to execute the PDK ruby with provided arguments
func Execute(args []string) (int, error) {
	i := getPDKInfo()
	executable := buildExecutable(i.RubyExecutable)
	args = buildCommandArgs(args, i.RubyExecutable, i.PDKExecutable)
	env := os.Environ()
	env = append(env, fmt.Sprintf("SSL_CERT_DIR=%s", i.CertDirectory), fmt.Sprintf("SSL_CERT_FILE=%s", i.CertPemFile))
	cmd := &exec.Cmd{
		Path:   executable,
		Args:   args,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Env:    env,
	}

	log.Trace().Msgf("args: %s", args)
	if err := cmd.Run(); err != nil {
		log.Fatal().Msgf("pdk failed with '%s'\n", err)

		return err.(*exec.ExitError).ExitCode(), err
	}

	return 0, nil
}

// getPDKInfo detects where the PDK Ruby installation is on the local file system
// It handles detecting the installation on Windows and other platforms
func getPDKInfo() *PDKInfo {
	rubyVersion := "2.4.10"
	installDir, err := getPDKInstallDirectory(true)
	if err != nil {
		log.Fatal().Msgf("error: %v", err)
	}

	i := &PDKInfo{
		RubyVersion:      rubyVersion,
		InstallDirectory: installDir,
		RubyExecutable:   filepath.Join(installDir, "private", "ruby", rubyVersion, "bin", "ruby"),
		PDKExecutable:    filepath.Join(installDir, "private", "ruby", rubyVersion, "bin", "pdk"),
		CertDirectory:    filepath.Join(installDir, "ssl", "certs"),
		CertPemFile:      filepath.Join(installDir, "ssl", "cert.pem"),
	}
	return i
}

// buildExecutable returns the executable to use to execute the PDK command
// On windows it returns `cmd.exe` as golang's exec needs a shell on Windows to run
// On *nix like platforms, it returns the Ruby exectuable provided as there is no need to specify a shell
func buildExecutable(rubyexe string) (executable string) {
	executable = rubyexe
	if runtime.GOOS == "windows" {
		exe, _ := exec.LookPath("cmd.exe")
		executable = exe
	}
	return executable
}

// buildCommandArgs takes an array of commandli arguments and the path to the
// ruby installation and returns a fully formatted commandline to execute
// On Windows, it prepends with `/c` so that cmd.exe executes the command properly
func buildCommandArgs(args []string, rubyexe, pdkexe string) []string {
	var a []string
	if runtime.GOOS == "windows" {
		a = append(a, "/c")
	}
	a = append(a, rubyexe, "-S", "--", pdkexe)
	a = append(a, args...)
	return a
}
