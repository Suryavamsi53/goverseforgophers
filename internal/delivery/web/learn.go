package web

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/suryavamsivaggu/goverse/pkg/markdown"
)

var mdRenderer = markdown.NewRenderer()

func RegisterLearnRoutes(r chi.Router) {
	r.Get("/learn", HandleLearnIndex)
	r.Get("/learn/{slug}", HandleLesson)
	r.Post("/api/progress/lesson/{slug}", HandleMarkProgress)
}

func HandleLearnIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := parseTemplates()
	err := tmpl.ExecuteTemplate(w, "base", map[string]interface{}{
		"Title": "Learn Golang - GoVerse",
		"Page":  "learn_index",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandleLesson(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.NotFound(w, r)
		return
	}

	// Read markdown file safely (prevent path traversal in real app)
	// For demo, we just read from content/lessons/{slug}.md
	cleanSlug := filepath.Base(slug)
	contentPath := filepath.Join("content", "lessons", cleanSlug+".md")
	
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

	tmpl := parseTemplates()
	err = tmpl.ExecuteTemplate(w, "base", map[string]interface{}{
		"Title":   strings.Title(strings.ReplaceAll(cleanSlug, "-", " ")) + " - GoVerse",
		"Page":    "lesson",
		"Content": template.HTML(htmlContent), // Safe because we trust our own markdown
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandleMarkProgress(w http.ResponseWriter, r *http.Request) {
	// For demo purposes, we'll just return the updated HTML snippet directly via HTMX.
	// In a real app, you would extract the user from the JWT and update the DB using domain.ProgressRepository.
	
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`
		<div class="w-full flex items-center justify-center px-4 py-2 bg-go-cyan/10 text-go-cyan rounded-lg font-medium border border-go-cyan/30">
			<svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
			Completed
		</div>
	`))
}
