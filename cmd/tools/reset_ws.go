package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/goverse?sslmode=disable"
	}
	
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	_, err = pool.Exec(context.Background(), "DELETE FROM workspaces WHERE ref_id = 'log-analyzer'")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deleted workspaces for log-analyzer")
}
