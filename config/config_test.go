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

func TestSetup(t *testing.T) {
	c := &Config{
		Default: "hello.world",
		Tasks: map[string]*Task{
			"hello": {
				Tasks: map[string]*Task{
					"world": {},
				},
			},
		},
	}

	assert.Empty(t, c.Tasks["hello"].Name)
	assert.Empty(t, c.Tasks["hello"].fullName)
	assert.Nil(t, c.Tasks["hello"].Config)
	assert.False(t, c.Tasks["hello"].IsDefault)
	assert.Empty(t, c.Tasks["hello"].Tasks["world"].Name)
	assert.Empty(t, c.Tasks["hello"].Tasks["world"].fullName)
	assert.Nil(t, c.Tasks["hello"].Tasks["world"].Config)

	c.setup()
	assert.Equal(t, "hello", c.Tasks["hello"].Name)
	assert.Equal(t, "hello", c.Tasks["hello"].fullName)
	assert.Equal(t, c, c.Tasks["hello"].Config)
	assert.True(t, c.Tasks["hello"].Tasks["world"].IsDefault)
	assert.Equal(t, "world", c.Tasks["hello"].Tasks["world"].Name)
	assert.Equal(t, "hello.world", c.Tasks["hello"].Tasks["world"].fullName)
	assert.Equal(t, c, c.Tasks["hello"].Tasks["world"].Config)
}

func TestVerify(t *testing.T) {
	assert.NoError(t, (&Task{Cmds: []string{"v"}}).Verify())
	assert.Error(t, (&Task{}).Verify())
}
