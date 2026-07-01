package web

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/suryavamsivaggu/goverse/internal/domain"
	"github.com/suryavamsivaggu/goverse/pkg/auth"
	"github.com/suryavamsivaggu/goverse/pkg/runner"
)

func (h *WebHandler) RegisterPracticeRoutes(r chi.Router) {
	r.Get("/practice", h.HandlePracticePage)
	r.Post("/api/v1/execute", h.HandleExecuteCode)
	r.Post("/api/v1/terminal", h.HandleTerminalCommand)
	r.Get("/api/v1/ws/terminal", h.HandleWSTerminal)
	r.Post("/api/v1/evaluate", h.HandleEvaluateCode)
	r.Post("/api/v1/format", h.HandleFormatCode)
	r.Get("/api/v1/workspace", h.HandleGetWorkspace)
	r.Post("/api/v1/workspace", h.HandleSaveWorkspace)
}

func (h *WebHandler) HandlePracticePage(w http.ResponseWriter, r *http.Request) {
	wsType := r.URL.Query().Get("type")
	if wsType == "" { wsType = "practice" }
	refID := r.URL.Query().Get("ref_id")
	if refID == "" { refID = "default" }

	tmpl := parseTemplates()
	err := tmpl.ExecuteTemplate(w, "base", map[string]interface{}{
		"Title":      "Go Sandbox Playground",
		"Page":       "practice",
		"IsEmbedded": r.URL.Query().Get("embed") == "true",
		"WsType":     wsType,
		"RefID":      refID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type ExecuteRequest struct {
	Files map[string]string `json:"files"`
}

func (h *WebHandler) HandleExecuteCode(w http.ResponseWriter, r *http.Request) {
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

func (h *WebHandler) HandleTerminalCommand(w http.ResponseWriter, r *http.Request) {
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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for the sandbox
	},
}

func (h *WebHandler) HandleWSTerminal(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	runner.HandleWSTerminalSession(conn)
}

type EvaluateRequest struct {
	ProblemSlug string `json:"problem_slug"`
	Code        string `json:"code"`
}

func (h *WebHandler) HandleEvaluateCode(w http.ResponseWriter, r *http.Request) {
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

type FormatRequest struct {
	Code string `json:"code"`
}

type FormatResponse struct {
	FormattedCode string `json:"formatted_code"`
	Error         string `json:"error"`
}

func (h *WebHandler) HandleFormatCode(w http.ResponseWriter, r *http.Request) {
	var req FormatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := runner.FormatCode(r.Context(), req.Code)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(FormatResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(FormatResponse{FormattedCode: result})
}

func (h *WebHandler) HandleGetWorkspace(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(userContextKey).(*auth.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	wsType := r.URL.Query().Get("type")
	refID := r.URL.Query().Get("ref_id")
	if wsType == "" {
		wsType = "practice"
	}
	if refID == "" {
		refID = "default"
	}

	ws, err := h.WorkspaceRepo.Get(r.Context(), claims.UserID, wsType, refID)
	if err != nil {
		// Return 404 so frontend knows to use default files
		http.Error(w, "Workspace not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ws)
}

func (h *WebHandler) HandleSaveWorkspace(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(userContextKey).(*auth.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req domain.Workspace
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	req.UserID = claims.UserID
	if req.Type == "" {
		req.Type = "practice"
	}
	if req.RefID == "" {
		req.RefID = "default"
	}

	err := h.WorkspaceRepo.Save(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
