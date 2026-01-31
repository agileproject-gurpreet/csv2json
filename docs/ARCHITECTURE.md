
# Architecture Document: csv2json

## Overview

The `csv2json` project is a Go-based tool designed to convert CSV files into JSON format via a RESTful API. It is structured for extensibility, maintainability, and ease of integration with databases (e.g., PostgreSQL).

## High-Level Architecture

- **API Layer** (`cmd/api/main.go`): Handles HTTP requests, routing, and response formatting.
- **Handler Layer** (`internal/handler/`): Contains logic for processing API requests, invoking services, and error handling.
- **Service Layer** (`internal/service/`): Implements business logic, orchestrates parsing and conversion, and interacts with the database layer.
- **Parser Layer** (`internal/parser/`): Responsible for parsing CSV files and transforming them into Go data structures.
- **Database Layer** (`internal/database/`): Manages database connections and operations (e.g., PostgreSQL integration).
- **Package Layer** (`pkg/csv2jsonx/`): Provides reusable conversion utilities for CSV to JSON transformation.

---

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         CLIENT                                   │
│                    (curl, Postman, Browser)                      │
└────────────┬────────────────────────────────────────────────────┘
             │
             │ HTTP Requests
             │
             ▼
┌─────────────────────────────────────────────────────────────────┐
│                     API SERVER (Go)                              │
│                                                                  │
│  ┌────────────────────────────────────────────────────────┐    │
│  │              HTTP Handlers                              │    │
│  │  - UploadCSV()                                          │    │
│  │  - GetAllData()                                         │    │
│  │  - GetDataByID()                                        │    │
│  │  - Health()                                             │    │
│  └─────────────┬──────────────────────────────────────────┘    │
│                │                                                 │
│                ▼                                                 │
│  ┌────────────────────────────────────────────────────────┐    │
│  │          Conversion Service                             │    │
│  │  - ProcessCSVReaderWithFilename()                       │    │
│  │  - GetAllData()                                         │    │
│  │  - GetDataByID()                                        │    │
│  └─────────────┬──────────────────────────────────────────┘    │
│                │                                                 │
│                ▼                                                 │
│  ┌────────────────────────────────────────────────────────┐    │
│  │            CSV Parser                                   │    │
│  │  - ParseCSV()                                           │    │
│  └─────────────┬──────────────────────────────────────────┘    │
│                │                                                 │
│                ▼                                                 │
│  ┌────────────────────────────────────────────────────────┐    │
│  │          Database Layer                                 │    │
│  │  - InsertCSVData()                                      │    │
│  │  - GetAllCSVData()                                      │    │
│  │  - GetCSVDataByID()                                     │    │
│  └─────────────┬──────────────────────────────────────────┘    │
│                │                                                 │
└────────────────┼─────────────────────────────────────────────────┘
                 │
                 │ SQL Queries
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                    PostgreSQL Database                           │
│                                                                  │
│  ┌────────────────────────────────────────────────────────┐    │
│  │                  csv_data Table                         │    │
│  │  ┌──────────────────────────────────────────────┐      │    │
│  │  │  id (SERIAL PRIMARY KEY)                      │      │    │
│  │  │  filename (VARCHAR)                           │      │    │
│  │  │  data (JSONB) ◄── Converted CSV as JSON       │      │    │
│  │  │  created_at (TIMESTAMP)                       │      │    │
│  │  └──────────────────────────────────────────────┘      │    │
│  └────────────────────────────────────────────────────────┘    │
│                                                                  │
│  Indexes:                                                        │
│  - idx_csv_data_created_at                                      │
│  - idx_csv_data_filename                                        │
└─────────────────────────────────────────────────────────────────┘
```

## Request Flow - Upload CSV

```
1. Client                    2. Handler                3. Service
   │                            │                         │
   │ POST /api/upload           │                         │
   │ multipart/form-data        │                         │
   │──────────────────────────►│                         │
   │                            │                         │
   │                            │ Extract file            │
   │                            │ Get filename            │
   │                            │                         │
   │                            │ ProcessCSVReader()      │
   │                            │───────────────────────►│
   │                            │                         │
                                                          │
4. Parser                      5. Database               │
   │                              │                       │
   │ ParseCSV()                   │                       │
   │◄─────────────────────────────│                       │
   │                              │                       │
   │ Return []map[string]string   │                       │
   │──────────────────────────────►                       │
   │                              │                       │
   │                              │ InsertCSVData()       │
   │                              │◄──────────────────────│
   │                              │                       │
   │                              │ Store as JSONB        │
   │                              │                       │
   │                              │ Return success        │
   │                              │───────────────────────►
   │                              │                       │
   │                              │ Return JSON           │
   │                              │◄──────────────────────│
   │                              │                       │
   │ HTTP 200 + JSON              │                       │
   │◄─────────────────────────────│                       │
```

## Request Flow - Retrieve Data

```
Client                   Handler                Service              Database
  │                         │                      │                    │
  │ GET /api/data           │                      │                    │
  │───────────────────────►│                      │                    │
  │                         │                      │                    │
  │                         │ GetAllData()         │                    │
  │                         │─────────────────────►│                    │
  │                         │                      │                    │
  │                         │                      │ GetAllCSVData()    │
  │                         │                      │───────────────────►│
  │                         │                      │                    │
  │                         │                      │  SELECT * FROM...  │
  │                         │                      │                    │
  │                         │                      │ Return records     │
  │                         │                      │◄───────────────────│
  │                         │                      │                    │
  │                         │ Return data          │                    │
  │                         │◄─────────────────────│                    │
  │                         │                      │                    │
  │ HTTP 200 + JSON         │                      │                    │
  │◄────────────────────────│                      │                    │
```


## Component Responsibilities


### 1. API Layer (`cmd/api/main.go`)
- **Responsibility**: Entry point for HTTP requests, sets up routing and server configuration.
- **Functions**:
  - Define RESTful endpoints
  - Start and manage the HTTP server
  - Integrate middleware (logging, CORS, etc.)

### 2. Handler Layer (`internal/handler/`)
- **Responsibility**: Handle HTTP requests and responses
- **Functions**:
  - Parse multipart form data
  - Extract files and parameters
  - Call service layer
  - Format HTTP responses
  - Log requests

### 3. Service Layer (`internal/service/`)
- **Responsibility**: Business logic and orchestration
- **Functions**:
  - Coordinate CSV parsing
  - Manage database operations
  - Handle errors
  - Return formatted data

### 4. Parser Layer (`internal/parser/`)
- **Responsibility**: CSV parsing
- **Functions**:
  - Read CSV headers
  - Parse rows into maps
  - Handle CSV format variations

### 5. Database Layer (`internal/database/`)
- **Responsibility**: PostgreSQL operations
- **Functions**:
  - Connection management
  - Schema initialization
  - CRUD operations
  - Query execution

### 6. Package Layer (`pkg/csv2jsonx/`)
- **Responsibility**: Reusable CSV to JSON conversion utilities
- **Functions**:
  - Convert parsed CSV data to JSON
  - Provide utility functions for data transformation

---


## Data Flow

1. **API Request**: User sends a CSV file via HTTP POST.
2. **Handler**: Validates and forwards the request to the service.
3. **Service**: Orchestrates parsing, conversion, and (optionally) database storage.
4. **Parser**: Reads and parses the CSV file.
5. **Converter**: Transforms parsed data into JSON.
6. **Database**: (Optional) Stores or retrieves conversion results.
7. **Response**: JSON result is returned to the user.

---

## Extensibility
- New parsers or converters can be added in `pkg/csv2jsonx/`.
- Additional database backends can be implemented in `internal/database/`.
- API endpoints can be extended in `cmd/api/main.go` and `internal/handler/`.

## Testing
- Unit tests are located in `internal/handler/tests/`, `internal/parser/tests/`, and `internal/service/tests/`.

---

## Data Transformations

```
CSV File                  Parsed Data              JSON Output
┌──────────────┐         ┌────────────────┐       ┌──────────────┐
│ name,age     │         │ []map[string]  │       │ [            │
│ John,30      │  ────►  │  string{       │ ────► │   {          │
│ Jane,25      │         │    {"name":    │       │     "name":  │
└──────────────┘         │     "John",    │       │     "John",  │
                         │     "age":"30"}│       │     "age":   │
                         │    {"name":    │       │     "30"     │
                         │     "Jane",    │       │   },         │
                         │     "age":"25"}│       │   {          │
                         │  }             │       │     "name":  │
                         └────────────────┘       │     "Jane",  │
                                                  │     "age":   │
                         Stored in PostgreSQL:   │     "25"     │
                         ┌────────────────┐       │   }          │
                         │ JSONB:         │       │ ]            │
                         │ [{"name":...}] │       └──────────────┘
                         └────────────────┘
```

## Technology Stack

```
┌─────────────────────────────────────────────────────────┐
│                    Application Layer                     │
│                                                          │
│  Language:        Go 1.25.4                             │
│  Web Framework:   net/http (standard library)           │
│  Router:          http.ServeMux                         │
│  Logging:         log (standard library)                │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│                    Database Layer                        │
│                                                          │
│  Database:        PostgreSQL 12+                        │
│  Driver:          github.com/lib/pq v1.10.9             │
│  Interface:       database/sql                          │
│  Data Format:     JSONB                                 │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│                    Infrastructure                        │
│                                                          │
│  Config:          Environment Variables (.env)          │
│  Connection Pool: 25 max open, 5 max idle               │
│  Indexes:         B-tree (filename, created_at)         │
└─────────────────────────────────────────────────────────┘
```

## Database Schema Details

```sql
┌──────────────────────────────────────────────────────────┐
│                      csv_data                            │
├──────────────┬────────────────┬──────────────────────────┤
│   Column     │      Type      │       Constraints        │
├──────────────┼────────────────┼──────────────────────────┤
│ id           │ SERIAL         │ PRIMARY KEY              │
│ filename     │ VARCHAR(255)   │                          │
│ data         │ JSONB          │ NOT NULL                 │
│ created_at   │ TIMESTAMP      │ DEFAULT CURRENT_TIMESTAMP│
└──────────────┴────────────────┴──────────────────────────┘

Indexes:
  - PRIMARY KEY (id)
  - idx_csv_data_created_at (B-tree on created_at)
  - idx_csv_data_filename (B-tree on filename)

Sample data structure:
{
  "id": 1,
  "filename": "sales.csv",
  "data": [
    {"name": "Product A", "price": "100", "quantity": "50"},
    {"name": "Product B", "price": "200", "quantity": "30"}
  ],
  "created_at": "2026-01-27T10:30:00Z"
}
```

## Scalability Considerations

### Connection Pooling
- Max Open Connections: 25
- Max Idle Connections: 5
- Connection Lifetime: 5 minutes

### Database Optimization
- JSONB indexing for fast queries
- B-tree indexes on frequently queried columns
- Efficient JSON storage format

### API Performance
- Direct streaming of responses
- Minimal memory overhead
- Concurrent request handling

## Security Features

- Environment-based configuration
- SSL/TLS support for database connections
- Input validation
- SQL injection prevention (prepared statements)
- Connection pooling prevents resource exhaustion
