# CSV2JSON API Usage Guide

## Overview
This API allows you to upload CSV files which will be processed, converted to JSON, and stored in the database with comprehensive logging.

## Endpoints

### 1. Upload CSV File
**Endpoint:** `POST /api/upload`

**Description:** Upload a CSV file for processing. The file will be parsed, converted to JSON, and stored in the database.

**Request:**
- Method: `POST`
- Content-Type: `multipart/form-data`
- Body: Form data with a file field named `file`

**Example using curl:**
```bash
curl -X POST http://localhost:8080/api/upload \
  -F "file=@data.csv"
```

**Example using PowerShell:**
```powershell
$filePath = "C:\path\to\your\data.csv"
$uri = "http://localhost:8080/api/upload"

$form = @{
    file = Get-Item -Path $filePath
}

Invoke-RestMethod -Uri $uri -Method Post -Form $form
```

**Success Response:**
```json
{
  "message": "CSV file processed successfully",
  "filename": "data.csv"
}
```

**Error Response:**
```json
{
  "error": "Failed to process CSV: file already processed"
}
```

### 2. Health Check
**Endpoint:** `GET /api/health`

**Description:** Check if the API server is running.

**Example:**
```bash
curl http://localhost:8080/api/health
```

**Response:**
```json
{
  "status": "healthy"
}
```

## Running the API

### 1. Setup Database
Make sure you have PostgreSQL running and the database created:
```sql
CREATE DATABASE csv2json;

CREATE TABLE csv_data (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file_hash VARCHAR(32) UNIQUE NOT NULL,
    json_data JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 2. Set Environment Variables
Copy `.env.example` to `.env` and configure:
```bash
cp .env.example .env
```

### 3. Install Dependencies
```bash
go mod tidy
```

### 4. Run the Server
```bash
go run cmd/api/main.go
```

The server will start on port 8080 (or the port specified in your environment).

## Logs

The API provides comprehensive logging for all operations:

```
[CSV2JSON-API] 2026/01/27 14:30:15 main.go:22: Starting CSV2JSON API server...
[CSV2JSON-API] 2026/01/27 14:30:15 main.go:32: Connecting to database at localhost:5432...
[CSV2JSON-API] 2026/01/27 14:30:15 main.go:43: Database connection established
[CSV2JSON-API] 2026/01/27 14:30:15 main.go:54: Server starting on port 8080
[CSV2JSON-API] 2026/01/27 14:30:15 main.go:55: Available endpoints:
[CSV2JSON-API] 2026/01/27 14:30:15 main.go:56:   POST /api/upload - Upload CSV file
[CSV2JSON-API] 2026/01/27 14:30:15 main.go:57:   GET  /api/health - Health check
[CSV2JSON-API] 2026/01/27 14:30:20 csv_handler.go:30: Received CSV upload request
[CSV2JSON-API] 2026/01/27 14:30:20 csv_handler.go:48: Processing file: data.csv (size: 1024 bytes)
[CSV2JSON-API] 2026/01/27 14:30:20 csv_handler.go:58: Successfully processed CSV file: data.csv
```

## Testing the API

### Using curl (Linux/Mac/Git Bash)
```bash
# Upload a CSV file
curl -X POST http://localhost:8080/api/upload \
  -F "file=@examples/sample.csv"

# Check health
curl http://localhost:8080/api/health
```

### Using PowerShell (Windows)
```powershell
# Upload a CSV file
$response = Invoke-RestMethod -Uri "http://localhost:8080/api/upload" `
  -Method Post `
  -Form @{ file = Get-Item "examples\sample.csv" }

Write-Output $response

# Check health
Invoke-RestMethod -Uri "http://localhost:8080/api/health" -Method Get
```

## Features

- ✅ **File Upload**: Multi-part form data handling
- ✅ **Duplicate Detection**: Prevents processing the same file twice using MD5 hashing
- ✅ **Comprehensive Logging**: All operations are logged with timestamps
- ✅ **Error Handling**: Detailed error messages for debugging
- ✅ **Health Check**: Monitor API availability
- ✅ **Database Storage**: Stores parsed JSON data in PostgreSQL
