package main

import (
	"calltree/internal/core"
	"calltree/internal/languages/javascript"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("calltree v0.1")
		fmt.Println("Usage:")
		fmt.Println("  calltree <file>")
		os.Exit(0)
	}

	filePath := os.Args[1]

	fmt.Println("Running calltree on file:")
	fmt.Println(filePath)

	code, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	parser := javascript.NewParser()

	result, err := parser.Parse(code)
	if err != nil {
		panic(err)
	}

	fmt.Println("Language:", result.Language)
	fmt.Println("Functions found:", len(result.Functions))

	tree := core.BuildCallTree(result.Functions)

	for name, node := range tree {
		core.PrintTree(name, node, "", true)
		fmt.Println()
	}
}
