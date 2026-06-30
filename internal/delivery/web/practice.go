package web

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/suryavamsivaggu/goverse/internal/domain"
	"github.com/suryavamsivaggu/goverse/pkg/runner"
)

func RegisterPracticeRoutes(r chi.Router) {
	r.Get("/practice", HandlePracticePage)
	r.Post("/api/v1/execute", HandleExecuteCode)
	r.Post("/api/v1/terminal", HandleTerminalCommand)
	r.Post("/api/v1/evaluate", HandleEvaluateCode)
}

func HandlePracticePage(w http.ResponseWriter, r *http.Request) {
	tmpl := parseTemplates()
	err := tmpl.ExecuteTemplate(w, "base", map[string]interface{}{
		"Title": "Go Sandbox Playground",
		"Page":  "practice",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type ExecuteRequest struct {
	Files map[string]string `json:"files"`
}

func HandleExecuteCode(w http.ResponseWriter, r *http.Request) {
	var req ExecuteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Files) == 0 {
		http.Error(w, "Files cannot be empty", http.StatusBadRequest)
		return
	}

	result, err := runner.ExecuteCode(r.Context(), req.Files)
	if err != nil {
		// Log the error internally, return generic error to client
		http.Error(w, "Failed to execute code", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

type TerminalRequest struct {
	Command string            `json:"command"`
	Files   map[string]string `json:"files"`
}

func HandleTerminalCommand(w http.ResponseWriter, r *http.Request) {
	var req TerminalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Command == "" {
		http.Error(w, "Command cannot be empty", http.StatusBadRequest)
		return
	}

	result, err := runner.RunCommand(r.Context(), req.Command, req.Files)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

type EvaluateRequest struct {
	ProblemSlug string `json:"problem_slug"`
	Code        string `json:"code"`
}

func HandleEvaluateCode(w http.ResponseWriter, r *http.Request) {
	var req EvaluateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Code == "" || req.ProblemSlug == "" {
		http.Error(w, "Code and problem slug cannot be empty", http.StatusBadRequest)
		return
	}

	// In a real app, you would fetch this from postgres using domain.ProblemRepository
	// For this demo, we'll mock a simple problem
	mockProblem := &domain.PracticeProblem{
		Slug:  "two-sum",
		Title: "Two Sum",
		TestCases: []domain.TestCase{
			{Input: "[2,7,11,15], target=9", ExpectedOutput: "[0,1]"},
			{Input: "[3,2,4], target=6", ExpectedOutput: "[1,2]"},
		},
	}

	result, err := runner.EvaluateProblem(r.Context(), req.Code, mockProblem)
	if err != nil {
		http.Error(w, "Failed to evaluate code", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
