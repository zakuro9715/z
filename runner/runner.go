package runner

import (
	"os"
	"os/exec"
	"strings"

	"github.com/zakuro9715/z/config"
	"github.com/zakuro9715/z/log"
)

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

func prependPath(path string) error {
	return os.Setenv("PATH", path+":"+os.Getenv("PATH"))
}

func preparePath(task *config.Task) error {
	if task.Parent == nil {
		for _, path := range task.Config.Paths {
			if err := prependPath(path); err != nil {
				return err
			}
		}
	} else {
		if err := preparePath(task.Parent); err != nil {
			return err
		}
		for _, path := range task.Paths {
			if err := prependPath(path); err != nil {
				return err
			}
		}
	}
	return nil
}

func prepareEnv(task *config.Task) error {
	if err := setTaskDefaultEnvs(task); err != nil {
		return err
	}
	return preparePath(task)
}

func logEnv(task *config.Task) {
	if n, _ := log.Info("envs:"); n == 0 {
		return
	}

	envs := map[string]bool{}
	for ; task != nil; task = task.Parent {
		for k := range task.Envs {
			envs[k] = true
		}
	}
	for k := range envs {
		log.Info("  %v=%v", k, os.Getenv(k))
	}
}

type Runner struct {
	config *Config
}

type Config struct {
	Silent bool
}

func New(config *Config) *Runner {
	return &Runner{config}
}

func Run(config *Config, task *config.Task, args []string) error {
	return New(config).Run(task, args)
}

func (r Runner) Run(task *config.Task, args []string) error {
	if err := task.Verify(); err != nil {
		return err
	}

	if err := prepareEnv(task); err != nil {
		return err
	}

	args, err := task.ArgsConfig.ProcessArgs(args)
	if err != nil {
		return err
	}
	shell := task.GetShell()
	log.Info("shell:")
	log.Info("  " + shell)

	log.Info("task:")
	log.Info("  " + task.FullName)
	log.Info("args:")
	log.Info("  %v", args)
	logEnv(task)

	argsStr := strings.Join(args, " ")
	log.Info("run:")
	for _, command := range task.Cmds {
		if isScriptFile(command) {
			log.Info("  script: %v %v", command, argsStr)
			return r.runCmd(exec.Command(shell, append([]string{command}, args...)...))
		} else {
			log.Info("  command: %v %v", command, argsStr)
			cmd := exec.Command(shell, "-c", command+" "+strings.Join(args, " "))
			if err := r.runCmd(cmd); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *Runner) runCmd(cmd *exec.Cmd) error {
	cmd.Stdin = os.Stdin
	if !r.config.Silent {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
