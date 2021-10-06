package exec_runner

import (
	"github.com/puppetlabs/pdkgo/internal/pkg/utils"
	"os/exec"
)

type ExecI interface {
	Command(name string, arg ...string) *exec.Cmd
}

type Exec struct{}

func (e *Exec) Command(name string, arg ...string) *exec.Cmd {
	return utils.Command(name, arg...)
}
