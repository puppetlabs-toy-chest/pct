package exec_runner

import (
	"github.com/rs/zerolog/log"
	"os"
	"os/exec"
	"runtime"
)

type ExecI interface {
	Command(name string, arg ...string) *exec.Cmd
}

type Exec struct{}

func (e *Exec) Command(name string, arg ...string) *exec.Cmd {
	pathToExecutable := ""
	if runtime.GOOS == "windows" {
		pathToExecutable, _ = exec.LookPath("cmd.exe")
	} else {
		pathToExecutable, _ = exec.LookPath(name)
	}
	correctArgs := buildCommandArgs(name, arg)
	log.Debug().Msgf("Path to executable: %v", pathToExecutable)
	log.Debug().Msgf("Command args: %v", correctArgs)
	cmd := &exec.Cmd{
		Path: pathToExecutable,
		Args: correctArgs,
		Env:  os.Environ(),
	}
	return cmd
}

func buildCommandArgs(commandName string, args []string) []string {
	var a []string
	if runtime.GOOS == "windows" {
		a = append(a, "/c")
	}
	a = append(a, commandName)
	a = append(a, args...)
	return a
}
