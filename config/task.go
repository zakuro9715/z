package config

import (
	"errors"

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

type Task struct {
	IsDefault   bool
	Name        string
	fullName    string
	Shell       string `yaml:"shell"`
	Cmds        Cmds   `yaml:"run"`
	Config      *Config
	Parent      *Task
	Description string `yaml:"desc"`
	Hooks       Hooks  `yaml:"hooks"`
	Tasks       Tasks  `yaml:"tasks"`
}
type Tasks map[string]*Task

func (t *Task) setup(c *Config, parent *Task, name string) {
	t.Name = name
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
