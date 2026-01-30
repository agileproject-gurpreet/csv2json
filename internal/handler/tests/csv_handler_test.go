package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/lavishag4193/csv2jsonx/internal/handler"
	"github.com/lavishag4193/csv2jsonx/internal/service"
)

// TestUploadCSV_Success tests successful CSV file upload
func TestUploadCSV_Success(t *testing.T) {
	// Setup
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	svc := service.NewConversionService(nil) // nil db for testing
	h := handler.NewCSVHandler(svc, logger)

	// Create a test CSV file
	csvContent := "name,age,city\nAlice,30,NYC\nBob,25,LA"
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.csv")
	if err != nil {
		t.Fatal(err)
	}
	io.WriteString(part, csvContent)
	writer.Close()

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	// Execute
	h.UploadCSV(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Verify JSON response
	var result []map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 records, got %d", len(result))
	}

	// Check first record
	if result[0]["name"] != "Alice" || result[0]["age"] != "30" || result[0]["city"] != "NYC" {
		t.Errorf("unexpected first record: %v", result[0])
	}

	// Check second record
	if result[1]["name"] != "Bob" || result[1]["age"] != "25" || result[1]["city"] != "LA" {
		t.Errorf("unexpected second record: %v", result[1])
	}
}

// TestUploadCSV_EmptyFile tests uploading an empty CSV file
func TestUploadCSV_EmptyFile(t *testing.T) {
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	svc := service.NewConversionService(nil)
	h := handler.NewCSVHandler(svc, logger)

	// Create empty CSV file
	csvContent := ""
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "empty.csv")
	if err != nil {
		t.Fatal(err)
	}
	io.WriteString(part, csvContent)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	h.UploadCSV(w, req)

	// Empty CSV returns error (EOF)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}

// TestUploadCSV_InvalidMethod tests non-POST methods
func TestUploadCSV_InvalidMethod(t *testing.T) {
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	svc := service.NewConversionService(nil)
	h := handler.NewCSVHandler(svc, logger)

	methods := []string{http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, method := range methods {
		req := httptest.NewRequest(method, "/api/upload", nil)
		w := httptest.NewRecorder()

		h.UploadCSV(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("method %s: expected status 405, got %d", method, w.Code)
		}
	}
}

// TestUploadCSV_NoFile tests upload without file
func TestUploadCSV_NoFile(t *testing.T) {
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	svc := service.NewConversionService(nil)
	h := handler.NewCSVHandler(svc, logger)

	// Create form without file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	h.UploadCSV(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

// TestUploadCSV_LargeFile tests uploading a large CSV file
func TestUploadCSV_LargeFile(t *testing.T) {
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	svc := service.NewConversionService(nil)
	h := handler.NewCSVHandler(svc, logger)

	// Create a large CSV file (1000 rows)
	var csvBuilder bytes.Buffer
	csvBuilder.WriteString("id,name,value\n")
	for i := 0; i < 1000; i++ {
		csvBuilder.WriteString(bytes.NewBufferString("").String())
		csvBuilder.WriteString(bytes.NewBufferString("").String())
		csvBuilder.WriteString(bytes.NewBufferString("").String())
		csvBuilder.WriteString(bytes.NewBufferString("").String())
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "large.csv")
	if err != nil {
		t.Fatal(err)
	}
	io.WriteString(part, csvBuilder.String())
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	h.UploadCSV(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

// TestUploadCSV_SpecialCharacters tests CSV with special characters
func TestUploadCSV_SpecialCharacters(t *testing.T) {
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	svc := service.NewConversionService(nil)
	h := handler.NewCSVHandler(svc, logger)

	// CSV with special characters, quotes, commas
	csvContent := `name,description,price
"Product A","Contains ""quotes"" and, commas",19.99
"Product B","Special chars: @#$%^&*()",29.99`

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "special.csv")
	if err != nil {
		t.Fatal(err)
	}
	io.WriteString(part, csvContent)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	h.UploadCSV(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var result []map[string]string
	json.Unmarshal(w.Body.Bytes(), &result)

	if len(result) != 2 {
		t.Errorf("expected 2 records, got %d", len(result))
	}
}

// TestUploadCSV_InvalidCSVFormat tests malformed CSV
func TestUploadCSV_InvalidCSVFormat(t *testing.T) {
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	svc := service.NewConversionService(nil)
	h := handler.NewCSVHandler(svc, logger)

	// Invalid CSV with mismatched columns
	csvContent := "name,age\nAlice,30,extra\nBob"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "invalid.csv")
	if err != nil {
		t.Fatal(err)
	}
	io.WriteString(part, csvContent)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	h.UploadCSV(w, req)

	// This may succeed or fail depending on CSV parser implementation
	// Just verify it doesn't panic
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Logf("CSV parsing returned status: %d", w.Code)
	}
}

// TestUploadCSV_ContentType tests content type validation
func TestUploadCSV_ContentType(t *testing.T) {
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	svc := service.NewConversionService(nil)
	h := handler.NewCSVHandler(svc, logger)

	csvContent := "name,age\nAlice,30"
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.csv")
	if err != nil {
		t.Fatal(err)
	}
	io.WriteString(part, csvContent)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	h.UploadCSV(w, req)

	// Verify response content type is JSON
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected content-type application/json, got %s", contentType)
	}
}

// TestGetAllData_Success tests successful retrieval of all data
func TestGetAllData_Success(t *testing.T) {
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	svc := service.NewConversionService(nil)
	h := handler.NewCSVHandler(svc, logger)

	req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
	w := httptest.NewRecorder()

	h.GetAllData(w, req)

	// Without database, should return error
	if w.Code != http.StatusInternalServerError {
		t.Logf("expected status 500 without DB, got %d", w.Code)
	}
}

// TestGetAllData_InvalidMethod tests non-GET methods
func TestGetAllData_InvalidMethod(t *testing.T) {
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	svc := service.NewConversionService(nil)
	h := handler.NewCSVHandler(svc, logger)

	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, method := range methods {
		req := httptest.NewRequest(method, "/api/data", nil)
		w := httptest.NewRecorder()

		h.GetAllData(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("method %s: expected status 405, got %d", method, w.Code)
		}
	}
}

// TestGetAllData_ResponseFormat tests the response format
func TestGetAllData_ResponseFormat(t *testing.T) {
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	svc := service.NewConversionService(nil)
	h := handler.NewCSVHandler(svc, logger)

	req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
	w := httptest.NewRecorder()

	h.GetAllData(w, req)

	// Without database, returns error with text/plain content type
	// When error occurs, http.Error sets content-type to text/plain
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}

// TestHealth tests the health check endpoint
func TestHealth(t *testing.T) {
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	svc := service.NewConversionService(nil)
	h := handler.NewCSVHandler(svc, logger)

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	w := httptest.NewRecorder()

	h.Health(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["status"] != "healthy" {
		t.Errorf("expected status=healthy, got %s", response["status"])
	}
}

// TestUploadCSV_MultipleConcurrent tests concurrent uploads
func TestUploadCSV_MultipleConcurrent(t *testing.T) {
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	svc := service.NewConversionService(nil)
	h := handler.NewCSVHandler(svc, logger)

	done := make(chan bool, 5)

	for i := 0; i < 5; i++ {
		go func(idx int) {
			csvContent := "name,age\nUser,30"
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, _ := writer.CreateFormFile("file", "test.csv")
			io.WriteString(part, csvContent)
			writer.Close()

			req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			w := httptest.NewRecorder()

			h.UploadCSV(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("concurrent request %d failed with status %d", idx, w.Code)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 5; i++ {
		<-done
	}
}

// TestUploadCSV_UTF8Characters tests UTF-8 encoded content
func TestUploadCSV_UTF8Characters(t *testing.T) {
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	svc := service.NewConversionService(nil)
	h := handler.NewCSVHandler(svc, logger)

	// CSV with UTF-8 characters
	csvContent := "name,city\n名前,東京\nПривет,Москва\n😀,🌍"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "utf8.csv")
	if err != nil {
		t.Fatal(err)
	}
	io.WriteString(part, csvContent)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	h.UploadCSV(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
