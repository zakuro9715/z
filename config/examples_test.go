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
