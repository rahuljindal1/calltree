package cli

import (
	"calltree/internal/core"
	"calltree/internal/languages/javascript"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	depthOnly  int
	jsonOutput bool
	rootsOnly  bool
)

func init() {
	analyzeCmd.Flags().IntVarP(
		&depthOnly,
		"depth",
		"d",
		-1,
		"Maximum call depth",
	)

	analyzeCmd.Flags().BoolVar(
		&jsonOutput,
		"json",
		false,
		"Output as JSON",
	)

	analyzeCmd.Flags().BoolVar(
		&rootsOnly,
		"roots-only",
		false,
		"Print only root functions",
	)

	rootCmd.AddCommand(analyzeCmd)
}

var analyzeCmd = &cobra.Command{
	Use:   "analyze <file>",
	Short: "Analyze a source file",
	Args:  cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		return analyzeFile(args[0])
	},
}

func analyzeFile(filePath string) error {
	code, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	parser := javascript.NewParser()

	result, err := parser.Parse(code)
	if err != nil {
		return err
	}

	tree := core.BuildCallTree(result.Functions)

	if rootsOnly {
		printRootsOnly(tree, result.Functions)
		return nil
	}

	printAll(tree)
	return nil
}

func printRootsOnly(
	tree map[string]*core.TreeNode,
	functions map[string]*core.Function,
) {
	roots := core.FindRoots(functions)

	for _, name := range roots {
		node := tree[name]
		if node == nil {
			continue
		}

		core.PrintTree(
			node,
			"",
			true,
			0,
			depthOnly,
		)

		fmt.Println()
	}
}

func printAll(tree map[string]*core.TreeNode) {
	for _, node := range tree {
		core.PrintTree(
			node,
			"",
			true,
			0,
			depthOnly,
		)

		fmt.Println()
	}
}
