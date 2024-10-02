package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// HTTP server that supports graceful shutdown
func startServer() *http.Server {
	srv := &http.Server{Addr: ":8080"}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // Simulate work for 2 seconds
		fmt.Fprintf(w, "Hello, world!")
	})

	// Start server in a separate goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	return srv
}

// Handles graceful shutdown logic
func gracefulShutdown(srv *http.Server) {
	// Create a channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a shutdown signal
	<-quit
	log.Println("Shutdown signal received, shutting down gracefully...")

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}

func main() {
	// Start the HTTP server
	srv := startServer()

	// Handle graceful shutdown
	gracefulShutdown(srv)
}
