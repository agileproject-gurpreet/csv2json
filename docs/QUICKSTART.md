# Quick Start Guide - PostgreSQL Integration

## Prerequisites

- Go 1.25.4+ installed
- PostgreSQL 12+ installed and running
- Git (for cloning the repository)

## Setup (5 minutes)

### Step 1: Clone and Setup
```bash
git clone <repository-url>
cd csv2json-api
```

### Step 2: Create Database
```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE csv2json;

# Exit psql
\q
```

### Step 3: Configure Environment
```bash
# Copy example environment file
cp .env.example .env

# Edit .env with your credentials
# Use your favorite text editor
notepad .env  # Windows
nano .env     # Linux/Mac
```

Update these values in `.env`:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=YOUR_PASSWORD_HERE
DB_NAME=csv2json
DB_SSLMODE=disable
PORT=8080
```

### Step 4: Install Dependencies
```bash
go mod download
```

### Step 5: Run the Application
```bash
go run cmd/api/main.go
```

You should see:
```
[CSV2JSON-API] Starting CSV2JSON API server...
[CSV2JSON-API] Successfully connected to PostgreSQL database
[CSV2JSON-API] Database schema initialized
[CSV2JSON-API] Server starting on port 8080
[CSV2JSON-API] Available endpoints:
[CSV2JSON-API]   POST /api/upload     - Upload CSV file
[CSV2JSON-API]   GET  /api/data       - Get all stored CSV data
[CSV2JSON-API]   GET  /api/data/id    - Get CSV data by ID (requires ?id=<id>)
[CSV2JSON-API]   GET  /api/health     - Health check
```

## Testing (2 minutes)

### Test 1: Health Check
```bash
curl http://localhost:8080/api/health
```

Expected response:
```json
{"status":"healthy"}
```

### Test 2: Upload CSV
```bash
curl -X POST http://localhost:8080/api/upload \
  -F "file=@examples/sample.csv"
```

Expected response:
```json
[
  {
    "name": "John",
    "age": "30",
    "city": "New York"
  },
  ...
]
```

### Test 3: Retrieve All Data
```bash
curl http://localhost:8080/api/data
```

Expected response:
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

### Test 4: Retrieve Specific Record
```bash
curl http://localhost:8080/api/data/id?id=1
```

Expected response:
```json
{
  "id": 1,
  "filename": "sample.csv",
  "data": [...],
  "created_at": "2026-01-27T10:30:00Z"
}
```

## Verify Database (Optional)

```bash
# Connect to database
psql -U postgres -d csv2json

# Check table exists
\dt

# View data
SELECT id, filename, created_at FROM csv_data;

# View JSON data
SELECT data FROM csv_data WHERE id = 1;

# Exit
\q
```

## Troubleshooting

### Issue: "failed to connect to database"
**Solution**: 
1. Verify PostgreSQL is running:
   ```bash
   # Windows
   Get-Service postgresql*
   
   # Linux
   sudo service postgresql status
   
   # Mac
   brew services list
   ```
2. Check credentials in `.env`
3. Ensure database exists: `psql -U postgres -l`

### Issue: "database does not exist"
**Solution**: Create the database manually:
```bash
psql -U postgres -c "CREATE DATABASE csv2json;"
```

### Issue: Port 8080 already in use
**Solution**: Change the port in `.env`:
```env
PORT=8081
```

### Issue: "role does not exist"
**Solution**: Create the PostgreSQL user:
```bash
psql -U postgres
CREATE USER your_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE csv2json TO your_user;
\q
```

## What's Happening Under the Hood?

1. **Application starts** → Connects to PostgreSQL
2. **Schema initialization** → Creates `csv_data` table if it doesn't exist
3. **CSV upload** → Parses CSV → Converts to JSON → Saves to database
4. **Data retrieval** → Queries PostgreSQL → Returns JSON

## Data Flow

```
CSV File → API Upload Endpoint
    ↓
Parse CSV to Map
    ↓
Convert to JSON
    ↓
Store in PostgreSQL (JSONB)
    ↓
Return JSON to Client
```

## File Structure

```
csv2json-api/
├── cmd/api/main.go              # Application entry point
├── internal/
│   ├── database/
│   │   └── postgres.go          # PostgreSQL integration
│   ├── handler/
│   │   └── csv_handler.go       # HTTP handlers
│   ├── parser/
│   │   └── csv_parser.go        # CSV parsing logic
│   └── service/
│       └── conversion_service.go # Business logic
├── docs/
│   ├── setup.sql                # Database setup script
│   ├── postgres_integration.md  # Detailed guide
│   └── IMPLEMENTATION_SUMMARY.md # Technical summary
├── .env.example                 # Environment template
└── README.md                    # Main documentation
```

## Next Steps

- Read [docs/postgres_integration.md](postgres_integration.md) for advanced features
- Explore JSONB querying capabilities
- Add authentication for production use
- Set up database backups

## Production Checklist

Before deploying to production:

- [ ] Change default passwords
- [ ] Enable SSL/TLS (`DB_SSLMODE=require`)
- [ ] Set up database backups
- [ ] Configure firewall rules
- [ ] Add API authentication
- [ ] Set up monitoring and logging
- [ ] Use environment-specific `.env` files
- [ ] Never commit `.env` to version control

## Support

For issues or questions:
1. Check [docs/postgres_integration.md](postgres_integration.md)
2. Review [docs/IMPLEMENTATION_SUMMARY.md](IMPLEMENTATION_SUMMARY.md)
3. Check PostgreSQL logs: `tail -f /var/log/postgresql/postgresql-*.log`
4. Check application logs in the console output

---

**Ready to go!** 🚀 Your CSV2JSON API with PostgreSQL integration is now running!
