package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func runLastAnalysis(path string) error {
	cfg, err := loadLastRun()
	if err != nil {
		return fmt.Errorf(
			"no previous run found (run analyze once first)",
		)
	}

	applyConfig(cfg)
	return analyzePath(path)
}

func runInteractiveAnalyze(path string) error {
	reader := bufio.NewReader(os.Stdin)

	// 🔁 RERUN PROMPT (FIRST)
	if cfg, err := loadLastRun(); err == nil {
		printConfig(cfg)

		fmt.Print("Re-run this configuration? (y/n): ")
		ans, _ := readLine(reader)

		if ans == "" || strings.EqualFold(ans, "y") {
			applyConfig(cfg)
			fmt.Println("Re-running previous analysis...")
			return analyzePath(path)
		}

		fmt.Println("Starting a new interactive analysis...")
	}

	fmt.Println("Calltree Interactive Analysis")
	fmt.Println("------------------------------")

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		fmt.Print("Scan recursively? (Y/n), Default = Y: ")
		ans, _ := readLine(reader)

		if ans == "" || strings.EqualFold(ans, "y") {
			recursive = true
		} else {
			recursive = false
		}

		fmt.Printf("Exclude directories [%s]: ", strings.Join(excludeDirs, ","))
		if v, _ := readLine(reader); v != "" {
			excludeDirs = splitCSV(v)
		}

		fmt.Printf("File extensions [%s]: ", strings.Join(extensions, ","))
		if v, _ := readLine(reader); v != "" {
			extensions = splitCSV(v)
		}
	}

	fmt.Print("Focus on function (optional): ")
	if v, _ := readLine(reader); v != "" {
		focusFn = v
	}

	fmt.Printf("Max depth (-1 = unlimited) [%d]: ", depthOnly)
	if v, _ := readLine(reader); v != "" {
		if d, err := strconv.Atoi(v); err == nil {
			depthOnly = d
		}
	}

	fmt.Println("Output format:")
	fmt.Println("  1) Tree")
	fmt.Println("  2) JSON")
	fmt.Print("Select [1]: ")
	if v, _ := readLine(reader); v == "2" {
		jsonOutput = true
		fmt.Print("JSON output file (optional): ")
		jsonFile, _ = readLine(reader)
	}

	fmt.Print("Show file names? (y/N): ")
	showFile = strings.EqualFold(mustRead(reader), "y")

	fmt.Print("Include built-in calls (map, includes, Number, etc.)? (y/N), Default = NO: ")
	includeBuiltins = strings.EqualFold(mustRead(reader), "y")

	cfg := AnalyzeConfig{
		Recursive:       recursive,
		ExcludeDirs:     excludeDirs,
		Extensions:      extensions,
		FocusFn:         focusFn,
		Depth:           depthOnly,
		JSON:            jsonOutput,
		JSONFile:        jsonFile,
		ShowFile:        showFile,
		RootsOnly:       rootsOnly,
		IncludeBuiltins: includeBuiltins,
	}

	_ = saveLastRun(cfg)

	return analyzePath(path)
}

func applyConfig(cfg *AnalyzeConfig) {
	recursive = cfg.Recursive
	excludeDirs = cfg.ExcludeDirs
	extensions = cfg.Extensions
	focusFn = cfg.FocusFn
	depthOnly = cfg.Depth
	jsonOutput = cfg.JSON
	jsonFile = cfg.JSONFile
	showFile = cfg.ShowFile
	rootsOnly = cfg.RootsOnly
	includeBuiltins = cfg.IncludeBuiltins
}

func readLine(r *bufio.Reader) (string, error) {
	line, err := r.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func mustRead(r *bufio.Reader) string {
	v, _ := readLine(r)
	return v
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
