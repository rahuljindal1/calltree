package core

import "fmt"

func PrintTree(node *TreeNode, prefix string, isLast bool, currentDepth int,
	maxDepth int) {
	if maxDepth >= 0 && currentDepth > maxDepth {
		return
	}

	if prefix == "" {
		fmt.Println(node.Name)
	} else {
		branch := "├─ "
		if isLast {
			branch = "└─ "
		}
		fmt.Println(prefix + branch + node.Name)
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
			child,
			nextPrefix,
			last,
			currentDepth+1,
			maxDepth,
		)
	}
}
