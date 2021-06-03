package testutils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"

	"github.com/rs/zerolog/log"
)

func SkipAcceptanceTest(t *testing.T) {
	if _, present := os.LookupEnv("TEST_ACCEPTANCE"); !present {
		t.Skip("Skipping, Acceptance test")
	}
}

// Run Command takes a command to execute and the directory in which to execute the command.
// if wd is and empty string it will default to the current working directory
func RunCommand(cmdString string, wd string) (stdout string, stderr string, exitCode int) {
	cmds := strings.Split(cmdString, " ")
	cmd := exec.Command(cmds[0], cmds[1:]...) //#nosec, nolint:gosec //code only used for testing
	if wd != "" {
		cmd.Dir = wd
	}
	out, err := cmd.CombinedOutput()
	exitCode = 0

	if err != nil {
		stderr = err.Error()
		// todo: double check that error statuss work on Windows
		if msg, ok := err.(*exec.ExitError); ok { // there is error code
			exitCode = msg.Sys().(syscall.WaitStatus).ExitStatus()
		}
	}

	stdout = string(out)

	return stdout, stderr, exitCode
}

// Wraps RunCommand for PCT calls, locating the correct binary for the executing OS
func RunPctCommand(cmdString string, wd string) (stdout string, stderr string, exitCode int) {
	// where is pct built for this current OS?
	postfix := ""
	if runtime.GOOS == "windows" {
		postfix = ".exe"
	}

	pctPath := fmt.Sprintf("../dist/pct_%s_%s/pct%s", runtime.GOOS, runtime.GOARCH, postfix)
	absPath, err := filepath.Abs(pctPath)

	if err != nil {
		log.Error().Msgf("Unable to run create path for %s: %s", pctPath, err.Error())
	}

	log.Debug().Msgf("Testing Command: pct %s", cmdString)

	executeString := fmt.Sprintf("%s %s", absPath, cmdString)

	return RunCommand(executeString, wd)
}
