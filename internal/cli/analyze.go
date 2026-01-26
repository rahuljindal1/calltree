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

		filePath := args[0]

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

		for name, node := range tree {
			core.PrintTree(
				name,
				node,
				"",
				true,
			)
			fmt.Println()
		}

		return nil
	},
}
