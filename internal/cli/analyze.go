package cli

import (
	"calltree/internal/core"
	"calltree/internal/languages/javascript"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	depthOnly  int
	jsonOutput bool
	rootsOnly  bool
	jsonFile   string
	focusFn    string
	showFile   bool
	recursive  bool

	excludeDirs []string
	extensions  []string
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

	analyzeCmd.Flags().StringVar(
		&jsonFile,
		"json-file",
		"",
		"Write JSON output to the specified file",
	)

	analyzeCmd.Flags().StringVar(
		&focusFn,
		"focus",
		"",
		"Show call tree starting from a specific function",
	)

	analyzeCmd.Flags().BoolVar(
		&showFile,
		"show-file",
		false,
		"Show source file name for each function",
	)

	analyzeCmd.Flags().BoolVarP(
		&recursive,
		"recursive",
		"r",
		false,
		"Scan directories recursively",
	)

	analyzeCmd.Flags().StringSliceVar(
		&excludeDirs,
		"exclude-dir",
		[]string{"node_modules"},
		"Directories to exclude (comma-separated)",
	)

	analyzeCmd.Flags().StringSliceVar(
		&extensions,
		"ext",
		[]string{".js"},
		"File extensions to include (comma-separated)",
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

func analyzeFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	var functions map[string]*core.Function

	if info.IsDir() {
		if !recursive {
			return fmt.Errorf("path %q is a directory, use --recursive", path)
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

		if jsonOutput {
			return core.PrintJSON(
				map[string]*core.TreeNode{focusFn: node},
				functions,
				rootsOnly,
				depthOnly,
				jsonFile,
			)
		}

		core.PrintTree(
			node,
			"",
			true,
			0,
			depthOnly,
			showFile,
		)
		fmt.Println()
		return nil
	}

	if jsonOutput {
		return core.PrintJSON(
			tree,
			functions,
			rootsOnly,
			depthOnly,
			jsonFile,
		)
	}

	if rootsOnly {
		printRootsOnly(tree, functions)
		return nil
	}

	printAll(tree)
	return nil
}

func analyzeDirectory(root string) (map[string]*core.Function, error) {
	functions := make(map[string]*core.Function)

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			for _, dir := range excludeDirs {
				if d.Name() == dir {
					return filepath.SkipDir
				}
			}
			return nil
		}

		matched := false
		for _, ext := range extensions {
			if strings.HasSuffix(path, ext) {
				matched = true
				break
			}
		}
		if !matched {
			return nil
		}

		fileFunctions, err := analyzeSingleFile(path)
		if err != nil {
			return err
		}

		for name, fn := range fileFunctions {
			functions[name] = fn
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return functions, nil
}

func analyzeSingleFile(filePath string) (map[string]*core.Function, error) {
	code, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	parser := javascript.NewParser()
	result, err := parser.Parse(code, filepath.Base(filePath))
	if err != nil {
		return nil, err
	}

	return result.Functions, nil
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
			showFile,
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
			showFile,
		)

		fmt.Println()
	}
}
