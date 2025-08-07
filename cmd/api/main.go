package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"library-management/internal/config"
	"library-management/internal/database"
	"library-management/internal/handler"
	"library-management/internal/repository/postgres"
	"library-management/internal/service"
	"library-management/pkg/logger"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize logger
	log := logger.New()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration", "error", err)
	}

	// Connect to database
	log.Info("Connecting to database...")
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatal("Database ping failed", "error", err)
	}
	log.Info("Database connection established")

	// Initialize database schema
	log.Info("Initializing database...")
	if err := database.InitializeDatabase(db); err != nil {
		log.Fatal("Failed to initialize database", "error", err)
	}
	log.Info("Database initialization completed")

	// Initialize layers
	bookRepo := postgres.NewBookRepository(db)
	bookService := service.NewBookService(bookRepo)
	handlers := handler.NewHandlers(bookService, log)

	// Setup router
	router := mux.NewRouter()
	handler.SetupRoutes(router, handlers)

	// Configure server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      router,
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	// Start server in goroutine
	go func() {
		log.Info("Starting server", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start", "error", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", "error", err)
	}

	log.Info("Server exited")
}
