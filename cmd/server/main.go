package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/suryavamsivaggu/goverse/internal/delivery/web"
	"github.com/suryavamsivaggu/goverse/internal/repository"
	"github.com/suryavamsivaggu/goverse/internal/usecase"
	"github.com/suryavamsivaggu/goverse/pkg/auth"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Static files
	fileServer := http.FileServer(http.Dir("./ui/assets"))
	r.Handle("/assets/*", http.StripPrefix("/assets/", fileServer))

	// Connect to Database
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://goverse_admin:goverse_password@localhost:5433/goverse_db?sslmode=disable"
	}
	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbPool.Close()

	userRepo := repository.NewPostgresUserRepository(dbPool)
	courseRepo := repository.NewPostgresCourseRepository(dbPool)
	progressRepo := repository.NewPostgresProgressRepository(dbPool)
	projectRepo := repository.NewPostgresProjectRepository(dbPool)
	workspaceRepo := repository.NewPostgresWorkspaceRepository(dbPool)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "super-secret-key-change-in-prod"
	}
	jwtManager := auth.NewJWTManager(jwtSecret, 24*time.Hour)
	authUseCase := usecase.NewAuthUseCase(userRepo, jwtManager)

	// Register web routes
	web.RegisterRoutes(r, userRepo, courseRepo, progressRepo, projectRepo, workspaceRepo, authUseCase, jwtManager)

	// Server setup
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Starting GoVerse server on port %s...", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting gracefully")
}
