package config

import (
	"errors"
	"strings"

	"github.com/valyala/fasttemplate"
	"github.com/zakuro9715/z/yaml"
)

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
	dummy       bool
	IsDefault   bool
	Name        string
	FullName    string
	Shell       string          `yaml:"shell"`
	Cmds        yaml.StringList `yaml:"run"`
	Envs        Envs            `yaml:"env"`
	Paths       yaml.StringList `yaml:"path"`
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
	var str string
	if err := yaml.Unmarshal(data, &str); err == nil {
		task.Cmds = []string{str}
		task.Description = str
		return nil
	}
	var strs []string
	if err := yaml.Unmarshal(data, &strs); err == nil {
		task.Cmds = strs
		task.Description = strings.Join(strs, "\n")
		return nil
	}
	if err := yaml.Unmarshal(data, &task.task); err != nil {
		return err
	}
	if len(task.Description) == 0 {
		task.Description = strings.Join(task.Cmds, "\n")
	}
	return nil
}

func (t *Task) setup(c *Config, parent *Task, name string) {
	if t == nil {
		return
	}
	names := strings.SplitN(name, ".", 2)
	if len(names) > 1 {
		sub := *t // copy
		*t = Task{
			task{
				dummy: true,
				Tasks: map[string]*Task{names[1]: &sub},
			},
		}
	}
	t.Name = names[0]
	t.Config = c
	t.Parent = parent
	if parent != nil {
		t.FullName = parent.FullName + "." + t.Name
	} else {
		t.FullName = t.Name
	}
	if t.FullName == c.Default {
		t.IsDefault = true
	}
	t.expandVars()

	for name, sub := range t.Tasks {
		sub.setup(c, t, name)
	}
}

func template(s string, m map[string]interface{}) string {
	return fasttemplate.New(s, "{{", "}}").ExecuteString(m)
}

func (t *Task) expandVars() {
	vars := map[string]interface{}{}
	for k, v := range t.Config.Vars {
		vars[k] = v
	}
	t.Description = template(t.Description, vars)
	for i, cmd := range t.Cmds {
		t.Cmds[i] = template(cmd, vars)
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
