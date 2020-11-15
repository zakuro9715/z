package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/zakuro9715/z/config"
)

type TaskRunner struct {
	task *config.Task
}

func NewTaskRunner(t *config.Task) *TaskRunner {
	return &TaskRunner{t}
}

func runWithOsStdio(cmd *exec.Cmd) error {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (r *TaskRunner) Run(args []string) error {
	if err := r.task.Verify(); err != nil {
		return err
	}

	shell := r.task.GetShell()
	if len(r.task.Script) > 0 {
		return runWithOsStdio(exec.Command(shell, r.task.Script))
	}
	for _, command := range r.task.Cmds {
		cmd := exec.Command(shell, "-c", command+" "+strings.Join(args, " "))
		if err := runWithOsStdio(cmd); err != nil {
			return err
		}
	}
	return nil
}
