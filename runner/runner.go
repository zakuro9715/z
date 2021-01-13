package runner

import (
	"os"
	"os/exec"
	"strings"

	"github.com/zakuro9715/z/config"
)

type TaskRunner struct {
	task *config.Task
}

func New(t *config.Task) *TaskRunner {
	return &TaskRunner{t}
}

func runWithOsStdio(cmd *exec.Cmd) error {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func isScriptFile(name string) bool {
	// err or dir or executable
	if f, err := os.Stat(name); err != nil || f.IsDir() || (f.Mode()&0111 == 0111) {
		return false
	}
	return true
}

func (r *TaskRunner) Run(args []string) error {
	if err := r.task.Verify(); err != nil {
		return err
	}

	shell := r.task.GetShell()
	for _, command := range r.task.Cmds {
		if isScriptFile(command) {
			return runWithOsStdio(exec.Command(shell, append([]string{command}, args...)...))
		} else {
			cmd := exec.Command(shell, "-c", command+" "+strings.Join(args, " "))
			if err := runWithOsStdio(cmd); err != nil {
				return err
			}
		}
	}
	return nil
}
