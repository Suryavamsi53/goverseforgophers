package web

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/suryavamsivaggu/goverse/internal/domain"
)

type WebHandler struct {
	UserRepo domain.UserRepository
}

func RegisterRoutes(r chi.Router, userRepo domain.UserRepository) {
	h := &WebHandler{
		UserRepo: userRepo,
	}

	r.Get("/", h.HandleLandingPage)
	r.Get("/dashboard", h.HandleDashboard)
	RegisterLearnRoutes(r)
	RegisterPracticeRoutes(r)
}

func parseTemplates() *template.Template {
	// Parse base layout and all pages/partials
	tmpl := template.New("")
	tmpl, _ = tmpl.ParseGlob(filepath.Join("ui", "templates", "layouts", "*.html"))
	tmpl, _ = tmpl.ParseGlob(filepath.Join("ui", "templates", "pages", "*.html"))
	tmpl, _ = tmpl.ParseGlob(filepath.Join("ui", "templates", "partials", "*.html"))
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

func (h *WebHandler) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	// For now, hardcode the authenticated user ID (the one we seeded in DB)
	userID := "11111111-1111-1111-1111-111111111111"

	user, err := h.UserRepo.GetByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	profile, err := h.UserRepo.GetProfile(r.Context(), userID)
	if err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	tmpl := parseTemplates()
	err = tmpl.ExecuteTemplate(w, "base", map[string]interface{}{
		"Title":   "Dashboard - GoVerse",
		"Page":    "dashboard",
		"User":    user,
		"Profile": profile,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
