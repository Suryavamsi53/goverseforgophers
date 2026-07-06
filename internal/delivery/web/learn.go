package web

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/suryavamsivaggu/goverse/internal/domain"
	"github.com/suryavamsivaggu/goverse/pkg/auth"
	"github.com/suryavamsivaggu/goverse/pkg/markdown"
)

var mdRenderer = markdown.NewRenderer()

func (h *WebHandler) RegisterLearnRoutes(r chi.Router) {
	r.Get("/learn", h.HandleLearnIndex)
	r.Get("/learn/{course}/{lesson}", h.HandleLesson)
	r.Post("/api/progress/lesson/{lesson}", h.HandleMarkProgress)
}

func (h *WebHandler) HandleLearnIndex(w http.ResponseWriter, r *http.Request) {
	courses, err := h.CourseRepo.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	groupedCourses := make(map[string][]map[string]interface{})
	groupOrder := []string{"beginner", "intermediate", "advanced", "expert"}

	for _, c := range courses {
		lessons, err := h.CourseRepo.GetLessonsByCourseID(r.Context(), c.ID)
		startSlug := "introduction" // Fallback
		if err == nil && len(lessons) > 0 {
			startSlug = lessons[0].Slug
		}
		
		diff := c.Difficulty
		if diff == "" {
			diff = "intermediate"
		}

		groupedCourses[diff] = append(groupedCourses[diff], map[string]interface{}{
			"Course":      c,
			"StartSlug":   startSlug,
			"LessonCount": len(lessons),
		})
	}

	tmpl := parseTemplates()
	data := h.getBaseTemplateData(r, "Learn Golang - GoVerse", "learn_index")
	data["GroupedCourses"] = groupedCourses
	data["GroupOrder"] = groupOrder
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *WebHandler) HandleLesson(w http.ResponseWriter, r *http.Request) {
	courseSlug := chi.URLParam(r, "course")
	lessonSlug := chi.URLParam(r, "lesson")
	
	if courseSlug == "" || lessonSlug == "" {
		http.NotFound(w, r)
		return
	}

	// Read markdown file safely
	cleanCourse := filepath.Base(courseSlug)
	cleanLesson := filepath.Base(lessonSlug)
	contentPath := filepath.Join("content", "lessons", cleanCourse, cleanLesson+".md")
	
	mdBytes, err := os.ReadFile(contentPath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Render Markdown to HTML
	htmlContent, err := mdRenderer.Render(mdBytes)
	if err != nil {
		http.Error(w, "Error rendering content", http.StatusInternalServerError)
		return
	}
	
	// Extract Table of Contents
	toc, err := mdRenderer.ExtractTOC(mdBytes)
	if err != nil {
		http.Error(w, "Error parsing TOC", http.StatusInternalServerError)
		return
	}

	// Fetch lesson details from DB to check completion status
	lesson, err := h.CourseRepo.GetLessonBySlug(r.Context(), cleanLesson)
	var isCompleted bool
	var userID string
	if claims, ok := r.Context().Value(userContextKey).(*auth.Claims); ok {
		userID = claims.UserID
	}
	
	if err == nil && lesson != nil {
		progressList, progressErr := h.ProgressRepo.GetProgress(r.Context(), userID)
		if progressErr == nil {
			for _, p := range progressList {
				if p.EntityType == "lesson" && p.EntityID == lesson.ID && p.Status == "completed" {
					isCompleted = true
					break
				}
			}
		}
	}

	// Get course info and lessons list for sidebar
	course, err := h.CourseRepo.GetBySlug(r.Context(), cleanCourse)
	var lessons []domain.Lesson
	var courseTitle string
	if err == nil && course != nil {
		courseTitle = course.Title
		lessons, _ = h.CourseRepo.GetLessonsByCourseID(r.Context(), course.ID)
	} else {
		// Fallback/Default
		courseTitle = strings.Title(strings.ReplaceAll(cleanCourse, "-", " "))
	}

	// Fetch user progress for sidebar item completion statuses
	completedLessons := make(map[string]bool)
	progressList, progressErr := h.ProgressRepo.GetProgress(r.Context(), userID)
	if progressErr == nil {
		for _, p := range progressList {
			if p.EntityType == "lesson" && p.Status == "completed" {
				completedLessons[p.EntityID] = true
			}
		}
	}

	// Prepare sidebar lessons and find prev/next
	var prevLesson, nextLesson *domain.Lesson
	var sidebarLessons []map[string]interface{}
	
	for i, l := range lessons {
		isActive := l.Slug == cleanLesson
		isComp := completedLessons[l.ID]
		
		sidebarLessons = append(sidebarLessons, map[string]interface{}{
			"Slug":        l.Slug,
			"Title":       l.Title,
			"IsActive":    isActive,
			"IsCompleted": isComp,
		})
		
		if isActive {
			if i > 0 {
				prevLesson = &lessons[i-1]
			}
			if i < len(lessons)-1 {
				nextLesson = &lessons[i+1]
			}
		}
	}

	tmpl := parseTemplates()
	data := h.getBaseTemplateData(r, strings.Title(strings.ReplaceAll(cleanLesson, "-", " ")) + " - GoVerse", "lesson")
	data["Content"] = template.HTML(htmlContent) // Safe because we trust our own markdown
	data["TOC"] = toc
	data["Course"] = cleanCourse
	data["CourseTitle"] = courseTitle
	data["Slug"] = cleanLesson
	data["IsCompleted"] = isCompleted
	data["Lessons"] = sidebarLessons
	data["PrevLesson"] = prevLesson
	data["NextLesson"] = nextLesson
	
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *WebHandler) HandleMarkProgress(w http.ResponseWriter, r *http.Request) {
	lessonSlug := chi.URLParam(r, "lesson")
	cleanLesson := filepath.Base(lessonSlug)
	
	var userID string
	if claims, ok := r.Context().Value(userContextKey).(*auth.Claims); ok {
		userID = claims.UserID
	}
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	// Find lesson ID by slug
	lesson, err := h.CourseRepo.GetLessonBySlug(r.Context(), cleanLesson)
	if err != nil {
		http.Error(w, "Lesson not found", http.StatusNotFound)
		return
	}
	
	err = h.ProgressRepo.MarkCompleted(r.Context(), userID, "lesson", lesson.ID)
	if err != nil {
		http.Error(w, "Failed to save progress", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`
		<div class="w-full flex items-center justify-center px-4 py-2 bg-go-cyan/10 text-go-cyan rounded-lg font-medium border border-go-cyan/30">
			<svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
			Completed
		</div>
	`))
}
