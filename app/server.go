package app

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
	"github.com/mytheresa/go-hiring-challenge/app/category"
	"github.com/mytheresa/go-hiring-challenge/app/container"
	"github.com/mytheresa/go-hiring-challenge/app/middlewares"
)

func RunServer() {
	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	container.CreateDIContainer()
	defer container.ShutDown()

	// Set up routing
	mux := http.NewServeMux()

	categoryHandler := category.NewCategoryHandler(container.Container.CategoryRepository)
	mux.HandleFunc("GET /categories", categoryHandler.ListCategories)
	mux.HandleFunc("POST /categories", categoryHandler.CreateCategory)

	catalogHandler := catalog.NewCatalogHandler(container.Container.ProductRepository)
	mux.HandleFunc("GET /catalog", catalogHandler.ListCatalog)
	mux.HandleFunc("GET /catalog/{code}", catalogHandler.GetProductDetails)

	var handler http.Handler = mux
	handler = middlewares.RecoverPanic(middlewares.LogAccessMiddleware(mux))

	// Set up the HTTP server
	// TODO: better way to handle how we load env variables!"
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", os.Getenv("HTTP_HOST"), os.Getenv("HTTP_PORT")),
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
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Error while shutting down server: %v", err)
		return
	}
	stop()
}
