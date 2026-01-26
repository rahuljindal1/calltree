package core

import "fmt"

func PrintTree(name string, node *TreeNode, prefix string, isLast bool) {
	if prefix == "" {
		fmt.Println(name)
	} else {
		branch := "├─ "
		if isLast {
			branch = "└─ "
		}
		fmt.Println(prefix + branch + name)
	}

	children := make([]string, 0, len(node.Children))
	for k := range node.Children {
		children = append(children, k)
	}

	for i, child := range children {
		last := i == len(children)-1

		nextPrefix := prefix
		if isLast {
			nextPrefix += "   "
		} else {
			nextPrefix += "│  "
		}

		PrintTree(
			child,
			node.Children[child],
			nextPrefix,
			last,
		)
	}
}
