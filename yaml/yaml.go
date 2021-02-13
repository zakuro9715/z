package yaml

import (
	"github.com/goccy/go-yaml"
)

var (
	decodeOptions = []yaml.DecodeOption{yaml.DisallowDuplicateKey()}
)

func Unmarshal(data []byte, v interface{}) error {
	return yaml.UnmarshalWithOptions(data, v, decodeOptions...)
}
