# Architecture Overview

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

### 1. HTTP Handlers (`internal/handler/`)
- **Responsibility**: Handle HTTP requests and responses
- **Functions**:
  - Parse multipart form data
  - Extract files and parameters
  - Call service layer
  - Format HTTP responses
  - Log requests

### 2. Conversion Service (`internal/service/`)
- **Responsibility**: Business logic
- **Functions**:
  - Coordinate CSV parsing
  - Manage database operations
  - Handle errors
  - Return formatted data

### 3. CSV Parser (`internal/parser/`)
- **Responsibility**: CSV parsing
- **Functions**:
  - Read CSV headers
  - Parse rows into maps
  - Handle CSV format variations

### 4. Database Layer (`internal/database/`)
- **Responsibility**: PostgreSQL operations
- **Functions**:
  - Connection management
  - Schema initialization
  - CRUD operations
  - Query execution

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
