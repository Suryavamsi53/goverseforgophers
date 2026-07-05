package web

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/suryavamsivaggu/goverse/internal/domain"
	"github.com/suryavamsivaggu/goverse/pkg/auth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type contextKey string
const userContextKey contextKey = "user_claims"

type AuthHandler struct {
	UseCase  domain.AuthUseCase
	JWT      *auth.JWTManager
	UserRepo domain.UserRepository
}

func RegisterAuthRoutes(r chi.Router, authUseCase domain.AuthUseCase, jwtManager *auth.JWTManager, userRepo domain.UserRepository) {
	ah := &AuthHandler{
		UseCase:  authUseCase,
		JWT:      jwtManager,
		UserRepo: userRepo,
	}

	r.Get("/login", ah.HandleLoginPage)
	r.Post("/login", ah.HandleLogin)
	r.Get("/register", ah.HandleRegisterPage)
	r.Post("/register", ah.HandleRegister)
	r.Get("/logout", ah.HandleLogout)
	r.Get("/auth/github", ah.HandleGitHubLogin)
	r.Get("/auth/github/callback", ah.HandleGitHubCallback)
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
					MaxAge:   -1,
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

func getGithubOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/auth/github/callback",
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
}

func (ah *AuthHandler) HandleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	url := getGithubOAuthConfig().AuthCodeURL("state-string-goverse", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (ah *AuthHandler) HandleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != "state-string-goverse" {
		http.Error(w, "invalid oauth state", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	token, err := getGithubOAuthConfig().Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := getGithubOAuthConfig().Client(r.Context(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		http.Error(w, "failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var ghUser struct {
		Login string `json:"login"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&ghUser); err != nil {
		http.Error(w, "failed to parse user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if ghUser.Email == "" {
		emailsResp, err := client.Get("https://api.github.com/user/emails")
		if err == nil {
			defer emailsResp.Body.Close()
			var emails []struct {
				Email   string `json:"email"`
				Primary bool   `json:"primary"`
			}
			if err := json.NewDecoder(emailsResp.Body).Decode(&emails); err == nil {
				for _, e := range emails {
					if e.Primary {
						ghUser.Email = e.Email
						break
					}
				}
				if ghUser.Email == "" && len(emails) > 0 {
					ghUser.Email = emails[0].Email
				}
			}
		}
	}

	if ghUser.Email == "" {
		ghUser.Email = ghUser.Login + "@github.com"
	}

	user, err := ah.UserRepo.GetByEmail(r.Context(), ghUser.Email)
	if err != nil {
		_, regErr := ah.UseCase.Register(r.Context(), ghUser.Login, ghUser.Email, "oauth-"+token.AccessToken[:10])
		if regErr != nil {
			http.Error(w, "failed to register user: "+regErr.Error(), http.StatusInternalServerError)
			return
		}
		user, _ = ah.UserRepo.GetByEmail(r.Context(), ghUser.Email)
	}

	jwtToken, err := ah.JWT.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		http.Error(w, "failed to generate token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    jwtToken,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
