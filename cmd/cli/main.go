package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/suryavamsivaggu/goverse/pkg/runner"
)

func main() {
	// Configure slog
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	slog.SetDefault(slog.New(handler))

	// Parse flags
	fileFlag := flag.String("file", "", "Path to the Go file to execute")
	flag.Parse()

	if *fileFlag == "" {
		fmt.Println("Usage: goverse-cli -file <path_to_go_file>")
		os.Exit(1)
	}

	// Read file
	codeBytes, err := os.ReadFile(*fileFlag)
	if err != nil {
		slog.Error("Failed to read file", "file", *fileFlag, "error", err)
		os.Exit(1)
	}

	code := string(codeBytes)
	
	slog.Info("Executing file locally", "file", *fileFlag)

	// Execute Code
	ctx := context.Background()
	result, err := runner.ExecuteCode(ctx, map[string]string{"main.go": code, "go.mod": "module example\ngo 1.21\n"})
	if err != nil {
		slog.Error("Execution failed completely", "error", err)
		os.Exit(1)
	}

	// Print results
	fmt.Println("========== EXECUTION RESULT ==========")
	if result.Error != "" {
		fmt.Printf("ERROR:\n%s\n", result.Error)
	} else {
		fmt.Printf("OUTPUT:\n%s\n", result.Output)
	}
	fmt.Printf("TIME: %d ms\n", result.TimeMs)
	fmt.Println("======================================")
}
