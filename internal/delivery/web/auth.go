package web

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/suryavamsivaggu/goverse/internal/domain"
	"github.com/suryavamsivaggu/goverse/pkg/auth"
)

type contextKey string
const userContextKey contextKey = "user_claims"

type AuthHandler struct {
	UseCase domain.AuthUseCase
	JWT     *auth.JWTManager
}

func RegisterAuthRoutes(r chi.Router, authUseCase domain.AuthUseCase, jwtManager *auth.JWTManager) {
	ah := &AuthHandler{
		UseCase: authUseCase,
		JWT:     jwtManager,
	}

	r.Get("/login", ah.HandleLoginPage)
	r.Post("/login", ah.HandleLogin)
	r.Get("/register", ah.HandleRegisterPage)
	r.Post("/register", ah.HandleRegister)
	r.Get("/logout", ah.HandleLogout)
}

func AuthMiddleware(jwtManager *auth.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			claims, err := jwtManager.VerifyToken(cookie.Value)
			if err != nil {
				http.SetCookie(w, &http.Cookie{
					Name:     "token",
					Value:    "",
					Path:     "/",
					Expires:  time.Now().Add(-1 * time.Hour),
					HttpOnly: true,
				})
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (ah *AuthHandler) HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl := parseTemplates()
	err := tmpl.ExecuteTemplate(w, "base", map[string]interface{}{
		"Title": "Login - GoVerse",
		"Page":  "login",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ah *AuthHandler) HandleRegisterPage(w http.ResponseWriter, r *http.Request) {
	tmpl := parseTemplates()
	err := tmpl.ExecuteTemplate(w, "base", map[string]interface{}{
		"Title": "Register - GoVerse",
		"Page":  "register",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ah *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	token, err := ah.UseCase.Login(r.Context(), email, password)
	if err != nil {
		tmpl := parseTemplates()
		tmpl.ExecuteTemplate(w, "base", map[string]interface{}{
			"Title": "Login - GoVerse",
			"Page":  "login",
			"Error": err.Error(),
		})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (ah *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	_, err := ah.UseCase.Register(r.Context(), username, email, password)
	if err != nil {
		tmpl := parseTemplates()
		tmpl.ExecuteTemplate(w, "base", map[string]interface{}{
			"Title": "Register - GoVerse",
			"Page":  "register",
			"Error": err.Error(),
		})
		return
	}

	// Auto-login
	token, _ := ah.UseCase.Login(r.Context(), email, password)
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (ah *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
