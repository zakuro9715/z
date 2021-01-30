package config

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/zakuro9715/z/log"
)

type oneOrMoreStr []string

func (v *oneOrMoreStr) UnmarshalYAML(data []byte) error {
	var str string
	if err := yaml.Unmarshal(data, &str); err == nil {
		*v = []string{str}
		return nil
	}

	ss := []string{}
	err := yaml.Unmarshal(data, &ss)
	*v = ss
	return err
}

type Envs map[string]string

type config struct {
	Shell   string `yaml:"shell"`
	Default string `yaml:"default"`
	Tasks   Tasks  `yaml:"tasks"`
	Envs    Envs   `yaml:"env"`
}

type Config struct {
	config
	allTasks Tasks
}

func (c *Config) UnmarshalYAML(data []byte) error {
	err := yaml.Unmarshal(data, &c.config)
	return err
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

	if n, _ := log.Info("config:"); n > 0 {
		log.Info("  path:")
		log.Info("    " + filename)
		log.Info("  default:")
		log.Info("    " + config.Default)
		log.Info("  tasks:")
		for _, task := range config.allTasks {
			log.Info("    " + task.FullName)
		}
	}
	return config, nil
}

func (c *Config) collectTasksRecursive(task *Task) {
	if !task.dummy {
		c.allTasks[task.FullName] = task
	}
	for _, t := range task.Tasks {
		c.collectTasksRecursive(t)
	}
}

func (c *Config) setup() {
	c.allTasks = Tasks{}
	for name, t := range c.Tasks {
		t.setup(c, nil, name)
		c.collectTasksRecursive(t)
	}
}

func (c *Config) FindTask(splitedFullName ...string) (*Task, error) {
	fullName := strings.Join(splitedFullName, ".")
	task, ok := c.allTasks[fullName]
	if ok {
		return task, nil
	}
	return nil, fmt.Errorf("Unknown task: %v", fullName)
}
