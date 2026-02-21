package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	depthOnly       int
	jsonOutput      bool
	rootsOnly       bool
	jsonFile        string
	focusFn         string
	showFile        bool
	recursive       bool
	rerun           bool
	includeBuiltins bool

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
	analyzeCmd.Flags().BoolVar(&includeBuiltins, "include-builtins", false, "Include language built-in calls (map, includes, Number, etc.)")

	rootCmd.AddCommand(analyzeCmd)
}

var analyzeCmd = &cobra.Command{
	Use:   "analyze <path>",
	Short: "Analyze source code, Use -h options for more details.",
	Long: `Analyze source code starting from the given file or directory.

When a single path is provided, the command runs in INTERACTIVE MODE:
- You will be guided through analysis options (depth, focus function, output format, etc.)
- No analysis is executed until you confirm the configuration

When multiple args are provided, the command runs in DIRECT MODE:
- Analysis starts immediately
- Flags are applied as provided
- No interactive prompts are shown

Examples:
analyze src/
	→ Interactive analysis for the src/ directory

analyze main.js
	→ Interactive analysis for main.go

analyze main.js -d=1
	→ Direct analysis of multiple paths
`,
	Args: cobra.MinimumNArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		// -----------------------------------
		// Interactive mode
		// -----------------------------------
		if len(args) == 1 && !hasUserFlags(cmd) {
			if rerun {
				return runLastAnalysis(path)
			}
			return runInteractiveAnalyze(path)
		}

		// -----------------------------------
		// Direct CLI mode
		// -----------------------------------
		return analyzePath(path)
	},
}

func hasUserFlags(cmd *cobra.Command) bool {
	cmd.Flags().Visit(func(f *pflag.Flag) {})
	return cmd.Flags().NFlag() > 0
}
