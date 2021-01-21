package runner

import (
	"os"
	"os/exec"
	"strings"

	"github.com/zakuro9715/z/config"
	"github.com/zakuro9715/z/log"
)

type TaskRunner struct {
}

func New() *TaskRunner {
	return &TaskRunner{}
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

func setDefaultEnvs(envs config.Envs) error {
	for k, v := range envs {
		if got := os.Getenv(k); len(got) > 0 {
			continue
		}
		if err := os.Setenv(k, v); err != nil {
			return err
		}
	}
	return nil
}

func setTaskDefaultEnvs(task *config.Task) error {
	if err := setDefaultEnvs(task.Envs); err != nil {
		return err
	}
	if task.Parent == nil {
		return setDefaultEnvs(task.Config.Envs)
	}
	return setTaskDefaultEnvs(task.Parent)
}

func (r *TaskRunner) Run(task *config.Task, args []string) error {
	if err := task.Verify(); err != nil {
		return err
	}

	if err := setTaskDefaultEnvs(task); err != nil {
		return err
	}

	args, err := task.ArgsConfig.ProcessArgs(args)
	if err != nil {
		return err
	}
	shell := task.GetShell()
	log.Infof("shell: %v\n", shell)

	for _, command := range task.Cmds {
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
