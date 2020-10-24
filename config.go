package main

import (
	"io/ioutil"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Shell string `yaml:"shell"`
	Tasks Tasks  `yaml:"tasks"`
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
	for _, t := range config.Tasks {
		t.setConfig(config)
	}
	return config, nil
}

type Hooks struct {
	Pre  string `yaml:"pre"`
	Post string `yaml:"post"`
}
