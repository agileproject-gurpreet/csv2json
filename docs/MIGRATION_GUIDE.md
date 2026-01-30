# Migration Checklist - Upgrading to PostgreSQL Integration

If you have an existing csv2json-api deployment without PostgreSQL, follow this checklist to upgrade.

## Pre-Migration

### 1. Backup Current System
- [ ] Take a full backup of your current deployment
- [ ] Document current API endpoints and usage
- [ ] Test rollback procedures

### 2. Review Requirements
- [ ] PostgreSQL 12+ available
- [ ] Database server accessible from application
- [ ] Sufficient disk space for data storage
- [ ] Network connectivity configured

### 3. Prepare Environment
- [ ] Install PostgreSQL if not present
- [ ] Create database user with appropriate privileges
- [ ] Configure firewall rules if necessary
- [ ] Test database connectivity

## Migration Steps

### Step 1: Update Code
```bash
# Pull latest changes
git pull origin main

# or download the updated files
# - internal/database/postgres.go (new)
# - cmd/api/main.go (updated)
# - internal/service/conversion_service.go (updated)
# - internal/handler/csv_handler.go (updated)
# - go.mod (updated)
```

### Step 2: Install Dependencies
```bash
go mod download
```

### Step 3: Create Database
```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE csv2json;

# Create user (optional, if not using postgres user)
CREATE USER csv2json_user WITH PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE csv2json TO csv2json_user;

# Exit
\q
```

### Step 4: Configure Environment
```bash
# Create .env file
cp .env.example .env
```

Edit `.env` with your settings:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=csv2json_user
DB_PASSWORD=secure_password
DB_NAME=csv2json
DB_SSLMODE=disable  # Use 'require' in production
PORT=8080
```

### Step 5: Test Build
```bash
go build ./cmd/api
```

### Step 6: Run Database Initialization
The application will automatically create tables on first run, or you can run manually:
```bash
psql -U csv2json_user -d csv2json -f docs/setup.sql
```

### Step 7: Test Locally
```bash
# Start the application
go run cmd/api/main.go

# In another terminal, test upload
curl -X POST http://localhost:8080/api/upload \
  -F "file=@examples/sample.csv"

# Test retrieval
curl http://localhost:8080/api/data
```

### Step 8: Deploy to Production
```bash
# Build production binary
go build -o csv2json-api ./cmd/api

# Stop existing service
sudo systemctl stop csv2json-api

# Replace binary
sudo cp csv2json-api /usr/local/bin/

# Update environment
sudo cp .env /etc/csv2json-api/.env

# Start service
sudo systemctl start csv2json-api

# Check status
sudo systemctl status csv2json-api
```

## Post-Migration Verification

### Checklist
- [ ] Application starts without errors
- [ ] Database connection successful
- [ ] Schema created correctly
- [ ] Upload endpoint works
- [ ] Data stored in database
- [ ] Retrieval endpoints work
- [ ] Health check responds
- [ ] Logs show no errors

### Verification Commands

```bash
# Check application logs
tail -f /var/log/csv2json-api/app.log

# Check PostgreSQL
psql -U csv2json_user -d csv2json

# List tables
\dt

# Check data
SELECT id, filename, created_at FROM csv_data ORDER BY created_at DESC LIMIT 5;

# Exit
\q
```

### API Endpoint Tests

```bash
# Health check
curl http://your-server:8080/api/health

# Upload test
curl -X POST http://your-server:8080/api/upload \
  -F "file=@test.csv"

# Get all data
curl http://your-server:8080/api/data

# Get by ID
curl http://your-server:8080/api/data/id?id=1
```

## Breaking Changes

### None! 
The migration is **backward compatible**:
- Existing API endpoints remain unchanged
- Response formats are the same
- No changes required to client code

### New Features
- ✅ Data persistence in PostgreSQL
- ✅ New endpoints: `/api/data` and `/api/data/id`
- ✅ Automatic schema management
- ✅ JSONB storage for efficient querying

## Rollback Plan

If issues occur, rollback steps:

### Option 1: Revert Code
```bash
# Checkout previous version
git checkout <previous-commit-hash>

# Rebuild
go build ./cmd/api

# Deploy old binary
sudo systemctl stop csv2json-api
sudo cp csv2json-api /usr/local/bin/
sudo systemctl start csv2json-api
```

### Option 2: Disable Database
If you want to keep the new code but disable database:
- Comment out database initialization in `main.go`
- Service will continue to work without persistence

## Troubleshooting

### Issue: Application won't start
**Symptoms**: 
- Error: "failed to connect to database"

**Solutions**:
1. Check PostgreSQL is running:
   ```bash
   sudo systemctl status postgresql
   ```
2. Verify credentials in `.env`
3. Test connection:
   ```bash
   psql -U csv2json_user -d csv2json
   ```

### Issue: Permission denied
**Symptoms**:
- Error: "failed to create schema"

**Solutions**:
```sql
-- Connect as postgres superuser
psql -U postgres

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE csv2json TO csv2json_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO csv2json_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO csv2json_user;
```

### Issue: Port conflicts
**Symptoms**:
- Error: "bind: address already in use"

**Solutions**:
1. Change port in `.env`:
   ```env
   PORT=8081
   ```
2. Or stop conflicting service

### Issue: Database size growing too large
**Solutions**:
1. Implement data retention policy:
   ```sql
   DELETE FROM csv_data WHERE created_at < NOW() - INTERVAL '30 days';
   ```
2. Set up automated cleanup cron job
3. Implement pagination in retrieval endpoints

## Performance Optimization

### After Migration

1. **Monitor Database Size**
   ```sql
   SELECT pg_size_pretty(pg_database_size('csv2json'));
   ```

2. **Add Additional Indexes** (if needed)
   ```sql
   CREATE INDEX idx_csv_data_gin ON csv_data USING GIN (data);
   ```

3. **Configure Connection Pool**
   Edit `internal/database/postgres.go`:
   ```go
   db.SetMaxOpenConns(50)  // Increase for high traffic
   db.SetMaxIdleConns(10)
   ```

4. **Enable Query Logging** (for debugging)
   ```sql
   ALTER DATABASE csv2json SET log_statement = 'all';
   ```

## Security Hardening

### Post-Migration Security Steps

1. **Enable SSL**
   ```env
   DB_SSLMODE=require
   ```

2. **Restrict Database Access**
   ```sql
   -- Revoke public access
   REVOKE ALL ON DATABASE csv2json FROM PUBLIC;
   
   -- Grant only to specific user
   GRANT CONNECT ON DATABASE csv2json TO csv2json_user;
   ```

3. **Use Strong Passwords**
   ```sql
   ALTER USER csv2json_user WITH PASSWORD 'very-strong-password-here';
   ```

4. **Configure pg_hba.conf**
   ```
   # Only allow local connections
   host    csv2json    csv2json_user    127.0.0.1/32    md5
   ```

5. **Implement API Authentication** (recommended)
   - Add API keys
   - Use JWT tokens
   - Implement rate limiting

## Monitoring

### Set Up Monitoring

1. **Application Metrics**
   - Monitor API response times
   - Track upload success/failure rates
   - Log database connection errors

2. **Database Metrics**
   ```sql
   -- Active connections
   SELECT count(*) FROM pg_stat_activity WHERE datname = 'csv2json';
   
   -- Table size
   SELECT pg_size_pretty(pg_total_relation_size('csv_data'));
   
   -- Row count
   SELECT COUNT(*) FROM csv_data;
   ```

3. **Set Up Alerts**
   - Database connection failures
   - Disk space warnings
   - High CPU usage
   - Memory consumption

## Support

If you encounter issues:

1. Check logs: Application and PostgreSQL
2. Review documentation in `docs/`
3. Test with minimal configuration
4. Verify prerequisites are met

---

## Success Criteria

Migration is successful when:
- ✅ Application starts without errors
- ✅ Database connection established
- ✅ CSV uploads store data in PostgreSQL
- ✅ Data retrieval works correctly
- ✅ All existing functionality preserved
- ✅ New endpoints accessible
- ✅ No data loss
- ✅ Performance acceptable

**Congratulations!** Your csv2json-api now has full PostgreSQL integration! 🎉
