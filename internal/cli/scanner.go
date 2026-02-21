package cli

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/rahuljindal1/calltree/internal/core"
	"github.com/rahuljindal1/calltree/internal/languages/typescript"
)

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

		for _, ext := range extensions {
			if strings.HasSuffix(path, ext) {
				fns, err := analyzeSingleFile(path)
				if err != nil {
					return err
				}
				for name, fn := range fns {
					functions[name] = fn
				}
				break
			}
		}
		return nil
	})

	return functions, err
}

func analyzeSingleFile(filePath string) (map[string]*core.Function, error) {
	code, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	parser := typescript.NewParser(core.ParseOptions{IncludeBuiltins: includeBuiltins})

	result, err := parser.Parse(code, filepath.Base(filePath))
	if err != nil {
		return nil, err
	}

	return result.Functions, nil
}
