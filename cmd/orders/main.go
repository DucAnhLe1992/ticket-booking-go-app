package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	port := ":3001"
	if p := os.Getenv("ORDERS_PORT"); p != "" {
		port = ":" + p
	}
	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("orders ok"))
	})
	log.Printf("orders service listening on %s\n", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal(err)
	}
}
