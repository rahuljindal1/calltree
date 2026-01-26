package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func runInteractiveAnalyze() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Calltree Interactive Analysis")
	fmt.Println("------------------------------")

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

	fmt.Println("\nRunning analysis...\n")

	cmd := buildCommand(path)

	fmt.Println("Reusable command:")
	fmt.Println(cmd)
	fmt.Println()

	return analyzeFile(path)
}

func buildCommand(path string) string {
	args := []string{"go run cmd/calltree/main.go", "analyze"}

	args = append(args, strconv.Quote(path))

	if recursive {
		args = append(args, "--recursive")
	}

	if len(excludeDirs) > 0 {
		args = append(args, "--exclude-dir="+strings.Join(excludeDirs, ","))
	}

	if len(extensions) > 0 {
		args = append(args, "--ext="+strings.Join(extensions, ","))
	}

	if focusFn != "" {
		args = append(args, "--focus="+focusFn)
	}

	if depthOnly != 0 {
		args = append(args, "--depth="+strconv.Itoa(depthOnly))
	}

	if jsonOutput {
		args = append(args, "--json")

		if jsonFile != "" {
			args = append(args, "--json-file="+jsonFile)
		}
	}

	if showFile {
		args = append(args, "--show-file")
	}

	return strings.Join(args, " ")
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
