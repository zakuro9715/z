package config

import (
	"errors"
	"io/ioutil"

	"github.com/goccy/go-yaml"
)

/* -- config -- */
type Config struct {
	Shell   string `yaml:"shell"`
	Default string `yaml:"default"`
	Tasks   Tasks  `yaml:"tasks"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config := new(Config)
	if err = yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}
	config.setup()
	return config, nil
}

func (c *Config) setup() {
	for name, t := range c.Tasks {
		t.setup(c, nil, name)
	}
}

/* -- task -- */

type Hooks struct {
	Pre  string `yaml:"pre"`
	Post string `yaml:"post"`
}

type Cmds []string
type Tasks map[string]*Task
type Task struct {
	IsDefault   bool
	Name        string
	fullName    string
	Shell       string `yaml:"shell"`
	Cmds        Cmds   `yaml:"run"`
	RunFile     string `yaml:"run_file"`
	Config      *Config
	Parent      *Task
	Description string `yaml:"desc"`
	Hooks       Hooks  `yaml:"hooks"`
	Tasks       Tasks  `yaml:"tasks"`
}

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

func (t *Task) Verify() error {
	if len(t.Cmds) > 0 && len(t.RunFile) > 0 {
		return errors.New(
			"You can only use either run or script. But both are specified.",
		)
	} else if len(t.Cmds) == 0 && len(t.RunFile) == 0 {
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
