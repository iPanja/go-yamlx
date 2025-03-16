package utils

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

func PrintYAMLStructure(data interface{}) {
	fmt.Printf("%+v\n", data)
}

func NewScalarNode() *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: "",
	}
}

func ParseYaml(yamlText string) *yaml.Node {
	var Node yaml.Node
	_ = yaml.Unmarshal([]byte(yamlText), &Node)

	return &Node
}
