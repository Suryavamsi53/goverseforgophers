package runner

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type ExecutionResult struct {
	Output string            `json:"output"`
	Error  string            `json:"error"`
	TimeMs int64             `json:"time_ms"`
	Files  map[string]string `json:"files,omitempty"` // New/modified files
}

// ExecuteCode writes the source code to a temporary file, runs it with `go run`,
// captures the standard output and standard error, and returns the result.
func ExecuteCode(ctx context.Context, files map[string]string) (*ExecutionResult, error) {
	slog.Info("Starting code execution session", "file_count", len(files))

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "goverse_run_*")
	if err != nil {
		slog.Error("Failed to create temp directory", "error", err)
		return nil, err
	}
	defer os.RemoveAll(tempDir) // Clean up afterwards

	// Write all files to the temp directory
	hasGoMod := false
	for name, content := range files {
		// Basic sanitization to prevent path traversal
		if strings.Contains(name, "..") || strings.HasPrefix(name, "/") {
			continue
		}
		if name == "go.mod" {
			hasGoMod = true
		}
		if name == "" {
			name = "main.go"
		}
		path := filepath.Join(tempDir, name)
		
		// Create parent directories if they don't exist
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			slog.Error("Failed to create parent directories", "error", err)
			return nil, err
		}
		
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			slog.Error("Failed to write file", "path", path, "error", err)
			return nil, err
		}
	}

	// Initialize Go module in the temp directory so it can run independently (if not provided)
	if !hasGoMod {
		if err := os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte("module example\n\ngo 1.21\n"), 0644); err != nil {
			slog.Error("Failed to write go.mod", "error", err)
			return nil, err
		}
	}

	// Prepare the command to run the code
	// Adding a hard timeout for execution to prevent infinite loops
	runCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(runCtx, "docker", "run", "--rm",
		"--memory", "256m",
		"--cpus", "0.5",
		"--network", "host",
		"-v", tempDir+":/app:z",
		"-w", "/app",
		"golang:1.21-alpine",
		"go", "run", ".")

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
			result.Error = "Execution timed out (10s limit). Note: For web servers, please use the Terminal!"
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

// RunCommand writes the files to a temporary directory and executes a given shell command.
func RunCommand(ctx context.Context, command string, files map[string]string) (*ExecutionResult, error) {
	slog.Info("Starting terminal command", "command", command)

	tempDir, err := os.MkdirTemp("", "goverse_term_*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir)

	hasGoMod := false
	for name, content := range files {
		if strings.Contains(name, "..") || strings.HasPrefix(name, "/") {
			continue
		}
		if name == "go.mod" {
			hasGoMod = true
		}
		if name == "" {
			name = "main.go"
		}
		path := filepath.Join(tempDir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return nil, err
		}
	}

	if !hasGoMod {
		os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte("module example\n\ngo 1.21\n"), 0644)
	}

	runCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(runCtx, "docker", "run", "--rm",
		"--memory", "256m",
		"--cpus", "0.5",
		"--network", "host",
		"-v", tempDir+":/app:z",
		"-w", "/app",
		"golang:1.21-alpine",
		"sh", "-c", command)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err = cmd.Run()
	duration := time.Since(start).Milliseconds()

	result := &ExecutionResult{
		Output: stdout.String(),
		TimeMs: duration,
		Files:  make(map[string]string),
	}

	// Read back any files created/modified by the command
	filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		// Skip go.mod if we auto-generated it and it wasn't there before
		rel, _ := filepath.Rel(tempDir, path)
		if rel == "go.mod" && !hasGoMod {
			return nil
		}
		
		contentBytes, err := os.ReadFile(path)
		if err == nil {
			result.Files[rel] = string(contentBytes)
		}
		return nil
	})

	if err != nil {
		if runCtx.Err() == context.DeadlineExceeded {
			result.Error = "Command timed out (10s limit). Note: For long-running processes, please use the Terminal!"
		} else {
			result.Error = stderr.String()
			if result.Error == "" {
				result.Error = err.Error()
			}
		}
	} else if stderr.Len() > 0 {
		result.Error = stderr.String()
	}

	return result, nil
}

// FormatCode runs gofmt on the provided code string
func FormatCode(ctx context.Context, code string) (string, error) {
	cmd := exec.CommandContext(ctx, "gofmt")
	cmd.Stdin = strings.NewReader(code)
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		if stderr.Len() > 0 {
			// e.g. main.go:3:1: expected 'package', found 'EOF'
			return "", fmt.Errorf("%s", stderr.String())
		}
		return "", err
	}
	
	return stdout.String(), nil
}
