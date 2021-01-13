package config

import (
	"errors"
	"strings"

	"github.com/goccy/go-yaml"
)

type Cmds []string

func (cmds *Cmds) UnmarshalYAML(data []byte) error {
	var str string
	if err := yaml.Unmarshal(data, &str); err == nil {
		*cmds = []string{str}
		return nil
	}

	ss := []string{}
	err := yaml.Unmarshal(data, &ss)
	*cmds = ss
	return err
}

type Hooks struct {
	Pre  string `yaml:"pre"`
	Post string `yaml:"post"`
}

type ArgsConfig struct {
	Required bool   `yaml:"required"`
	Default  string `yaml:"default"`
}

func (v ArgsConfig) ProcessArgs(args []string) ([]string, error) {
	if len(args) == 0 {
		if len(v.Default) > 0 {
			return strings.Split(v.Default, " "), nil
		} else if v.Required {
			return nil, errors.New("args is required")
		}
	}
	return args, nil
}

type task struct {
	IsDefault   bool
	Name        string
	fullName    string
	Shell       string `yaml:"shell"`
	Cmds        Cmds   `yaml:"run"`
	Config      *Config
	Parent      *Task
	Description string     `yaml:"desc"`
	Hooks       Hooks      `yaml:"hooks"`
	Tasks       Tasks      `yaml:"tasks"`
	ArgsConfig  ArgsConfig `yaml:"args"`
}

type Task struct {
	task
}
type Tasks map[string]*Task

func (task *Task) UnmarshalYAML(data []byte) error {
	err := yaml.Unmarshal(data, &task.task)
	return err
}

func (t *Task) setup(c *Config, parent *Task, name string) {
	names := strings.SplitN(name, ".", 2)
	if len(names) > 1 {
		sub := *t // copy
		*t = Task{
			task{Tasks: map[string]*Task{names[1]: &sub}},
		}
	}
	t.Name = names[0]
	t.Config = c
	t.Parent = parent
	if parent != nil {
		t.fullName = parent.fullName + "." + t.Name
	} else {
		t.fullName = t.Name
	}
	if t.fullName == c.Default {
		t.IsDefault = true
	}
	for name, sub := range t.Tasks {
		sub.setup(c, t, name)
	}
}

func (t *Task) Verify() error {
	if len(t.Cmds) == 0 {
		return errors.New("Nothing to run")
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
