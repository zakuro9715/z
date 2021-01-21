package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerify(t *testing.T) {
	assert.NoError(t, (&Task{task{Cmds: []string{"v"}}}).Verify())
	assert.NoError(t, (&Task{task{AliasTo: "v"}}).Verify())
	assert.Error(t, (&Task{}).Verify())
	assert.Error(t, (&Task{task{Cmds: []string{"v"}, AliasTo: "v"}}).Verify())
}
