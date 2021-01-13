package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerify(t *testing.T) {
	assert.NoError(t, (&Task{Cmds: []string{"v"}}).Verify())
	assert.Error(t, (&Task{}).Verify())
}
