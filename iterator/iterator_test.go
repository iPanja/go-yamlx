package iterator

import (
	"fmt"
	"github.com/iPanja/go-yamlx/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gopkg.in/yaml.v3"
)

var _ = Describe("Iterator", func() {
	Describe("FromNode", func() {
		It("returns the root node", func() {
			node := &yaml.Node{}
			next := FromNode(node)

			item, ok := next()
			Expect(item).To(Equal(node))
			Expect(ok).To(BeTrue())

			item, ok = next()
			Expect(item).To(BeNil())
			Expect(ok).To(BeFalse())
		})
	})

	DescribeTable("FromNodes",
		func(yaml string, values ...*yaml.Node) {
			doc := utils.ParseYaml(yaml)
			it := FromNodes(doc.Content...)

			for _, value := range values {
				node, ok := it()

				Expect(ok).To(BeTrue())
				Expect(node.Kind).To(Equal(value.Kind))
				Expect(node.Value).To(Equal(value.Value))
			}

			_, ok := it()
			Expect(ok).To(BeFalse())
		},
		// Description, yaml string, values ...*yaml.Node
		// These tests don't do a lot since we aren't going through their children
		Entry("scalar", "a", scalarNode("a")),
		Entry("sequence", "[a, b, c]", &yaml.Node{Kind: yaml.SequenceNode, Value: ""}),
		Entry("map", "{a: b}", &yaml.Node{Kind: yaml.MappingNode, Value: ""}),
	)

	Describe("ResolveAlias", func() {
		It("Both alias nodes return the underlying scalar node", func() {
			targetScalarNode := &yaml.Node{Kind: yaml.ScalarNode, Value: "destination"}
			aliasNode := &yaml.Node{Kind: yaml.AliasNode, Value: "", Alias: targetScalarNode}

			content := []*yaml.Node{
				targetScalarNode,
				aliasNode,
				aliasNode,
			}

			iter := FromNodes(content...).ResolveAliases()

			for i := 0; i < 3; i++ {
				item, ok := iter()
				Expect(item).To(Equal(targetScalarNode))
				Expect(ok).To(BeTrue())
			}

			item, ok := iter()
			Expect(item).To(BeNil())
			Expect(ok).To(BeFalse())
		})
	})

	Describe("ResolveAlias", func() {
		It("Both alias nodes return the underlying scalar node", func() {
			targetScalarNode := &yaml.Node{Kind: yaml.ScalarNode, Value: "destination"}
			aliasNode := &yaml.Node{Kind: yaml.AliasNode, Value: "", Alias: targetScalarNode}

			content := []*yaml.Node{
				targetScalarNode,
				aliasNode,
				aliasNode,
			}

			iter := FromNodes(content...).ResolveAliases()

			for i := 0; i < 3; i++ {
				item, ok := iter()
				Expect(item).To(Equal(targetScalarNode))
				Expect(ok).To(BeTrue())
			}

			item, ok := iter()
			Expect(item).To(BeNil())
			Expect(ok).To(BeFalse())
		})
	})

	//Describe("EmbedMerges", func() {
	//	It("The merge key should be replaced with the alias' content", func() {
	//		targetMergeContents := &yaml.Node{
	//			Kind:  yaml.SequenceNode,
	//			Value: "",
	//			Content: []*yaml.Node{
	//				&yaml.Node{Kind: yaml.ScalarNode, Value: "destinationKey"},
	//				&yaml.Node{Kind: yaml.ScalarNode, Value: "destinationValue"},
	//			},
	//		}
	//
	//		content := []*yaml.Node{
	//			scalarNode("a"),
	//			&yaml.Node{Kind: yaml.ScalarNode, Value: "<<", Tag: "!!merge"},
	//			&yaml.Node{Kind: yaml.AliasNode, Value: "", Alias: targetMergeContents},
	//			scalarNode("b"),
	//		}
	//
	//		result := []*yaml.Node{
	//			scalarNode("a"),
	//			scalarNode("destinationKey"),
	//			scalarNode("destinationValue"),
	//			scalarNode("b"),
	//		}
	//
	//		iter := FromNodes(content...).EmbedMerges()
	//
	//		for item, ok := iter(); ok && len(result) > 0; item, ok = iter() {
	//			nextResult := result[0]
	//			result = result[1:]
	//
	//			Expect(item).To(Equal(nextResult))
	//			Expect(ok).To(BeTrue())
	//		}
	//
	//		item, ok := iter()
	//		Expect(item).To(BeNil())
	//		Expect(ok).To(BeFalse())
	//	})
	//})

	DescribeTable("RecurseNodes",
		func(yaml string, values ...*yaml.Node) {
			doc := toYAML(yaml)
			next := FromNode(doc).RecurseNodes()

			for _, value := range values {
				node, ok := next()

				Expect(ok).To(BeTrue())
				Expect(node.Kind).To(Equal(value.Kind))
				Expect(node.Value).To(Equal(value.Value))
			}

			_, ok := next()
			Expect(ok).To(BeFalse())
		},

		Entry("scalar", "a",
			docNode, scalarNode("a")),

		Entry("sequence", "[a, b, c]",
			docNode, seqNode, scalarNode("a"), scalarNode("b"), scalarNode("c")),

		Entry("map", "{a: b}",
			docNode, mapNode, scalarNode("a"), scalarNode("b")),
	)

	Describe("RecurseNodes", func() {
		It("Extensive Test", func() {
			extensiveData := `
doc:
  var: &name user
  block: &bv
    - u
    - v
    - *name
  map:
    - a: b
    - name: *name
      <<: *bv
      d: e
  seq:
    - y
    - z
    - *name
`
			result := []*yaml.Node{}

			doc := utils.ParseYaml(extensiveData)
			iter := FromNode(doc).EmbedMerges().RecurseNodes()

			for item, ok := iter(); ok; item, ok = iter() {
				result = append(result, item)
				fmt.Printf("[%v]\t%v\n", item.Tag, item.Value)
				// Expect(item).To(Equal(nextResult))
				// Expect(ok).To(BeTrue())
			}

			item, ok := iter()
			Expect(item).To(BeNil())
			Expect(ok).To(BeFalse())
		})
	})

	//DescribeTable("Values",
	//	func(yaml string, values ...*yaml.Node) {
	//		doc := toYAML(yaml)
	//		next := FromNode(doc).
	//			Values(). // the root is the document node
	//			Values()
	//
	//		for _, value := range values {
	//			node, ok := next()
	//
	//			Expect(ok).To(BeTrue())
	//			Expect(node.Kind).To(Equal(value.Kind))
	//			Expect(node.Value).To(Equal(value.Value))
	//		}
	//
	//		_, ok := next()
	//		Expect(ok).To(BeFalse())
	//
	//	},
	//
	//	Entry("scalar", nil /* no values */),
	//
	//	Entry("sequence", "[a, b, c]",
	//		scalarNode("a"), scalarNode("b"), scalarNode("c")),
	//
	//	Entry("map", "a: b\nc: d",
	//		scalarNode("a"), scalarNode("b"), scalarNode("c"), scalarNode("d"),
	//	),
	//)
	//
	//Describe("Filter", func() {
	//	It("passes items through satisfying the predicate", func() {
	//		next := FromNode(docNode).Filter(All)
	//		node, ok := next()
	//
	//		Expect(ok).To(BeTrue())
	//		Expect(node).To(Equal(docNode))
	//	})
	//
	//	It("does not pass items that do not satisfy the predicate", func() {
	//		next := FromNode(docNode).Filter(None)
	//		_, ok := next()
	//		Expect(ok).To(BeFalse())
	//	})
	//
	//	It("predicate is not invoked when there are no items", func() {
	//		empty := Iterator(func() (*yaml.Node, bool) {
	//			return nil, false
	//		})
	//
	//		next := empty.Filter(func(*yaml.Node) bool {
	//			Fail("unexpected invocation of the filter")
	//			return true
	//		})
	//
	//		_, ok := next()
	//		Expect(ok).To(BeFalse())
	//	})
	//})
	//
	//Describe("MapKeys", func() {
	//	It("returns the keys of a map", func() {
	//		next := FromNode(toYAML("a: b\nc: d\ne: f")).
	//			RecurseNodes().
	//			Filter(WithKind(yaml.MappingNode)).
	//			MapKeys()
	//
	//		for _, value := range []string{"a", "c", "e"} {
	//			node, ok := next()
	//			Expect(ok).To(BeTrue())
	//			Expect(node.Value).To(Equal(value))
	//		}
	//
	//		_, ok := next()
	//		Expect(ok).To(BeFalse())
	//	})
	//
	//	It("returns nothing for sequences", func() {
	//		next := FromNode(toYAML("[a, b, c, d]")).
	//			RecurseNodes().
	//			Filter(WithKind(yaml.SequenceNode)).
	//			MapKeys()
	//
	//		_, ok := next()
	//		Expect(ok).To(BeFalse())
	//	})
	//})
	//
	//Describe("MapValues", func() {
	//	It("returns the keys of a map", func() {
	//		next := FromNode(toYAML("a: b\nc: d\ne: f")).
	//			RecurseNodes().
	//			Filter(WithKind(yaml.MappingNode)).
	//			MapValues()
	//
	//		for _, value := range []string{"b", "d", "f"} {
	//			node, ok := next()
	//			Expect(ok).To(BeTrue())
	//			Expect(node.Value).To(Equal(value))
	//		}
	//
	//		_, ok := next()
	//		Expect(ok).To(BeFalse())
	//	})
	//
	//	It("returns nothing for sequences", func() {
	//		next := FromNode(toYAML("[a, b, c, d]")).
	//			RecurseNodes().
	//			Filter(WithKind(yaml.SequenceNode)).
	//			MapValues()
	//
	//		_, ok := next()
	//		Expect(ok).To(BeFalse())
	//	})
	//})
	//
	//Describe("Iterate", func() {
	//	It("custom iterators can be supplied", func() {
	//		repeater := func(next Iterator) Iterator {
	//			return func() (node *yaml.Node, ok bool) {
	//				node, ok = next()
	//				if ok {
	//					node = scalarNode(strings.Repeat(node.Value, 2))
	//				}
	//				return
	//			}
	//		}
	//
	//		next := FromNodes(scalarNode("a")).
	//			Iterate(repeater).
	//			Iterate(repeater)
	//
	//		node, ok := next()
	//		Expect(ok).To(BeTrue())
	//		Expect(node.Value).To(Equal("aaaa"))
	//	})
	//})
	//
	//Describe("ValuesForMap", func() {
	//	It("returns the values of a map matching the key/value predicates", func() {
	//		next := FromNode(toYAML("a: b\nc: d\ne: f")).
	//			RecurseNodes().
	//			Filter(WithKind(yaml.MappingNode)).
	//			ValuesForMap(All, func(node *yaml.Node) bool {
	//				return node.Value == "d"
	//			})
	//
	//		node, ok := next()
	//		Expect(ok).To(BeTrue())
	//		Expect(node.Value).To(Equal("d"))
	//
	//		_, ok = next()
	//		Expect(ok).To(BeFalse())
	//	})
	//
	//	It("returns nothing for sequences", func() {
	//		next := FromNode(toYAML("[a, b, c, d]")).
	//			RecurseNodes().
	//			Filter(WithKind(yaml.SequenceNode)).
	//			ValuesForMap(All, All)
	//
	//		_, ok := next()
	//		Expect(ok).To(BeFalse())
	//	})
	//})
	//
	//Describe("FromIterators", func() {
	//	It("merges multiple iterators into a single stream", func() {
	//		next := FromIterators(
	//			FromNode(&yaml.Node{Value: "a"}),
	//			FromNode(&yaml.Node{Value: "b"}),
	//			FromNode(&yaml.Node{Value: "c"}),
	//		)
	//
	//		for _, value := range []string{"a", "b", "c"} {
	//			node, ok := next()
	//			Expect(ok).To(BeTrue(), value+" to be present")
	//			Expect(node.Value).To(Equal(value))
	//		}
	//
	//		_, ok := next()
	//		Expect(ok).To(BeFalse())
	//	})
	//})
})

var mapNode = &yaml.Node{
	Kind: yaml.MappingNode,
}

var seqNode = &yaml.Node{
	Kind: yaml.SequenceNode,
}

var docNode = &yaml.Node{
	Kind: yaml.DocumentNode,
}

func scalarNode(val string) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: val,
	}
}

func aliasNode(targetNode *yaml.Node) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.AliasNode,
		Value: "",
		Alias: targetNode,
	}
}

func toYAML(s string) *yaml.Node {
	var node yaml.Node

	err := yaml.Unmarshal([]byte(s), &node)
	ExpectWithOffset(1, err).ToNot(HaveOccurred())

	return &node
}
