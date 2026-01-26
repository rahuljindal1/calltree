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

	for i, child := range node.Children {
		last := i == len(node.Children)-1

		nextPrefix := prefix
		if isLast {
			nextPrefix += "   "
		} else {
			nextPrefix += "│  "
		}

		PrintTree(
			child.Name,
			child,
			nextPrefix,
			last,
		)
	}
}
