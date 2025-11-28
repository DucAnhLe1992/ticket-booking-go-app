package main

import (
	"log"
	"os"
	"time"

	"github.com/DucAnhLe1992/ticket-booking-go-app/internal/expiration"
)

func main() {
	interval := time.Minute
	log.Println("expiration worker starting")
	for {
		if err := expiration.RunOnce(); err != nil {
			log.Printf("expiration.RunOnce error: %v", err)
		}
		time.Sleep(interval)
		if os.Getenv("ONE_SHOT") == "true" {
			break
		}
	}
}
