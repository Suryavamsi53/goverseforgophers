package web

import (
	"encoding/json"
	"html/template"
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
	r.Post("/api/v1/user/ide-settings", h.HandleUpdateIDESettings)
}

func (h *WebHandler) HandlePracticePage(w http.ResponseWriter, r *http.Request) {
	wsType := r.URL.Query().Get("type")
	if wsType == "" { wsType = "practice" }
	refID := r.URL.Query().Get("ref_id")
	if refID == "" { refID = "default" }

	var settings *domain.UserSettings
	claims, ok := r.Context().Value(userContextKey).(*auth.Claims)
	if ok {
		settings, _ = h.UserRepo.GetSettings(r.Context(), claims.UserID)
	}
	if settings == nil {
		settings = &domain.UserSettings{
			EditorSettings: map[string]interface{}{"fontSize": 14, "minimap": false, "wordWrap": "off", "tabSize": 4, "theme": "goverseDark"},
			Extensions: map[string]bool{"gemini": false, "dracula": false, "vim": false, "go-snippets": false, "go-linter": false, "html-preview": false, "html-snippets": false, "rest-client": false, "docker": false, "sql-tools": false},
		}
	}
	settingsJSON, _ := json.Marshal(settings)

	var defaultFilesJSON string = "{}"
	if wsType == "project" && refID != "" && h.ProjectRepo != nil {
		proj, err := h.ProjectRepo.GetBySlug(r.Context(), refID)
		if err == nil && proj != nil {
			defaultFiles := map[string]string{
				"main.go": proj.StarterCode,
			}
			b, _ := json.Marshal(defaultFiles)
			defaultFilesJSON = string(b)
		}
	}

	tmpl := parseTemplates()
	data := h.getBaseTemplateData(r, "Go Sandbox Playground", "practice")
	data["IsEmbedded"] = r.URL.Query().Get("embed") == "true"
	data["WsType"] = wsType
	data["RefID"] = refID
	data["SettingsJSON"] = template.HTML(settingsJSON)
	data["DefaultFilesJSON"] = template.HTML(defaultFilesJSON)
	
	err := tmpl.ExecuteTemplate(w, "base", data)
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
		if wsType == "project" && h.ProjectRepo != nil {
			proj, pErr := h.ProjectRepo.GetBySlug(r.Context(), refID)
			if pErr == nil && proj != nil {
				ws = &domain.Workspace{
					UserID: claims.UserID,
					Type:   wsType,
					RefID:  refID,
					Files: map[string]string{
						"main.go": proj.StarterCode,
					},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(ws)
				return
			}
		}
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

func (h *WebHandler) HandleUpdateIDESettings(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(userContextKey).(*auth.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req domain.UserSettings
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	req.UserID = claims.UserID

	if err := h.UserRepo.UpdateSettings(r.Context(), &req); err != nil {
		http.Error(w, "Failed to save settings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
