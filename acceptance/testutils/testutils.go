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

var app string

func SetAppName(appName string) {
	app = appName
}

func SkipAcceptanceTest(t *testing.T) {
	if _, present := os.LookupEnv("TEST_ACCEPTANCE"); !present {
		t.Skip("Skipping, Acceptance test")
	}
}

func SkipTestInNonCIEnv(t *testing.T) {
	if _, ci := os.LookupEnv("CI"); !ci {
		t.Skip("Skipping test execution in non-CI environment")
	}
}

// Run Command takes a command to execute and the directory in which to execute the command.
// if wd is and empty string it will default to the current working directory
func RunCommand(cmdString string, wd string) (stdout string, stderr string, exitCode int) {
	cmds := strings.Split(cmdString, " ")
	cmd := exec.Command(cmds[0], cmds[1:]...) // #nosec // used only for testing
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

// Wraps RunCommand for App calls, locating the correct binary for the executing OS
func RunAppCommand(cmdString string, wd string) (stdout string, stderr string, exitCode int) {
	// where is app built for this current OS?
	postfix := ""
	if runtime.GOOS == "windows" {
		postfix = ".exe"
	}

	appPath := fmt.Sprintf("../../dist/%s_%s_%s/%s%s", app, runtime.GOOS, runtime.GOARCH, app, postfix)
	absPath, err := filepath.Abs(appPath)

	if err != nil {
		log.Error().Msgf("Unable to run create path for %s: %s", appPath, err.Error())
	}

	log.Debug().Msgf("Testing Command: %s %s", app, cmdString)

	executeString := fmt.Sprintf("%s %s", absPath, cmdString)

	return RunCommand(executeString, wd)
}

// On macOS systems, the `TempDir` func in the `testing` package will
// potentially return the symlink to the dir, rather than the actual
// path (`/private/folders/...` vs `/var/private/folders/...`).
func GetTmpDir(t *testing.T) string {
	dirName := t.TempDir()
	tmpDir, err := filepath.EvalSymlinks(dirName)

	if err != nil {
		panic("Could not create temp dir for test")
	}

	return tmpDir
}
