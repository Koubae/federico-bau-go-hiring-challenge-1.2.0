package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/mytheresa/go-hiring-challenge/app/container"
	"github.com/mytheresa/go-hiring-challenge/app/middlewares"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	container.CreateDIContainer()
	defer container.ShutDown()

	// Initialize handlers
	cat := catalog.NewCatalogHandler(container.Container.ProductRepository)

	// Set up routing
	mux := http.NewServeMux()
	mux.HandleFunc("GET /catalog", cat.HandleGet)

	var handler http.Handler = mux
	handler = middlewares.RecoverPanic(middlewares.LogAccessMiddleware(mux))

	// Set up the HTTP server
	// TODO: better way to handle how we load env variables!"
	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%s", os.Getenv("HTTP_PORT")),
		Handler: handler,
	}

	// Start the server
	go func() {
		log.Printf("Starting server on http://%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %s", err)
		}

		log.Println("Server stopped gracefully")
	}()

	<-ctx.Done()
	log.Println("Shutting down server...")
	srv.Shutdown(ctx)
	stop()
}
