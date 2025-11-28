package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/DucAnhLe1992/ticket-booking-go-app/internal/expiration"
	"github.com/DucAnhLe1992/ticket-booking-go-app/internal/pubsub"
)

func main() {
	// Configuration from environment
	redisAddr := os.Getenv("REDIS_URL")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		log.Fatal("NATS_URL is required")
	}

	// Connect to NATS
	natsClient, err := pubsub.NewNATS(natsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer natsClient.Close()

	log.Println("Connected to NATS")

	// Create expiration queue and worker
	queue := expiration.NewExpirationQueue(redisAddr, natsClient)
	defer queue.Close()

	worker := expiration.NewExpirationWorker(redisAddr, natsClient, 10)

	// Register NATS listeners to schedule expirations
	if err := expiration.RegisterNATSListeners(context.Background(), natsClient, queue); err != nil {
		log.Fatalf("Failed to register NATS listeners: %v", err)
	}

	log.Println("NATS listeners registered")

	// Start worker in a goroutine
	go func() {
		log.Println("Starting expiration worker...")
		if err := worker.Start(); err != nil {
			log.Fatalf("Worker error: %v", err)
		}
	}()

	log.Println("Expiration service is running")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down expiration service...")
	worker.Stop()
	log.Println("Expiration service stopped")
}
