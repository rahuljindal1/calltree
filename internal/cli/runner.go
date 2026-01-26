package cli

import (
	"fmt"
	"os"

	"calltree/internal/core"
)

func analyzeFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	var functions map[string]*core.Function

	if info.IsDir() {
		if !recursive {
			return fmt.Errorf("path %q is a directory, enable recursive scan", path)
		}
		functions, err = analyzeDirectory(path)
	} else {
		functions, err = analyzeSingleFile(path)
	}
	if err != nil {
		return err
	}

	tree := core.BuildCallTree(functions)

	if focusFn != "" {
		node, ok := tree[focusFn]
		if !ok {
			return fmt.Errorf("function %q not found", focusFn)
		}
		return renderTree(map[string]*core.TreeNode{focusFn: node}, functions)
	}

	return renderTree(tree, functions)
}
