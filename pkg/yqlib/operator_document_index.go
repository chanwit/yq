package yqlib

import (
	"container/list"
	"fmt"

	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func getDocumentIndexOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)
		node := &yaml.Node{Kind: yaml.ScalarNode, Value: fmt.Sprintf("%v", candidate.Document), Tag: "!!int"}
		scalar := candidate.CreateChild(nil, node)
		results.PushBack(scalar)
	}
	return context.ChildContext(results), nil
}
