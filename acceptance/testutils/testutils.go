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

// Run Command takes a command to execute and the directory in which to execute the command.
// if wd is and empty string it will default to the current working directory
func RunCommand(cmdString string, wd string) (stdout string, stderr string, exitCode int) {
	cmds := strings.Split(cmdString, " ")

	cmds = toolArgsAsSingleArg(cmds) // Remove when GH-52 is resolved

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

// Remove when GH-52 is resolved
// This function assumes that '--toolArgs' is the final argument passed to the app
func toolArgsAsSingleArg(cmds []string) []string {
	toolArgsFragmentIndex := 0
	var reassembled strings.Builder

	for i, arg := range cmds {
		if strings.HasPrefix(arg, "--toolArgs=") {
			toolArgsFragmentIndex = i
		}
		if toolArgsFragmentIndex > 0 {
			reassembled.WriteString(fmt.Sprintf("%s ", arg))
		}
	}

	if toolArgsFragmentIndex > 0 {
		var cmdAndArgs []string
		cmdAndArgs = append(cmdAndArgs, cmds[0])
		cmdAndArgs = append(cmdAndArgs, cmds[1:toolArgsFragmentIndex]...)
		cmdAndArgs = append(cmdAndArgs, reassembled.String())
		return cmdAndArgs
	}

	return cmds
}
