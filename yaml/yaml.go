package yaml

import (
	"github.com/goccy/go-yaml"
)

var (
	decodeOptions = []yaml.DecodeOption{yaml.DisallowDuplicateKey()}
)

type StringList []string

func (v *StringList) UnmarshalYAML(data []byte) error {
	var str string
	if err := yaml.Unmarshal(data, &str); err == nil {
		*v = []string{str}
		return nil
	}

	ss := []string{}
	err := yaml.Unmarshal(data, &ss)
	*v = ss
	return err
}

func Unmarshal(data []byte, v interface{}) error {
	return yaml.UnmarshalWithOptions(data, v, decodeOptions...)
}
