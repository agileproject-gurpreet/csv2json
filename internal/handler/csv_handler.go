package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/agileproject-gurpreet/csv2json/internal/service"
)

type CSVHandler struct {
	service *service.ConversionService
	logger  *log.Logger
}

func NewCSVHandler(service *service.ConversionService, logger *log.Logger) *CSVHandler {
	return &CSVHandler{
		service: service,
		logger:  logger,
	}
}

// UploadCSV handles CSV file upload
func (h *CSVHandler) UploadCSV(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.Printf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.logger.Println("Received CSV upload request")

	// Parse multipart form (32MB max)
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		h.logger.Printf("Failed to parse form: %v", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get the file from form
	file, header, err := r.FormFile("file")
	if err != nil {
		h.logger.Printf("Failed to get file from form: %v", err)
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	h.logger.Printf("Processing file: %s (size: %d bytes)", header.Filename, header.Size)

	// Process the CSV file
	jsonData, err := h.service.ProcessCSVReaderWithFilename(file, header.Filename)
	if err != nil {
		h.logger.Printf("Failed to process CSV file '%s': %v", header.Filename, err)
		http.Error(w, fmt.Sprintf("Failed to process CSV: %v", err), http.StatusInternalServerError)
		return
	}

	h.logger.Printf("Successfully processed CSV file: %s, converted %d bytes to JSON", header.Filename, len(jsonData))

	// Send JSON data response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// Health check endpoint
func (h *CSVHandler) Health(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("Health check requested")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}

// GetAllData retrieves all CSV data from the database
func (h *CSVHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.logger.Printf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.logger.Println("Received get all data request")

	data, err := h.service.GetAllData()
	if err != nil {
		h.logger.Printf("Failed to retrieve data: %v", err)
		http.Error(w, fmt.Sprintf("Failed to retrieve data: %v", err), http.StatusInternalServerError)
		return
	}

	h.logger.Printf("Successfully retrieved %d records", len(data))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

// GetDataByID retrieves a specific CSV data record by ID
func (h *CSVHandler) GetDataByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.logger.Printf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse ID from query parameter
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	h.logger.Printf("Received get data by ID request: %d", id)

	data, err := h.service.GetDataByID(id)
	if err != nil {
		h.logger.Printf("Failed to retrieve data for ID %d: %v", id, err)
		http.Error(w, fmt.Sprintf("Failed to retrieve data: %v", err), http.StatusNotFound)
		return
	}

	h.logger.Printf("Successfully retrieved record ID: %d", id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
