package config

import (
	"io/ioutil"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/zakuro9715/z/log"
)

type Envs map[string]string

type Config struct {
	Shell    string `yaml:"shell"`
	Default  string `yaml:"default"`
	Tasks    Tasks  `yaml:"tasks"`
	Envs     Envs   `yaml:"env"`
	allTasks Tasks
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

func (c *Config) FindTask(splitedFullName ...string) *Task {
	return c.allTasks[strings.Join(splitedFullName, ".")]
}
