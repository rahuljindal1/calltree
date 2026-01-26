package cli

import (
	"fmt"

	"calltree/internal/core"
)

func renderTree(tree map[string]*core.TreeNode, functions map[string]*core.Function) error {
	if jsonOutput {
		return core.PrintJSON(tree, functions, rootsOnly, depthOnly, jsonFile)
	}

	if rootsOnly {
		printRootsOnly(tree, functions)
		return nil
	}

	printAll(tree)
	return nil
}

func printRootsOnly(tree map[string]*core.TreeNode, functions map[string]*core.Function) {
	for _, name := range core.FindRoots(functions) {
		if node := tree[name]; node != nil {
			core.PrintTree(node, "", true, 0, depthOnly, showFile)
			fmt.Println()
		}
	}
}

func printAll(tree map[string]*core.TreeNode) {
	for _, node := range tree {
		core.PrintTree(node, "", true, 0, depthOnly, showFile)
		fmt.Println()
	}
}
