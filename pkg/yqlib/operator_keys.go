package yqlib

import (
	"container/list"
	"fmt"

	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func keysOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- keysOperator")

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		node := unwrapDoc(candidate.Node)
		var targetNode *yaml.Node
		if node.Kind == yaml.MappingNode {
			targetNode = getMapKeys(node)
		} else if node.Kind == yaml.SequenceNode {
			targetNode = getIndicies(node)
		} else {
			return Context{}, fmt.Errorf("Cannot get keys of %v, keys only works for maps and arrays", node.Tag)
		}

		result := candidate.CreateChild(nil, targetNode)
		results.PushBack(result)
	}

	return context.ChildContext(results), nil
}

func getMapKeys(node *yaml.Node) *yaml.Node {
	contents := make([]*yaml.Node, 0)
	for index := 0; index < len(node.Content); index = index + 2 {
		contents = append(contents, node.Content[index])
	}
	return &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq", Content: contents}
}

func getIndicies(node *yaml.Node) *yaml.Node {
	var contents = make([]*yaml.Node, len(node.Content))

	for index := range node.Content {
		contents[index] = &yaml.Node{
			Kind:  yaml.ScalarNode,
			Tag:   "!!int",
			Value: fmt.Sprintf("%v", index),
		}
	}

	return &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq", Content: contents}
}
