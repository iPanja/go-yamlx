package iterator

import (
	"github.com/iPanja/go-yamlx/utils"
	"gopkg.in/yaml.v3"
)

type (
	Iterator func() (*yaml.Node, bool)
	//Predicate func(*yaml.Node) bool
)

// Entrypoint

// FromNode creates an Iterator on a single node
func FromNode(node *yaml.Node) Iterator {
	return FromNodes(node)
}

// FromNodes creates an Iterator on all given nodes
func FromNodes(nodes ...*yaml.Node) Iterator {
	index := 0

	return func() (node *yaml.Node, ok bool) {
		// We have reached the end
		ok = index < len(nodes)
		if !ok {
			return
		}

		node = nodes[index]
		index++

		return
	}
}

// ResolveAliases will replace any Alias node with its underlying node
func (iter Iterator) ResolveAliases() Iterator {
	return func() (node *yaml.Node, ok bool) {
		node, ok = iter()
		if !ok {
			return
		}

		if node.Kind == yaml.AliasNode {
			node = node.Alias
		}

		return
	}
}

// EmbedMerges will embed the aliases (or list of) and ensure subsequent merge blocks and aliases are resolved
//func (iter Iterator) EmbedMerges() Iterator {
//	return func() (node *yaml.Node, ok bool) {
//		return lazyNodeProcess(iter)
//	}
//}

type MergeState struct {
	InMergeBlock bool
	Aliases      []*yaml.Node
}

// EmbedMergesInMaps will update the contents of any maps that end up being iterated over.
//
// Their contents will reflect the final result of applying the merge key (both a single alias or a sequence are supported) if one is found
func (iter Iterator) EmbedMergesInMaps() Iterator {
	return func() (node *yaml.Node, ok bool) {
		node, ok = iter()
		if !ok {
			return nil, false
		}

		// Embeds only work within mapping nodes
		if node.Kind != yaml.MappingNode {
			return node, true
		}

		var mergeState MergeState
		processMappingNode(node, &mergeState)

		// If we're in a merge block, return a new node with merged content
		// We do not want to modify the true mapping node, but instead return a duplicate, fully rendered, node
		if mergeState.InMergeBlock {
			mergedNode := &yaml.Node{
				Kind:    node.Kind,
				Content: append([]*yaml.Node{}, node.Content...), // Copy original
			}

			// Merge alias nodes into the new node
			for _, alias := range mergeState.Aliases {
				utils.MergeNodes(mergedNode, alias) // Follows YAML embedding precedence
			}

			return mergedNode, true
		}

		// Otherwise, return the original node
		return node, true
	}
}

// processMappingNode will extract the mergeState if found within the given map
func processMappingNode(node *yaml.Node, mergeState *MergeState) {
	// Resolve aliases in the node itself
	//if node.Kind == yaml.AliasNode {
	//	node = node.Alias
	//}

	// Process the node's content to handle merge keys
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		if key.Value == "<<" {
			mergeValue := node.Content[i+1]
			if mergeValue.Kind == yaml.AliasNode {
				// <<: *alias
				mergeState.InMergeBlock = true
				mergeState.Aliases = append(mergeState.Aliases, mergeValue.Alias)
			} else if mergeValue.Kind == yaml.SequenceNode {
				// <<: [*alias_one, *alias_two, ...]
				for _, aliasNode := range mergeValue.Content {
					if aliasNode.Kind == yaml.AliasNode {
						mergeState.InMergeBlock = true
						mergeState.Aliases = append(mergeState.Aliases, aliasNode.Alias)
					}
				}
			}

			// Remove the merge key and its value (don't modify the original node)
			node.Content = append(node.Content[:i], node.Content[i+2:]...)
			i -= 2 // Adjust the index after removal
		}
	}
}

// RecurseNodes recursively traverse the nodes' children/content via BFS
//
// All aliases will be resolved and all merge keys will be embedded
func (iter Iterator) RecurseNodes() Iterator {
	var queue []*yaml.Node

	return func() (node *yaml.Node, ok bool) {
		// Check stack
		if len(queue) > 0 {
			// node, ok = FromNode(stack[len(stack)-1])()
			node = queue[0]
			queue = queue[1:]
			ok = true
		} else {
			node, ok = iter()
			if !ok {
				return
			}
		}

		// Process node - resolve aliases and embed merge keys
		node = processNode(node)

		// Dive into children
		if len(node.Content) > 0 {
			childIter := FromNodes(node.Content...).ResolveAliases().EmbedMergesInMaps()

			for child, ok := childIter(); ok; child, ok = childIter() {
				queue = append(queue, child)
			}
		}

		return
	}
}

func processNode(node *yaml.Node) *yaml.Node {
	// Resolve aliases for the current node
	if node.Kind == yaml.AliasNode {
		node = node.Alias
	}

	// Embed merge keys if the node is a mapping node
	if node.Kind == yaml.MappingNode {
		var mergeState MergeState
		processMappingNode(node, &mergeState)
		if mergeState.InMergeBlock {
			mergedNode := &yaml.Node{
				Kind:    node.Kind,
				Content: append([]*yaml.Node{}, node.Content...), // Copy original content
			}
			for _, alias := range mergeState.Aliases {
				utils.MergeNodes(mergedNode, alias)
			}
			node = mergedNode
		}
	}

	return node
}
