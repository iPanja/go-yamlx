package iterator

import (
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
// without modifying the actual nodes.
func FromNodes(nodes ...*yaml.Node) Iterator {
	index := 0
	//var mergeNodes []*yaml.Node
	//var mergeIndex int

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

func (iter Iterator) EmbedMerges() Iterator {
	var mergeIter Iterator

	return func() (node *yaml.Node, ok bool) {
		// Check if we are currently iterating through a merge block's contents
		if mergeIter != nil {
			node, ok = mergeIter()
			if ok {
				return
			}

			mergeIter = nil
		}

		// Normal iteration
		node, ok = iter()
		if !ok {
			return
		}

		// Instead return first node of merge sequence, and use it for the future
		// Key: Scalar Node (Tag == !!merge)
		// Value: Alias (TODO: could be a sequence of aliases?)
		if node.Tag == "!!merge" {
			alias, aok := iter()
			if !aok {
				return // Yikes!
			}

			mergeIter = FromNodes(alias.Alias.Content...)

			node, ok = mergeIter()
			if !ok {
				node, ok = iter() // TODO: see if this is correct behavior...
			}
		}

		return
	}
}

// RecurseNodes recursively expands out on all children (sequence, mapping)
//
// # Searches DFS style
//
// - Resolves aliases
//
// - Embeds merge blocks
func (iter Iterator) RecurseNodes() Iterator {
	var stack []*yaml.Node

	return func() (node *yaml.Node, ok bool) {
		// Check stack
		if len(stack) > 0 {
			// node, ok = FromNode(stack[len(stack)-1])()
			node = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			ok = true
		} else {
			node, ok = iter()
			if !ok {
				return
			}
		}

		// Dive into children
		if len(node.Content) > 0 {
			childIter := FromNodes(node.Content...)
			for child, ok := childIter(); ok; child, ok = childIter() {
				stack = append([]*yaml.Node{child}, stack...)
			}
		}

		return
	}
}
