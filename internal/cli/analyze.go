package cli

import (
	"calltree/internal/core"
	"calltree/internal/languages/javascript"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	depthOnly  int
	jsonOutput bool
	rootsOnly  bool
)

type jsonNode struct {
	Name     string     `json:"name"`
	Children []jsonNode `json:"children,omitempty"`
}

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
		"Output call tree as JSON",
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

	if jsonOutput {
		return printJSON(tree, result.Functions)
	}

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

func printJSON(
	tree map[string]*core.TreeNode,
	functions map[string]*core.Function,
) error {

	var roots []string

	if rootsOnly {
		roots = core.FindRoots(functions)
	} else {
		for name := range tree {
			roots = append(roots, name)
		}
	}

	var output []jsonNode

	for _, name := range roots {
		node := tree[name]
		if node == nil {
			continue
		}

		output = append(
			output,
			toJSONNode(node, 0, depthOnly),
		)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	return enc.Encode(output)
}

func toJSONNode(
	node *core.TreeNode,
	currentDepth int,
	maxDepth int,
) jsonNode {

	if maxDepth >= 0 && currentDepth > maxDepth {
		return jsonNode{}
	}

	out := jsonNode{
		Name: node.Name,
	}

	for _, child := range node.Children {
		childNode := toJSONNode(
			child,
			currentDepth+1,
			maxDepth,
		)

		if childNode.Name != "" {
			out.Children = append(out.Children, childNode)
		}
	}

	return out
}
