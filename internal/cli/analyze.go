package cli

import (
	"bufio"
	"calltree/internal/core"
	"calltree/internal/languages/javascript"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
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
	Use:   "analyze",
	Short: "Analyze source code (interactive mode)",

	Args: cobra.NoArgs,

	RunE: func(cmd *cobra.Command, args []string) error {
		return runInteractiveAnalyze()
	},
}

func runInteractiveAnalyze() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Calltree Interactive Analysis")
	fmt.Println("------------------------------")

	// 1. Path
	fmt.Print("Enter file or directory path: ")
	path, err := readLine(reader)
	if err != nil {
		return err
	}

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		fmt.Print("Scan recursively? (y/N): ")
		ans, _ := readLine(reader)
		recursive = strings.EqualFold(ans, "y")

		// 3. Exclude dirs
		fmt.Printf("Exclude directories (comma-separated) [%s]: ",
			strings.Join(excludeDirs, ","),
		)
		excl, _ := readLine(reader)
		if excl != "" {
			excludeDirs = splitCSV(excl)
		}

		// 4. Extensions
		fmt.Printf("File extensions [%s]: ",
			strings.Join(extensions, ","),
		)
		ext, _ := readLine(reader)
		if ext != "" {
			extensions = splitCSV(ext)
		}
	}

	// 5. Focus
	fmt.Print("Focus on function (leave empty to skip): ")
	focus, _ := readLine(reader)
	if focus != "" {
		focusFn = focus
	}

	// 6. Depth
	fmt.Printf("Max depth (-1 = unlimited) [%d]: ", depthOnly)
	depthStr, _ := readLine(reader)
	if depthStr != "" {
		if d, err := strconv.Atoi(depthStr); err == nil {
			depthOnly = d
		}
	}

	// 7. Output format
	fmt.Println("Output format:")
	fmt.Println("  1) Tree (CLI)")
	fmt.Println("  2) JSON")
	fmt.Print("Select [1]: ")
	format, _ := readLine(reader)

	if format == "2" {
		jsonOutput = true
		fmt.Print("JSON output file (leave empty for stdout): ")
		jf, _ := readLine(reader)
		jsonFile = jf
	}

	// 8. Show file names
	fmt.Print("Show file names? (y/N): ")
	sf, _ := readLine(reader)
	showFile = strings.EqualFold(sf, "y")

	fmt.Println("\nRunning analysis...\n")

	return analyzeFile(path)
}

func readLine(r *bufio.Reader) (string, error) {
	line, err := r.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
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
