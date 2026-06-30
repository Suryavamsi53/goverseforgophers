package web

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/suryavamsivaggu/goverse/internal/domain"
)

type WebHandler struct {
	UserRepo     domain.UserRepository
	CourseRepo   domain.CourseRepository
	ProgressRepo domain.ProgressRepository
}

func RegisterRoutes(r chi.Router, userRepo domain.UserRepository, courseRepo domain.CourseRepository, progressRepo domain.ProgressRepository) {
	h := &WebHandler{
		UserRepo:     userRepo,
		CourseRepo:   courseRepo,
		ProgressRepo: progressRepo,
	}

	r.Get("/", h.HandleLandingPage)
	r.Get("/dashboard", h.HandleDashboard)
	r.Get("/roadmap", h.HandleRoadmap)
	h.RegisterLearnRoutes(r)
	RegisterPracticeRoutes(r)
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
	tmpl := template.New("")
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
