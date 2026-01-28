# PostgreSQL Integration Guide

## Overview

The CSV2JSON API now includes full PostgreSQL integration, automatically storing all uploaded CSV files as JSON data in a PostgreSQL database.

## How It Works

1. **Upload**: When a CSV file is uploaded via `/api/upload`, the API:
   - Parses the CSV file
   - Converts it to JSON format
   - Stores the JSON data in PostgreSQL as JSONB
   - Returns the JSON data to the client

2. **Storage**: Data is stored in the `csv_data` table with:
   - `id`: Auto-incrementing primary key
   - `filename`: Original filename
   - `data`: JSONB column containing the converted data
   - `created_at`: Timestamp of when the data was stored

3. **Retrieval**: You can query stored data using:
   - `/api/data` - Get all records
   - `/api/data/id?id=<id>` - Get specific record by ID

## Database Configuration

### Development Setup

For local development, use the default configuration in `.env`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=csv2json
DB_SSLMODE=disable
```

### Production Setup

For production, update the configuration with secure values:

```env
DB_HOST=your-production-host
DB_PORT=5432
DB_USER=your-db-user
DB_PASSWORD=your-secure-password
DB_NAME=csv2json
DB_SSLMODE=require
```

## Database Schema

The application automatically creates the required schema on startup:

```sql
CREATE TABLE IF NOT EXISTS csv_data (
    id SERIAL PRIMARY KEY,
    filename VARCHAR(255),
    data JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_csv_data_created_at ON csv_data(created_at);
CREATE INDEX IF NOT EXISTS idx_csv_data_filename ON csv_data(filename);
```

## Querying JSONB Data

PostgreSQL's JSONB type allows for powerful queries. Here are some examples:

### 1. Query all records
```sql
SELECT * FROM csv_data ORDER BY created_at DESC;
```

### 2. Search within JSON data
```sql
SELECT * FROM csv_data 
WHERE data @> '[{"name": "John"}]'::jsonb;
```

### 3. Extract specific fields
```sql
SELECT id, filename, 
       jsonb_array_length(data) as record_count,
       created_at 
FROM csv_data;
```

### 4. Filter by filename pattern
```sql
SELECT * FROM csv_data 
WHERE filename LIKE 'sales%';
```

## Connection Pooling

The application uses connection pooling for optimal performance:

- Max open connections: 25
- Max idle connections: 5
- Connection max lifetime: 5 minutes

These values can be adjusted in `internal/database/postgres.go` if needed.

## Error Handling

The application handles database errors gracefully:

- Connection failures are logged and the application exits
- Insert failures return HTTP 500 with error details
- Query failures return appropriate HTTP status codes

## Performance Considerations

### JSONB Indexing

For large datasets, consider adding GIN indexes:

```sql
CREATE INDEX idx_csv_data_gin ON csv_data USING GIN (data);
```

### Data Retention

Implement data cleanup for old records:

```sql
DELETE FROM csv_data 
WHERE created_at < NOW() - INTERVAL '90 days';
```

## Monitoring

### Check Database Size
```sql
SELECT pg_size_pretty(pg_database_size('csv2json')) as size;
```

### Count Records
```sql
SELECT COUNT(*) FROM csv_data;
```

### Recent Uploads
```sql
SELECT filename, created_at 
FROM csv_data 
ORDER BY created_at DESC 
LIMIT 10;
```

## Backup and Recovery

### Backup Database
```bash
pg_dump csv2json > backup.sql
```

### Restore Database
```bash
psql csv2json < backup.sql
```

## Troubleshooting

### Connection Issues

If you see "failed to connect to database" errors:

1. Verify PostgreSQL is running:
   ```bash
   # Windows
   Get-Service -Name postgresql*
   
   # Linux/Mac
   sudo service postgresql status
   ```

2. Check connection parameters in `.env`

3. Verify the database exists:
   ```sql
   \l
   ```

### Schema Issues

If tables are not created automatically:

1. Check logs for schema initialization errors
2. Manually run `docs/setup.sql`
3. Verify user permissions

### Performance Issues

If queries are slow:

1. Add appropriate indexes
2. Analyze query execution plans:
   ```sql
   EXPLAIN ANALYZE SELECT * FROM csv_data;
   ```
3. Consider partitioning for very large datasets

## Security Best Practices

1. **Never commit credentials**: Keep `.env` out of version control
2. **Use SSL in production**: Set `DB_SSLMODE=require`
3. **Limit user privileges**: Create a dedicated database user
4. **Regular backups**: Automate database backups
5. **Monitor access**: Enable PostgreSQL logging
