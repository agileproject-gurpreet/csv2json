# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned
- Support for custom CSV delimiters
- Batch upload functionality
- Data export endpoints (JSON to CSV)
- Query filtering for stored data

## [1.0.0] - 2026-01-28

### Added
- Initial release of csv2json converter
- REST API with CSV upload endpoint (`POST /api/upload`)
- Automatic CSV to JSON conversion with header detection
- PostgreSQL integration for data persistence
- JSONB storage for efficient querying
- Database connection pooling
- Health check endpoint (`GET /api/health`)
- Retrieve all stored data endpoint (`GET /api/data`)
- Retrieve data by ID endpoint (`GET /api/data/id?id={id}`)
- Automatic table creation on startup
- CSV parser with comprehensive error handling
- Support for various CSV formats
- Unit tests for CSV parser
- Complete API documentation
- Environment variable configuration
- Sample CSV file for testing
- API usage examples

### Database
- Created `csv_data` table with JSONB support
- Added automatic timestamp tracking
- Implemented indexed ID for fast lookups
- Setup scripts for database initialization

### Dependencies
- Go 1.25.4 support
- PostgreSQL driver (lib/pq v1.10.9)
- PostgreSQL 12+ compatibility

### Documentation
- README with comprehensive setup instructions
- API endpoint documentation with examples
- Database schema documentation
- Environment variable reference table
- Development and testing guidelines
- LICENSE file (MIT License)

### Project Structure
- Organized code with `cmd/`, `internal/`, and `pkg/` structure
- Separation of concerns: handlers, services, parsers, database
- Examples directory with sample files
- Docs directory with SQL setup scripts

### Security
- Parameterized SQL queries to prevent injection
- Environment-based configuration
- Configurable SSL mode for database connections

[unreleased]: https://github.com/lavishag4193/csv2jsonx/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/lavishag4193/csv2jsonx/releases/tag/v1.0.0
