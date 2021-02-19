package yaml

import (
	"strings"

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

type StringKeyValueList map[string]string

func (v *StringKeyValueList) UnmarshalYAML(data []byte) error {
	dict := map[string]string{}
	list := StringList{}
	if err := Unmarshal(data, &list); err == nil {
		for _, s := range list {
			parts := strings.SplitN(s, "=", 2)
			key := strings.TrimSpace(parts[0])
			switch len(parts) {
			case 1:
				dict[key] = ""
			case 2:
				dict[key] = strings.TrimSpace(parts[1])
			default:
				panic("unreachable code")
			}
		}
		*v = dict
		return nil
	}
	if err := Unmarshal(data, &dict); err != nil {
		return err
	}
	*v = dict
	return nil
}

func Unmarshal(data []byte, v interface{}) error {
	return goyaml.UnmarshalWithOptions(data, v, decodeOptions...)
}
