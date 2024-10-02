package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestGracefulShutdown(t *testing.T) {
	// Start the server
	srv := startServer()

	// Simulate a shutdown signal after a delay
	go func() {
		time.Sleep(1 * time.Second) // Wait for 1 second before shutdown
		shutdownServer(srv)
	}()

	// Simulate a client request
	go func() {
		resp, err := http.Get("http://localhost:8080/")
		if err != nil {
			t.Errorf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	}()

	// Allow enough time for the request and shutdown to complete
	time.Sleep(5 * time.Second)
}

// Programmatically sends a shutdown signal
func shutdownServer(srv *http.Server) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// Simulate a shutdown signal
	sig <- syscall.SIGINT

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
