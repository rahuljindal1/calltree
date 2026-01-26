package cli

import (
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
	rerun      bool

	excludeDirs []string
	extensions  []string
)

func init() {
	analyzeCmd.Flags().IntVarP(&depthOnly, "depth", "d", -1, "Maximum call depth")
	analyzeCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output call tree as JSON")
	analyzeCmd.Flags().BoolVar(&rootsOnly, "roots-only", false, "Print only root functions")
	analyzeCmd.Flags().StringVar(&jsonFile, "json-file", "", "Write JSON output to file")
	analyzeCmd.Flags().StringVar(&focusFn, "focus", "", "Focus on a specific function")
	analyzeCmd.Flags().BoolVar(&showFile, "show-file", false, "Show source file name")
	analyzeCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Scan directories recursively")
	analyzeCmd.Flags().StringSliceVar(&excludeDirs, "exclude-dir", []string{}, "Directories to exclude")
	analyzeCmd.Flags().StringSliceVar(&extensions, "ext", []string{}, "File extensions to include")
	analyzeCmd.Flags().BoolVar(&rerun, "rerun", false, "Re-run the last analysis configuration")

	rootCmd.AddCommand(analyzeCmd)
}

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze source code (interactive mode)",
	Args:  cobra.MaximumNArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		// -----------------------------------
		// Interactive mode
		// -----------------------------------
		if len(args) == 0 {
			if rerun {
				return runLastAnalysis()
			}
			return runInteractiveAnalyze()
		}

		// -----------------------------------
		// Direct CLI mode
		// -----------------------------------
		path := args[0]

		return analyzeFile(path)
	},
}
