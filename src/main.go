package main

import (
	"log"
	"net/http"

	"rpi-workload/api"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Create a new server implementation
	server := api.NewServer()

	// Create the strict handler
	strictHandler := api.NewStrictHandler(server, nil)

	// Register the handler with the router
	api.HandlerFromMux(strictHandler, r)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
