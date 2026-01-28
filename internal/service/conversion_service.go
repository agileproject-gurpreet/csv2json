package service

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/lavishag4193/csv2jsonx/internal/parser"
)

type ConversionService struct{}

func NewConversionService() *ConversionService {
	return &ConversionService{}
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

	return jsonData, nil
}
