package config

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/zakuro9715/z/log"
	"github.com/zakuro9715/z/yaml"
)

type Envs map[string]string

func (v *Envs) UnmarshalYAML(data []byte) error {
	dict := map[string]string{}
	list := yaml.StringList{}
	if err := yaml.Unmarshal(data, &list); err == nil {
		for _, s := range list {
			parts := strings.SplitN(s, "=", 2)
			key := strings.TrimSpace(parts[0])
			switch len(parts) {
			case 1:
				dict[key] = ""
			case 2:
				dict[key] = strings.TrimSpace(parts[1])
			default:
				panic("unreachable code")
			}
		}
		*v = dict
		return nil
	}
	if err := yaml.Unmarshal(data, &dict); err != nil {
		return err
	}
	*v = dict
	return nil
}

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
