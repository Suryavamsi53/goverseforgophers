package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/suryavamsivaggu/goverse/internal/domain"
)

type TestResult struct {
	Passed         bool   `json:"passed"`
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
	ActualOutput   string `json:"actual_output"`
	Error          string `json:"error,omitempty"`
}

type EvaluationResult struct {
	Success     bool         `json:"success"`
	TestResults []TestResult `json:"test_results"`
	TimeMs      int64        `json:"time_ms"`
	SystemError string       `json:"system_error,omitempty"`
}

// EvaluateProblem runs the user's code against the provided test cases.
// It assumes the user code provides a specific function that we test.
// For MVP, we will append a simple main function that runs the test cases
// and outputs JSON, which we then parse.
func EvaluateProblem(ctx context.Context, code string, problem *domain.PracticeProblem) (*EvaluationResult, error) {
	slog.Info("Starting problem evaluation", "problem_slug", problem.Slug, "code_length", len(code))

	tempDir, err := os.MkdirTemp("", "goverse_eval_*")
	if err != nil {
		slog.Error("Failed to create temp directory for evaluation", "error", err)
		return nil, err
	}
	defer os.RemoveAll(tempDir)

	// We need to inject a test runner into the code.
	// We'll replace their "func main" or assume they don't have one,
	// but for simplicity in MVP, we just append our own main block.
	
	// A better approach for a real system:
	// We wrap their code as a package, and write a separate test file `main_test.go`.
	// For this demo, let's write their code to `solution.go` and our test harness to `main.go`.

	solutionPath := filepath.Join(tempDir, "solution.go")
	// Make sure they declare package main
	if !strings.Contains(code, "package main") {
		code = "package main\n\n" + code
	}
	// Rename main function so it doesn't conflict if they wrote one
	code = strings.Replace(code, "func main()", "func userMain()", 1)
	
	if err := os.WriteFile(solutionPath, []byte(code), 0644); err != nil {
		return nil, err
	}

	// Generate main.go harness
	harnessCode := generateHarness(problem.TestCases)
	mainPath := filepath.Join(tempDir, "main.go")
	if err := os.WriteFile(mainPath, []byte(harnessCode), 0644); err != nil {
		return nil, err
	}

	modCmd := exec.CommandContext(ctx, "go", "mod", "init", "run")
	modCmd.Dir = tempDir
	if err := modCmd.Run(); err != nil {
		return nil, err
	}

	runCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(runCtx, "go", "run", "main.go", "solution.go")
	cmd.Dir = tempDir

	start := time.Now()
	out, err := cmd.CombinedOutput()
	duration := time.Since(start).Milliseconds()

	slog.Info("Problem evaluation completed", "duration_ms", duration, "problem_slug", problem.Slug, "success", err == nil)

	evalResult := &EvaluationResult{
		TimeMs: duration,
	}

	if err != nil {
		evalResult.SystemError = string(out)
		if evalResult.SystemError == "" {
			evalResult.SystemError = err.Error()
		}
		return evalResult, nil
	}

	// Parse JSON output from the harness
	if err := json.Unmarshal(out, &evalResult.TestResults); err != nil {
		evalResult.SystemError = "Failed to parse test results: " + err.Error() + "\nOutput: " + string(out)
		return evalResult, nil
	}

	evalResult.Success = true
	for _, tr := range evalResult.TestResults {
		if !tr.Passed {
			evalResult.Success = false
			break
		}
	}

	return evalResult, nil
}

func generateHarness(testCases []domain.TestCase) string {
	// Construct the test cases dynamically
	var testCasesGo string
	for _, tc := range testCases {
		testCasesGo += fmt.Sprintf(`{"%s", "%s"},`+"\n", strings.ReplaceAll(tc.Input, `"`, `\"`), strings.ReplaceAll(tc.ExpectedOutput, `"`, `\"`))
	}

	return `package main

import (
	"encoding/json"
	"fmt"
)

type TestResult struct {
	Passed         bool   ` + "`json:\"passed\"`" + `
	Input          string ` + "`json:\"input\"`" + `
	ExpectedOutput string ` + "`json:\"expected_output\"`" + `
	ActualOutput   string ` + "`json:\"actual_output\"`" + `
}

func main() {
	testCases := [][2]string{
		` + testCasesGo + `
	}

	results := []TestResult{}
	
	for _, tc := range testCases {
		actual := Solve(tc[0])
		passed := actual == tc[1]
		
		results = append(results, TestResult{
			Passed:         passed,
			Input:          tc[0],
			ExpectedOutput: tc[1],
			ActualOutput:   actual,
		})
	}
	
	out, _ := json.Marshal(results)
	fmt.Println(string(out)) 
}
`
}
