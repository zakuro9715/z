package main

import (
	"os"
	"os/exec"
	"strings"
)

type Tasks map[string]*Task
type Task struct {
	Cmds        []string `yaml:"run"`
	Config      *Config
	Description string `yaml:"desc"`
	Hooks       Hooks  `yaml:"hooks"`
	Tasks       Tasks  `yaml:"tasks"`
}

func (t *Task) setConfig(c *Config) {
	t.Config = c
	for _, sub := range t.Tasks {
		sub.setConfig(c)
	}
}

func (t *Task) Run(args []string) error {
	shell := t.Config.GetShell()
	for _, command := range t.Cmds {
		cmd := exec.Command(shell, "-c", command+" "+strings.Join(args, " "))
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) GetShell() string {
	if len(c.Shell) == 0 {
		return "sh"
	}
	return c.Shell
}
