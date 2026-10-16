package main

import (
	"context"
	"log"
	"market/internal/auth"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:1234@localhost:5432/market?sslmode=disable"
	}
	pool, err := pgxpool.New(context.Background(), databaseURL)

	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	mux := http.NewServeMux()

	auth := auth.NewAuthModule(pool)

	auth.RegisterRoutes(mux)

	go func() { log.Println("Started auth-service") }()

	err = http.ListenAndServe(":8080", mux)

	if err != nil {
		log.Fatal(err)
	}
}
