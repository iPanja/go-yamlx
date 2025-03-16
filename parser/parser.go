package parser

import (
	"gopkg.in/yaml.v3"
	"os"
)

func ReadYAML(filename string, out interface{}) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, out)
}

func WriteYAML(filename string, in interface{}) error {
	data, err := yaml.Marshal(in)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
