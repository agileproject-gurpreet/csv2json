# csv2json
Fast, simple, and extensible CSV to JSON converter for developers with PostgreSQL integration.

This is the public repository containing installation steps, usage examples, and documentation

## Features

- Upload CSV files via REST API
- Automatic CSV to JSON conversion
- PostgreSQL database storage for all converted data
- JSONB support for efficient querying
- Connection pooling and optimized database operations

## Prerequisites

- Go 1.24 or higher
- PostgreSQL 12 or higher

## Database Setup

1. Install PostgreSQL if not already installed
2. Create a database for the application:
```sql
CREATE DATABASE csv2json;
```

3. Configure database connection using environment variables (see `.env.example`)

The application will automatically create the required tables on startup.

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd csv2json-api
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your PostgreSQL credentials
```

4. Run the application:
```bash
go run cmd/api/main.go
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `postgres` |
| `DB_NAME` | Database name | `csv2json` |
| `DB_SSLMODE` | SSL mode for connection | `disable` |
| `PORT` | Server port | `8080` |

## API Endpoints

### Upload CSV
```
POST /api/upload
Content-Type: multipart/form-data
```

Upload a CSV file and convert it to JSON. The data is automatically saved to PostgreSQL.

**Example:**
```bash
curl -X POST http://localhost:8080/api/upload \
  -F "file=@sample.csv"
```

**Response:**
```json
[
  {
    "column1": "value1",
    "column2": "value2"
  }
]
```

### Get All Data
```
GET /api/data
```

Retrieve all stored CSV data from the database.

**Example:**
```bash
curl http://localhost:8080/api/data
```

**Response:**
```json
[
  {
    "id": 1,
    "filename": "sample.csv",
    "data": [...],
    "created_at": "2026-01-27T10:30:00Z"
  }
]
```

### Get Data by ID
```
GET /api/data/id?id={id}
```

Retrieve a specific CSV data record by its ID.

**Example:**
```bash
curl http://localhost:8080/api/data/id?id=1
```

**Response:**
```json
{
  "id": 1,
  "filename": "sample.csv",
  "data": [...],
  "created_at": "2026-01-27T10:30:00Z"
}
```

### Health Check
```
GET /api/health
```

Check if the API is running.

**Response:**
```json
{
  "status": "healthy"
}
```

## Database Schema

The application creates the following table:

```sql
CREATE TABLE csv_data (
    id SERIAL PRIMARY KEY,
    filename VARCHAR(255),
    data JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

The `data` column stores the converted JSON data as JSONB, allowing for efficient querying and indexing.

## Development

Run tests:
```bash
go test ./...
```

## License

See LICENSE file for details.

