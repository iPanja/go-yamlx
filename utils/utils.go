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

// MergeNodes will merge src -> dst while following proper YAML merge behavior
//
// This function will modify dst. Use a duplicate, dummy, dst node if needed.
func MergeNodes(dst, src *yaml.Node) {
	if dst.Kind != yaml.MappingNode || src.Kind != yaml.MappingNode {
		return // Only merge mapping nodes
	}

	// Create a map of keys in the destination node
	keyMap := make(map[string]int)
	for i := 0; i < len(dst.Content); i += 2 {
		key := dst.Content[i].Value
		keyMap[key] = i + 1 // Store the index of the value
	}

	// Merge keys from the source node into the destination node
	for i := 0; i < len(src.Content); i += 2 {
		key := src.Content[i].Value
		if key == "<<" {
			continue // Skip merge keys as we handled them
		}

		// If the key already exists in the destination, skip it
		// Keys in mapping nodes earlier in the sequence override keys specified in later mapping nodes
		if _, exists := keyMap[key]; exists {
			continue
		}

		dst.Content = append(dst.Content, src.Content[i], src.Content[i+1])
	}
}
