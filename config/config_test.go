package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func (tasks Tasks) FullNames() []string {
	names := make([]string, 0, len(tasks))
	for name := range tasks {
		names = append(names, name)
	}
	return names
}

func TestLoadConfigError(t *testing.T) {
	config, err := LoadConfig("nonexists")
	assert.Nil(t, config)
	assert.Error(t, err)
}

func TestSetup(t *testing.T) {
	c := &Config{
		config{
			Default: "hello.world",
			Tasks: map[string]*Task{
				"hello": {
					task{
						Tasks: map[string]*Task{
							"world": {},
						},
					},
				},
				"echo.hey": {},
			},
		},
		nil,
	}

	assert.Empty(t, c.Tasks["hello"].Name)
	assert.Empty(t, c.Tasks["hello"].FullName)
	assert.Nil(t, c.Tasks["hello"].Config)
	assert.False(t, c.Tasks["hello"].IsDefault)
	assert.Empty(t, c.Tasks["hello"].Tasks["world"].Name)
	assert.Empty(t, c.Tasks["hello"].Tasks["world"].FullName)
	assert.Nil(t, c.Tasks["hello"].Tasks["world"].Config)
	assert.Nil(t, c.allTasks)
	assert.Nil(t, c.Tasks["echo"])

	c.setup()

	assert.Equal(t, "hello", c.Tasks["hello"].Name)
	assert.Equal(t, "hello", c.Tasks["hello"].FullName)
	assert.Equal(t, c, c.Tasks["hello"].Config)
	assert.True(t, c.Tasks["hello"].Tasks["world"].IsDefault)
	assert.Equal(t, "world", c.Tasks["hello"].Tasks["world"].Name)
	assert.Equal(t, "hello.world", c.Tasks["hello"].Tasks["world"].FullName)
	assert.Equal(t, c, c.Tasks["hello"].Tasks["world"].Config)
	fullNames := []string{"echo.hey", "hello", "hello.world"}
	assert.ElementsMatch(t, fullNames, c.allTasks.FullNames())
}

func TestFindTask(t *testing.T) {
	c := &Config{
		config{
			Default: "hello.world",
			Tasks: map[string]*Task{
				"hello": {
					task{
						Tasks: map[string]*Task{
							"world": {},
						},
					},
				},
			},
		},
		nil,
	}
	c.setup()

	findTask := func(args ...string) *Task {
		t, _ := c.FindTask(args...)
		return t
	}
	assert.Equal(t, c.Tasks["hello"], findTask("hello"))
	assert.Equal(t, c.Tasks["hello"].Tasks["world"], findTask("hello.world"))
	assert.Equal(t, c.Tasks["hello"].Tasks["world"], findTask("hello", "world"))
	assert.Nil(t, findTask("null"))
}
