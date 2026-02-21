package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const lastRunDir = ".calltree"
const lastRunFile = "last-run.json"

type AnalyzeConfig struct {
	Recursive       bool     `json:"recursive"`
	ExcludeDirs     []string `json:"excludeDirs"`
	Extensions      []string `json:"extensions"`
	FocusFn         string   `json:"focusFn"`
	Depth           int      `json:"depth"`
	JSON            bool     `json:"json"`
	JSONFile        string   `json:"jsonFile"`
	ShowFile        bool     `json:"showFile"`
	RootsOnly       bool     `json:"rootsOnly"`
	IncludeBuiltins bool     `json:"includeBuiltins"`
}

func saveLastRun(cfg AnalyzeConfig) error {
	if err := os.MkdirAll(lastRunDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(
		filepath.Join(lastRunDir, lastRunFile),
		data,
		0644,
	)
}

func loadLastRun() (*AnalyzeConfig, error) {
	data, err := os.ReadFile(
		filepath.Join(lastRunDir, lastRunFile),
	)
	if err != nil {
		return nil, err
	}

	var cfg AnalyzeConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func printConfig(cfg *AnalyzeConfig) {
	fmt.Println("\nLast analysis configuration:")
	fmt.Println("----------------------------")
	fmt.Printf("Recursive   : %v\n", cfg.Recursive)
	fmt.Printf("ExcludeDirs : %v\n", cfg.ExcludeDirs)
	fmt.Printf("Extensions  : %v\n", cfg.Extensions)
	fmt.Printf("FocusFn     : %s\n", cfg.FocusFn)
	fmt.Printf("Depth       : %d\n", cfg.Depth)
	fmt.Printf("JSON        : %v\n", cfg.JSON)
	fmt.Printf("JSONFile    : %s\n", cfg.JSONFile)
	fmt.Printf("ShowFile    : %v\n", cfg.ShowFile)
	fmt.Printf("RootsOnly   : %v\n", cfg.RootsOnly)
	fmt.Printf("IncludeBuiltins   : %v\n", cfg.IncludeBuiltins)
	fmt.Println("----------------------------")
}
