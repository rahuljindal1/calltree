package core

import "fmt"

func PrintTree(
	node *TreeNode,
	prefix string,
	isLast bool,
	currentDepth int,
	maxDepth int,
	showFile bool,
) {
	if maxDepth >= 0 && currentDepth > maxDepth {
		return
	}

	label := node.Name
	if showFile && node.File != "" {
		label = fmt.Sprintf("%s (%s)", node.Name, node.File)
	}

	if prefix == "" {
		fmt.Println(label)
	} else {
		branch := "├─ "
		if isLast {
			branch = "└─ "
		}
		fmt.Println(prefix + branch + label)
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
			showFile,
		)
	}
}
