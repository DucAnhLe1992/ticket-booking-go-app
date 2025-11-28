package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	port := ":3003"
	if p := os.Getenv("TICKETS_PORT"); p != "" {
		port = ":" + p
	}
	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("tickets ok"))
	})
	log.Printf("tickets service listening on %s\n", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal(err)
	}
}
