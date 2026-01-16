package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"rpi-workload/api"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://user:password@localhost:5432/workload_db"
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	// Create a new server implementation
	server := api.NewServer(pool)

	// Initialize the schema
	if err := server.InitSchema(context.Background()); err != nil {
		log.Printf("Failed to initialize schema: %v", err)
		// We might not want to fatal here if the DB is just not reachable yet and we want to retry or something,
		// but for now, logging it is enough.
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Create the strict handler
	strictHandler := api.NewStrictHandler(server, nil)

	// Register the handler with the router
	api.HandlerFromMux(strictHandler, r)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
