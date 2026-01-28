package main

import (
	"log"
	"net/http"
	"os"

	"github.com/lavishag4193/csv2jsonx/internal/database"
	"github.com/lavishag4193/csv2jsonx/internal/handler"
	"github.com/lavishag4193/csv2jsonx/internal/service"
)

func main() {
	// Setup logger
	logger := log.New(os.Stdout, "[CSV2JSON-API] ", log.LstdFlags|log.Lshortfile)
	logger.Println("Starting CSV2JSON API server...")

	// Initialize PostgreSQL database (optional)
	var db *database.PostgresDB
	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "csv2json"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	db, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		logger.Printf("Warning: Failed to connect to database: %v", err)
		logger.Println("Running without database persistence - CSV conversion will still work!")
		db = nil
	} else {
		defer db.Close()
		logger.Println("Successfully connected to PostgreSQL database")

		// Initialize database schema
		if err := db.InitSchema(); err != nil {
			logger.Printf("Warning: Failed to initialize database schema: %v", err)
			db = nil
		} else {
			logger.Println("Database schema initialized")
		}
	}

	// Initialize service and handler
	svc := service.NewConversionService(db)
	csvHandler := handler.NewCSVHandler(svc, logger)

	// Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("/api/upload", csvHandler.UploadCSV)
	mux.HandleFunc("/api/data", csvHandler.GetAllData)
	mux.HandleFunc("/api/data/id", csvHandler.GetDataByID)
	mux.HandleFunc("/api/health", csvHandler.Health)

	// Start server
	port := getEnv("PORT", "8080")
	addr := ":" + port

	logger.Printf("Server starting on port %s", port)
	logger.Println("Available endpoints:")
	logger.Println("  POST /api/upload     - Upload CSV file")
	logger.Println("  GET  /api/data       - Get all stored CSV data")
	logger.Println("  GET  /api/data/id    - Get CSV data by ID (requires ?id=<id>)")
	logger.Println("  GET  /api/health     - Health check")

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
