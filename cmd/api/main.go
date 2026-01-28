package main

import (
	"log"
	"net/http"
	"os"

	"github.com/lavishag4193/csv2jsonx/internal/handler"
	"github.com/lavishag4193/csv2jsonx/internal/service"
)

func main() {
	// Setup logger
	logger := log.New(os.Stdout, "[CSV2JSON-API] ", log.LstdFlags|log.Lshortfile)
	logger.Println("Starting CSV2JSON API server...")

	// Initialize service and handler (no database required)
	svc := service.NewConversionService()
	csvHandler := handler.NewCSVHandler(svc, logger)

	// Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("/api/upload", csvHandler.UploadCSV)
	mux.HandleFunc("/api/health", csvHandler.Health)

	// Start server
	port := getEnv("PORT", "8080")
	addr := ":" + port

	logger.Printf("Server starting on port %s", port)
	logger.Println("Available endpoints:")
	logger.Println("  POST /api/upload - Upload CSV file")
	logger.Println("  GET  /api/health - Health check")

	if err := http.ListenAndServe(addr, mux); err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
