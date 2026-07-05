package web

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/suryavamsivaggu/goverse/internal/domain"
	"github.com/suryavamsivaggu/goverse/pkg/auth"
	"github.com/suryavamsivaggu/goverse/pkg/runner"
)

type WebHandler struct {
	UserRepo     domain.UserRepository
	CourseRepo   domain.CourseRepository
	ProgressRepo domain.ProgressRepository
	ProjectRepo   domain.ProjectRepository
	WorkspaceRepo domain.WorkspaceRepository
}

func RegisterRoutes(r chi.Router, userRepo domain.UserRepository, courseRepo domain.CourseRepository, progressRepo domain.ProgressRepository, projectRepo domain.ProjectRepository, workspaceRepo domain.WorkspaceRepository, authUseCase domain.AuthUseCase, jwtManager *auth.JWTManager) {
	h := &WebHandler{
		UserRepo:      userRepo,
		CourseRepo:    courseRepo,
		ProgressRepo:  progressRepo,
		ProjectRepo:   projectRepo,
		WorkspaceRepo: workspaceRepo,
	}

	r.Get("/", h.HandleLandingPage)
	r.Get("/roadmap", h.HandleRoadmap)

	// Auth routes
	RegisterAuthRoutes(r, authUseCase, jwtManager, userRepo)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware(jwtManager))
		r.Get("/dashboard", h.HandleDashboard)
		r.Get("/settings", h.HandleSettingsPage)
		r.Post("/api/v1/settings", h.HandleUpdateSettings)
		r.Get("/leaderboard", h.HandleLeaderboard)
		r.Get("/projects", h.HandleProjects)
		r.Get("/projects/{slug}", h.HandleProjectDetail)
		r.Post("/api/v1/projects/{slug}/submit", h.HandleProjectSubmit)
		h.RegisterLearnRoutes(r)
		h.RegisterPracticeRoutes(r)
	})
}

func getProjectRoot() string {
	if _, err := os.Stat("ui/templates"); err == nil {
		return "."
	}
	if _, err := os.Stat("../../ui/templates"); err == nil {
		return "../.."
	}
	return "."
}

func parseTemplates() *template.Template {
	// Parse base layout and all pages/partials
	root := getProjectRoot()
	
	// Create base template with functions
	tmpl := template.New("base").Funcs(template.FuncMap{
		"add":   func(a, b int) int { return a + b },
		"upper": strings.ToUpper,
	})
	
	tmpl = template.Must(tmpl.ParseGlob(filepath.Join(root, "ui", "templates", "layouts", "*.html")))
	tmpl = template.Must(tmpl.ParseGlob(filepath.Join(root, "ui", "templates", "pages", "*.html")))
	tmpl = template.Must(tmpl.ParseGlob(filepath.Join(root, "ui", "templates", "partials", "*.html")))
	return tmpl
}

func (h *WebHandler) HandleLandingPage(w http.ResponseWriter, r *http.Request) {
	tmpl := parseTemplates()
	err := tmpl.ExecuteTemplate(w, "base", map[string]interface{}{
		"Title": "GoVerse - The Ultimate Golang Learning Platform",
		"Page":  "index",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *WebHandler) HandleRoadmap(w http.ResponseWriter, r *http.Request) {
	tmpl := parseTemplates()
	err := tmpl.ExecuteTemplate(w, "base", map[string]interface{}{
		"Title": "Roadmap - GoVerse",
		"Page":  "roadmap",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *WebHandler) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(userContextKey).(*auth.Claims)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	userID := claims.UserID

	user, err := h.UserRepo.GetByID(r.Context(), userID)
	if err != nil {
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	profile, err := h.UserRepo.GetProfile(r.Context(), userID)
	if err != nil {
		profile = &domain.UserProfile{}
	}

	// Dynamic stats
	progressList, _ := h.ProgressRepo.GetProgress(r.Context(), userID)
	completedLessons := 0
	for _, p := range progressList {
		if p.EntityType == "lesson" && p.Status == "completed" {
			completedLessons++
		}
	}

	// Count total lessons
	totalLessons := 0
	courses, _ := h.CourseRepo.GetAll(r.Context())
	for _, c := range courses {
		lessons, _ := h.CourseRepo.GetLessonsByCourseID(r.Context(), c.ID)
		totalLessons += len(lessons)
	}

	lessonPercent := 0
	if totalLessons > 0 {
		lessonPercent = (completedLessons * 100) / totalLessons
	}

	// Find the next incomplete lesson to display in "Continue Learning"
	var nextLesson *domain.Lesson
	var nextCourse *domain.Course
	
	completedMap := make(map[string]bool)
	for _, p := range progressList {
		if p.EntityType == "lesson" && p.Status == "completed" {
			completedMap[p.EntityID] = true
		}
	}

	foundIncomplete := false
	for _, c := range courses {
		lessons, _ := h.CourseRepo.GetLessonsByCourseID(r.Context(), c.ID)
		for _, l := range lessons {
			if !completedMap[l.ID] {
				lCopy := l
				cCopy := c
				nextLesson = &lCopy
				nextCourse = &cCopy
				foundIncomplete = true
				break
			}
		}
		if foundIncomplete {
			break
		}
	}

	tmpl := parseTemplates()
	err = tmpl.ExecuteTemplate(w, "base", map[string]interface{}{
		"Title":            "Dashboard - GoVerse",
		"Page":             "dashboard",
		"User":             user,
		"Profile":          profile,
		"CompletedLessons": completedLessons,
		"TotalLessons":     totalLessons,
		"LessonPercent":    lessonPercent,
		"NextLesson":       nextLesson,
		"NextCourse":       nextCourse,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *WebHandler) HandleLeaderboard(w http.ResponseWriter, r *http.Request) {
	entries, err := h.UserRepo.GetLeaderboard(r.Context(), 50)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := parseTemplates()
	err = tmpl.ExecuteTemplate(w, "base", map[string]interface{}{
		"Title":   "Leaderboard - GoVerse",
		"Page":    "leaderboard",
		"Entries": entries,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *WebHandler) HandleProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.ProjectRepo.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := parseTemplates()
	err = tmpl.ExecuteTemplate(w, "base", map[string]interface{}{
		"Title":    "Projects - GoVerse",
		"Page":     "projects",
		"Projects": projects,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *WebHandler) HandleProjectDetail(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	project, err := h.ProjectRepo.GetBySlug(r.Context(), slug)
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	tmpl := parseTemplates()
	err = tmpl.ExecuteTemplate(w, "base", map[string]interface{}{
		"Title":   project.Title + " - GoVerse Projects",
		"Page":    "project_detail",
		"Project": project,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type ProjectSubmitRequest struct {
	Code string `json:"code"`
}

type ProjectSubmitResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (h *WebHandler) HandleProjectSubmit(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	project, err := h.ProjectRepo.GetBySlug(r.Context(), slug)
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	var req ProjectSubmitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := runner.EvaluateProject(r.Context(), req.Code, project)
	if err != nil {
		http.Error(w, "Failed to evaluate project", http.StatusInternalServerError)
		return
	}

	resp := ProjectSubmitResponse{
		Success: result.Success,
		Message: result.SystemError,
	}

	if result.Success {
		// Mark project as complete
		user := r.Context().Value("user").(*domain.User)
		_ = h.ProgressRepo.MarkCompleted(r.Context(), user.ID, "project", project.ID)
		resp.Message = "Congratulations! Your project passed all tests."
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *WebHandler) HandleSettingsPage(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(userContextKey).(*auth.Claims)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := h.UserRepo.GetByID(r.Context(), claims.UserID)
	if err != nil {
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	profile, err := h.UserRepo.GetProfile(r.Context(), user.ID)
	if err != nil {
		profile = &domain.UserProfile{}
	}

	tmpl := parseTemplates()
	err = tmpl.ExecuteTemplate(w, "base", map[string]interface{}{
		"Title":      "Settings - GoVerse",
		"Page":       "settings",
		"User":       user,
		"Profile":    profile,
		"IsLoggedIn": true,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *WebHandler) HandleUpdateSettings(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(userContextKey).(*auth.Claims)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	
	user, err := h.UserRepo.GetByID(r.Context(), claims.UserID)
	if err != nil {
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	
	var req domain.UserProfile
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	req.UserID = user.ID
	
	// We check if profile exists, if not create, else update
	existingProfile, err := h.UserRepo.GetProfile(r.Context(), user.ID)
	if err != nil {
		err = h.UserRepo.CreateProfile(r.Context(), &req)
	} else {
		// Preserve stats
		req.DailyStreak = existingProfile.DailyStreak
		req.TotalScore = existingProfile.TotalScore
		err = h.UserRepo.UpdateProfile(r.Context(), &req)
	}

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to save profile"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
