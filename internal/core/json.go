package core

import (
	"encoding/json"
	"os"
)

type JSONNode struct {
	Name     string     `json:"name"`
	Children []JSONNode `json:"children,omitempty"`
}

func PrintJSON(
	tree map[string]*TreeNode,
	functions map[string]*Function,
	rootsOnly bool,
	maxDepth int,
) error {

	var roots []string

	if rootsOnly {
		roots = FindRoots(functions)
	} else {
		for name := range tree {
			roots = append(roots, name)
		}
	}

	var output []JSONNode

	for _, name := range roots {
		node := tree[name]
		if node == nil {
			continue
		}

		output = append(output, toJSONNode(node, 0, maxDepth))
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(output)
}

func toJSONNode(
	node *TreeNode,
	currentDepth int,
	maxDepth int,
) JSONNode {

	if maxDepth >= 0 && currentDepth > maxDepth {
		return JSONNode{}
	}

	out := JSONNode{Name: node.Name}

	for _, child := range node.Children {
		childNode := toJSONNode(child, currentDepth+1, maxDepth)
		if childNode.Name != "" {
			out.Children = append(out.Children, childNode)
		}
	}

	return out
}
