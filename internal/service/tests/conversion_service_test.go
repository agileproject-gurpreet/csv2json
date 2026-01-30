package tests

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/agileproject-gurpreet/csv2json/internal/service"
)

// TestNewConversionService tests service creation
func TestNewConversionService(t *testing.T) {
	svc := service.NewConversionService(nil)
	if svc == nil {
		t.Fatal("expected non-nil service")
	}
}

// TestProcessCSVReader_Success tests successful CSV processing from reader
func TestProcessCSVReader_Success(t *testing.T) {
	svc := service.NewConversionService(nil)

	csvData := "name,age,city\nAlice,30,NYC\nBob,25,LA"
	reader := strings.NewReader(csvData)

	jsonData, err := svc.ProcessCSVReader(reader)
	if err != nil {
		t.Fatalf("ProcessCSVReader failed: %v", err)
	}

	// Verify JSON format
	var result []map[string]string
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 records, got %d", len(result))
	}

	// Verify first record
	if result[0]["name"] != "Alice" {
		t.Errorf("expected name=Alice, got %s", result[0]["name"])
	}
	if result[0]["age"] != "30" {
		t.Errorf("expected age=30, got %s", result[0]["age"])
	}
	if result[0]["city"] != "NYC" {
		t.Errorf("expected city=NYC, got %s", result[0]["city"])
	}

	// Verify second record
	if result[1]["name"] != "Bob" {
		t.Errorf("expected name=Bob, got %s", result[1]["name"])
	}
	if result[1]["age"] != "25" {
		t.Errorf("expected age=25, got %s", result[1]["age"])
	}
	if result[1]["city"] != "LA" {
		t.Errorf("expected city=LA, got %s", result[1]["city"])
	}
}

// TestProcessCSVReader_EmptyCSV tests empty CSV input
func TestProcessCSVReader_EmptyCSV(t *testing.T) {
	svc := service.NewConversionService(nil)

	csvData := ""
	reader := strings.NewReader(csvData)

	_, err := svc.ProcessCSVReader(reader)
	if err == nil {
		t.Fatal("expected error for empty CSV")
	}

	// Empty CSV returns EOF error
	if !strings.Contains(err.Error(), "EOF") {
		t.Errorf("expected EOF error, got: %v", err)
	}
}

// TestProcessCSVReader_OnlyHeaders tests CSV with only headers
func TestProcessCSVReader_OnlyHeaders(t *testing.T) {
	svc := service.NewConversionService(nil)

	csvData := "name,age,city"
	reader := strings.NewReader(csvData)

	jsonData, err := svc.ProcessCSVReader(reader)
	if err != nil {
		t.Fatalf("ProcessCSVReader failed: %v", err)
	}

	var result []map[string]string
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	// Only headers should return empty array
	if len(result) != 0 {
		t.Errorf("expected 0 records for headers-only CSV, got %d", len(result))
	}
}

// TestProcessCSVReader_SingleRecord tests CSV with single record
func TestProcessCSVReader_SingleRecord(t *testing.T) {
	svc := service.NewConversionService(nil)

	csvData := "name,age\nAlice,30"
	reader := strings.NewReader(csvData)

	jsonData, err := svc.ProcessCSVReader(reader)
	if err != nil {
		t.Fatalf("ProcessCSVReader failed: %v", err)
	}

	var result []map[string]string
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("expected 1 record, got %d", len(result))
	}

	if result[0]["name"] != "Alice" || result[0]["age"] != "30" {
		t.Errorf("unexpected record data: %v", result[0])
	}
}

// TestProcessCSVReader_MultipleColumns tests CSV with many columns
func TestProcessCSVReader_MultipleColumns(t *testing.T) {
	svc := service.NewConversionService(nil)

	csvData := "col1,col2,col3,col4,col5\nval1,val2,val3,val4,val5"
	reader := strings.NewReader(csvData)

	jsonData, err := svc.ProcessCSVReader(reader)
	if err != nil {
		t.Fatalf("ProcessCSVReader failed: %v", err)
	}

	var result []map[string]string
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("expected 1 record, got %d", len(result))
	}

	// Verify all columns are present
	if len(result[0]) != 5 {
		t.Errorf("expected 5 columns, got %d", len(result[0]))
	}
}

// TestProcessCSVReaderWithFilename_Success tests CSV processing with filename
func TestProcessCSVReaderWithFilename_Success(t *testing.T) {
	svc := service.NewConversionService(nil) // No DB

	csvData := "name,age\nAlice,30"
	reader := strings.NewReader(csvData)

	jsonData, err := svc.ProcessCSVReaderWithFilename(reader, "test.csv")
	if err != nil {
		t.Fatalf("ProcessCSVReaderWithFilename failed: %v", err)
	}

	var result []map[string]string
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("expected 1 record, got %d", len(result))
	}
}

// TestProcessCSVReaderWithFilename_NoDatabase tests CSV processing without database
func TestProcessCSVReaderWithFilename_NoDatabase(t *testing.T) {
	svc := service.NewConversionService(nil)

	csvData := "name,age\nAlice,30\nBob,25"
	reader := strings.NewReader(csvData)

	// Should succeed even without database
	jsonData, err := svc.ProcessCSVReaderWithFilename(reader, "test.csv")
	if err != nil {
		t.Fatalf("expected success without DB, got error: %v", err)
	}

	if jsonData == nil {
		t.Fatal("expected non-nil JSON data")
	}
}

// TestProcessCSVReader_SpecialCharacters tests CSV with special characters
func TestProcessCSVReader_SpecialCharacters(t *testing.T) {
	svc := service.NewConversionService(nil)

	csvData := `name,description
"Product A","Contains ""quotes"" and, commas"
"Product B","Special: @#$%"`

	reader := strings.NewReader(csvData)

	jsonData, err := svc.ProcessCSVReader(reader)
	if err != nil {
		t.Fatalf("ProcessCSVReader failed: %v", err)
	}

	var result []map[string]string
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 records, got %d", len(result))
	}

	// Verify special characters are preserved
	if !strings.Contains(result[0]["description"], "quotes") {
		t.Errorf("expected quotes in description, got: %s", result[0]["description"])
	}
}

// TestProcessCSVReader_UTF8Content tests CSV with UTF-8 characters
func TestProcessCSVReader_UTF8Content(t *testing.T) {
	svc := service.NewConversionService(nil)

	csvData := "name,city\n名前,東京\nПривет,Москва"
	reader := strings.NewReader(csvData)

	jsonData, err := svc.ProcessCSVReader(reader)
	if err != nil {
		t.Fatalf("ProcessCSVReader failed: %v", err)
	}

	var result []map[string]string
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 records, got %d", len(result))
	}

	// Verify UTF-8 characters are preserved
	if result[0]["name"] != "名前" {
		t.Errorf("expected 名前, got %s", result[0]["name"])
	}
	if result[0]["city"] != "東京" {
		t.Errorf("expected 東京, got %s", result[0]["city"])
	}
}

// TestProcessCSVReader_Whitespace tests CSV with whitespace
func TestProcessCSVReader_Whitespace(t *testing.T) {
	svc := service.NewConversionService(nil)

	csvData := "name,age\n  Alice  ,  30  \n  Bob  ,  25  "
	reader := strings.NewReader(csvData)

	jsonData, err := svc.ProcessCSVReader(reader)
	if err != nil {
		t.Fatalf("ProcessCSVReader failed: %v", err)
	}

	var result []map[string]string
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 records, got %d", len(result))
	}
}

// TestProcessCSVReader_NumericValues tests CSV with numeric values
func TestProcessCSVReader_NumericValues(t *testing.T) {
	svc := service.NewConversionService(nil)

	csvData := "id,price,quantity\n1,19.99,100\n2,29.99,50"
	reader := strings.NewReader(csvData)

	jsonData, err := svc.ProcessCSVReader(reader)
	if err != nil {
		t.Fatalf("ProcessCSVReader failed: %v", err)
	}

	var result []map[string]string
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 records, got %d", len(result))
	}

	// Verify numeric values are preserved as strings
	if result[0]["price"] != "19.99" {
		t.Errorf("expected price=19.99, got %s", result[0]["price"])
	}
}

// TestProcessCSVReader_BooleanValues tests CSV with boolean-like values
func TestProcessCSVReader_BooleanValues(t *testing.T) {
	svc := service.NewConversionService(nil)

	csvData := "name,active,verified\nAlice,true,false\nBob,false,true"
	reader := strings.NewReader(csvData)

	jsonData, err := svc.ProcessCSVReader(reader)
	if err != nil {
		t.Fatalf("ProcessCSVReader failed: %v", err)
	}

	var result []map[string]string
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 records, got %d", len(result))
	}

	// Verify boolean values are preserved as strings
	if result[0]["active"] != "true" {
		t.Errorf("expected active=true, got %s", result[0]["active"])
	}
}

// TestProcessCSVFile_Success tests file-based CSV processing
func TestProcessCSVFile_Success(t *testing.T) {
	svc := service.NewConversionService(nil)

	// Create temporary CSV file
	tmpFile, err := os.CreateTemp("", "test_*.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	csvContent := "name,age\nAlice,30\nBob,25"
	if _, err := tmpFile.WriteString(csvContent); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	jsonData, err := svc.ProcessCSVFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("ProcessCSVFile failed: %v", err)
	}

	var result []map[string]string
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 records, got %d", len(result))
	}
}

// TestProcessCSVFile_FileNotFound tests non-existent file
func TestProcessCSVFile_FileNotFound(t *testing.T) {
	svc := service.NewConversionService(nil)

	_, err := svc.ProcessCSVFile("/nonexistent/file.csv")
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}

	if !strings.Contains(err.Error(), "failed to open file") {
		t.Errorf("unexpected error message: %v", err)
	}
}

// TestProcessCSVFile_EmptyFile tests empty file
func TestProcessCSVFile_EmptyFile(t *testing.T) {
	svc := service.NewConversionService(nil)

	// Create empty temporary file
	tmpFile, err := os.CreateTemp("", "empty_*.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	_, err = svc.ProcessCSVFile(tmpFile.Name())
	if err == nil {
		t.Fatal("expected error for empty file")
	}

	// Empty file returns EOF error
	if !strings.Contains(err.Error(), "EOF") {
		t.Errorf("expected EOF error, got: %v", err)
	}
}

// TestGetAllData_NoDatabase tests GetAllData without database
func TestGetAllData_NoDatabase(t *testing.T) {
	svc := service.NewConversionService(nil)

	_, err := svc.GetAllData()
	if err == nil {
		t.Fatal("expected error when database not initialized")
	}

	if !strings.Contains(err.Error(), "database not initialized") {
		t.Errorf("unexpected error message: %v", err)
	}
}

// TestGetDataByID_NoDatabase tests GetDataByID without database
func TestGetDataByID_NoDatabase(t *testing.T) {
	svc := service.NewConversionService(nil)

	_, err := svc.GetDataByID(1)
	if err == nil {
		t.Fatal("expected error when database not initialized")
	}

	if !strings.Contains(err.Error(), "database not initialized") {
		t.Errorf("unexpected error message: %v", err)
	}
}

// TestProcessCSVReader_LargeCSV tests processing large CSV
func TestProcessCSVReader_LargeCSV(t *testing.T) {
	svc := service.NewConversionService(nil)

	// Generate large CSV (1000 rows)
	var sb strings.Builder
	sb.WriteString("id,name,value\n")
	for i := 0; i < 1000; i++ {
		sb.WriteString("1,test,value\n")
	}

	reader := strings.NewReader(sb.String())

	jsonData, err := svc.ProcessCSVReader(reader)
	if err != nil {
		t.Fatalf("ProcessCSVReader failed: %v", err)
	}

	var result []map[string]string
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if len(result) != 1000 {
		t.Errorf("expected 1000 records, got %d", len(result))
	}
}

// TestProcessCSVReader_EmptyValues tests CSV with empty values
func TestProcessCSVReader_EmptyValues(t *testing.T) {
	svc := service.NewConversionService(nil)

	csvData := "name,age,city\nAlice,,NYC\n,25,\nBob,30,LA"
	reader := strings.NewReader(csvData)

	jsonData, err := svc.ProcessCSVReader(reader)
	if err != nil {
		t.Fatalf("ProcessCSVReader failed: %v", err)
	}

	var result []map[string]string
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("expected 3 records, got %d", len(result))
	}

	// Verify empty values are handled
	if result[0]["age"] != "" {
		t.Errorf("expected empty age, got %s", result[0]["age"])
	}
}

// TestProcessCSVReader_QuotedFields tests CSV with quoted fields
func TestProcessCSVReader_QuotedFields(t *testing.T) {
	svc := service.NewConversionService(nil)

	csvData := `name,address
"John Doe","123 Main St, Apt 4"
"Jane Smith","456 Oak Ave"`

	reader := strings.NewReader(csvData)

	jsonData, err := svc.ProcessCSVReader(reader)
	if err != nil {
		t.Fatalf("ProcessCSVReader failed: %v", err)
	}

	var result []map[string]string
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 records, got %d", len(result))
	}

	// Verify quoted fields are handled correctly
	if result[0]["name"] != "John Doe" {
		t.Errorf("expected 'John Doe', got %s", result[0]["name"])
	}
}
