package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

type Tasks map[string]*Task
type Task struct {
	Shell       string   `yaml:"shell"`
	Cmds        []string `yaml:"run"`
	Script      string   `yaml:"script"`
	Config      *Config
	Parent      *Task
	Description string `yaml:"desc"`
	Hooks       Hooks  `yaml:"hooks"`
	Tasks       Tasks  `yaml:"tasks"`
}

func (t *Task) setup(c *Config, parent *Task) {
	t.Config = c
	t.Parent = parent
	for _, sub := range t.Tasks {
		sub.setup(c, t)
	}
}

func (t *Task) Verify() error {
	if len(t.Cmds) > 0 && len(t.Script) > 0 {
		return errors.New(
			"You can only use either run or script. But both are specified.",
		)
	} else if len(t.Cmds) == 0 && len(t.Script) == 0 {
		return errors.New("Nothing to run")
	}
	return nil
}

func runWithOsStdio(cmd *exec.Cmd) error {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (t *Task) Run(args []string) error {
	if err := t.Verify(); err != nil {
		return err
	}

	shell := t.GetShell()
	if len(t.Script) > 0 {
		return runWithOsStdio(exec.Command(shell, t.Script))
	}
	for _, command := range t.Cmds {
		cmd := exec.Command(shell, "-c", command+" "+strings.Join(args, " "))
		if err := runWithOsStdio(cmd); err != nil {
			return err
		}
	}
	return nil
}

func (t *Task) GetShell() string {
	if len(t.Shell) > 0 {
		return t.Shell
	}
	if t.Parent == nil {
		if len(t.Config.Shell) == 0 {
			return "sh"
		}
		return t.Config.Shell
	}
	return t.Parent.GetShell()
}
