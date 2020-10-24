package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExampleIsValid(t *testing.T) {
	assert.NoError(t,
		filepath.Walk("examples", func(path string, info os.FileInfo, err error) error {
			assert.NoError(t, err)
			if info.IsDir() {
				return nil
			}
			_, err = LoadConfig(path)
			return err
		}),
	)
}

func TestLoadHelloExample(t *testing.T) {
	expected := &Config{
		Shell: "bash",
		Tasks: map[string]*Task{
			"hello": {
				Cmds: []string{"echo hello"},
				Hooks: Hooks{
					Pre:  "echo saying hello",
					Post: "echo said hello",
				},
				Tasks: map[string]*Task{
					"world": {
						Cmds:        []string{"z hello -- world"},
						Description: "Hello World",
					},
					"cjk": {
						Cmds: []string{
							"z hello japan",
							"z hello korea",
							"z hello china",
						},
					},
				},
			},
		},
	}
	for _, task := range expected.Tasks {
		task.setConfig(expected)
	}

	actual, err := LoadConfig("examples/hello.yaml")
	assert.NoError(t, err)

	assert.Equal(t, expected, actual)
}
