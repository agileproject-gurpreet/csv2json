package service

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/lavishag4193/csv2jsonx/internal/database"
	"github.com/lavishag4193/csv2jsonx/internal/parser"
)

type ConversionService struct {
	db *database.PostgresDB
}

func NewConversionService(db *database.PostgresDB) *ConversionService {
	return &ConversionService{
		db: db,
	}
}

// ProcessCSVFile reads a CSV file and converts it to JSON
func (s *ConversionService) ProcessCSVFile(filePath string) ([]byte, error) {
	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Parse CSV
	records, err := parser.ParseCSV(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV: %w", err)
	}

	// Convert to JSON
	jsonData, err := json.Marshal(records)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to JSON: %w", err)
	}

	return jsonData, nil
}

// ProcessCSVReader reads a CSV from an io.Reader and converts it to JSON
func (s *ConversionService) ProcessCSVReader(r io.Reader) ([]byte, error) {
	return s.ProcessCSVReaderWithFilename(r, "")
}

// ProcessCSVReaderWithFilename reads a CSV from an io.Reader, converts it to JSON, and saves to database
func (s *ConversionService) ProcessCSVReaderWithFilename(r io.Reader, filename string) ([]byte, error) {
	// Parse CSV
	records, err := parser.ParseCSV(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV: %w", err)
	}

	// Convert to JSON
	jsonData, err := json.Marshal(records)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to JSON: %w", err)
	}

	// Save to database if db is available
	if s.db != nil {
		if err := s.db.InsertCSVData(filename, records); err != nil {
			return nil, fmt.Errorf("failed to save to database: %w", err)
		}
	}

	return jsonData, nil
}

// GetAllData retrieves all CSV data from the database
func (s *ConversionService) GetAllData() ([]map[string]interface{}, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	return s.db.GetAllCSVData()
}

// GetDataByID retrieves CSV data by ID from the database
func (s *ConversionService) GetDataByID(id int) (map[string]interface{}, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	return s.db.GetCSVDataByID(id)
}
