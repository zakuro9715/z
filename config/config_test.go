package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigError(t *testing.T) {
	config, err := LoadConfig("nonexists")
	assert.Nil(t, config)
	assert.Error(t, err)
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
