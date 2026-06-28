package runner

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type ExecutionResult struct {
	Output string `json:"output"`
	Error  string `json:"error"`
	TimeMs int64  `json:"time_ms"`
}

// ExecuteCode writes the source code to a temporary file, runs it with `go run`,
// captures the standard output and standard error, and returns the result.
func ExecuteCode(ctx context.Context, code string) (*ExecutionResult, error) {
	slog.Info("Starting code execution session", "code_length", len(code))

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "goverse_run_*")
	if err != nil {
		slog.Error("Failed to create temp directory", "error", err)
		return nil, err
	}
	defer os.RemoveAll(tempDir) // Clean up afterwards

	// If the code doesn't contain a main function, we inject a dummy one so it compiles.
	// This supports running LeetCode-style function snippets.
	importFmt := false
	if !strings.Contains(code, "func main()") {
		code += "\n\nfunc main() {\n\t// Automatically added main function for testing\n"
		if strings.Contains(code, "func Solve(") {
			code += "\tfmt.Println(\"Output:\", Solve(\"\"))\n"
			importFmt = true
		} else {
			code += "\tfmt.Println(\"Code compiled successfully, but no main() or Solve() function was found to execute.\")\n"
			importFmt = true
		}
		code += "}\n"
	}

	if importFmt && !strings.Contains(code, "\"fmt\"") {
		code = strings.Replace(code, "package main", "package main\nimport \"fmt\"\n", 1)
	}

	// Write code to main.go in the temp directory
	mainPath := filepath.Join(tempDir, "main.go")
	if err := os.WriteFile(mainPath, []byte(code), 0644); err != nil {
		return nil, err
	}

	// Initialize Go module in the temp directory so it can run independently
	modCmd := exec.CommandContext(ctx, "go", "mod", "init", "run")
	modCmd.Dir = tempDir
	if err := modCmd.Run(); err != nil {
		return nil, err
	}

	// Prepare the command to run the code
	// Adding a hard timeout for execution to prevent infinite loops
	runCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(runCtx, "go", "run", "main.go")
	cmd.Dir = tempDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err = cmd.Run()
	duration := time.Since(start).Milliseconds()

	slog.Info("Code execution completed", "duration_ms", duration, "success", err == nil)

	result := &ExecutionResult{
		Output: stdout.String(),
		TimeMs: duration,
	}

	// If there's an error, it could be a compile error or runtime panic
	if err != nil {
		if runCtx.Err() == context.DeadlineExceeded {
			result.Error = "Execution timed out (5s limit)"
		} else {
			result.Error = stderr.String()
			if result.Error == "" {
				result.Error = err.Error()
			}
		}
	} else if stderr.Len() > 0 {
		// Some programs might write to stderr without failing
		result.Error = stderr.String()
	}

	return result, nil
}
