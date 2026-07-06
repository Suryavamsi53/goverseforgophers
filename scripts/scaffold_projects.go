package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/suryavamsivaggu/goverse/internal/repository"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mockRepo := repository.NewMockProjectRepository()
	projects, _ := mockRepo.GetAll(context.Background())

	for _, p := range projects {
		tagsJSON, _ := json.Marshal(p.Tags)
		reqsJSON, _ := json.Marshal(p.Requirements)
		tipsJSON, _ := json.Marshal(p.Tips)

		_, err := db.Exec(`
			INSERT INTO projects (
				slug, title, description, scenario, difficulty, icon, color, 
				tags, requirements, tips, starter_code, test_file
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
			) ON CONFLICT (slug) DO UPDATE SET
				title = EXCLUDED.title,
				description = EXCLUDED.description,
				scenario = EXCLUDED.scenario,
				difficulty = EXCLUDED.difficulty,
				icon = EXCLUDED.icon,
				color = EXCLUDED.color,
				tags = EXCLUDED.tags,
				requirements = EXCLUDED.requirements,
				tips = EXCLUDED.tips,
				starter_code = EXCLUDED.starter_code,
				test_file = EXCLUDED.test_file
		`, p.Slug, p.Title, p.Description, p.Scenario, p.Difficulty, p.Icon, p.Color,
			tagsJSON, reqsJSON, tipsJSON, p.StarterCode, p.TestFile)

		if err != nil {
			log.Printf("Failed to insert project %s: %v\n", p.Slug, err)
		} else {
			fmt.Printf("Successfully scaffolded project: %s\n", p.Title)
		}
	}

	fmt.Println("Projects scaffolding completed!")
}
