package config

import (
	"io/ioutil"

	"github.com/goccy/go-yaml"
)

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
