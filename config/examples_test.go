package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExampleIsValid(t *testing.T) {
	assert.NoError(t,
		filepath.Walk("../examples", func(path string, info os.FileInfo, err error) error {
			assert.NoError(t, err)
			if info.IsDir() || !strings.HasSuffix(path, "yaml") {
				return nil
			}
			println(path)
			_, err = LoadConfig(path)
			return err
		}),
	)
}

func TestLoadHelloExample(t *testing.T) {
	expected := &Config{
		Shell:   "bash",
		Default: "hello",
		Tasks: map[string]*Task{
			"hello": {
				Cmds:        []string{"echo hello", "echo bye"},
				Description: "Say hello",
				Hooks: Hooks{
					Pre:  "echo saying hello",
					Post: "echo said hello",
				},
				ArgsConfig: ArgsConfig{
					Required: true,
					Default:  "you",
				},
				Tasks: map[string]*Task{
					"world": {
						Cmds: []string{"z hello -- world"},
					},
					"script": {
						Cmds: []string{"examples/hello.sh"},
					},
					"python": {
						Shell: "python",
						Cmds:  []string{"print('hello python')"},
					},
				},
			},
		},
	}
	expected.setup()

	actual, err := LoadConfig("../examples/hello.yaml")
	assert.NoError(t, err)

	assert.Equal(t, expected, actual)
}
