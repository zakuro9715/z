package yaml

import (
	"github.com/goccy/go-yaml"
	goyaml "github.com/goccy/go-yaml"
)

var (
	decodeOptions = []goyaml.DecodeOption{yaml.DisallowDuplicateKey()}
)

type StringList []string

func (v *StringList) UnmarshalYAML(data []byte) error {
	var str string
	if err := Unmarshal(data, &str); err == nil {
		*v = []string{str}
		return nil
	}

	ss := []string{}
	err := Unmarshal(data, &ss)
	*v = ss
	return err
}

func Unmarshal(data []byte, v interface{}) error {
	return goyaml.UnmarshalWithOptions(data, v, decodeOptions...)
}
