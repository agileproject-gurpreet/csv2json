package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	DB *sql.DB
}

// Config holds database configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB(config Config) (*PostgresDB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresDB{DB: db}, nil
}

// InitSchema creates the necessary tables if they don't exist
func (p *PostgresDB) InitSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS csv_data (
		id SERIAL PRIMARY KEY,
		filename VARCHAR(255),
		data JSONB NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE INDEX IF NOT EXISTS idx_csv_data_created_at ON csv_data(created_at);
	CREATE INDEX IF NOT EXISTS idx_csv_data_filename ON csv_data(filename);
	`

	_, err := p.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	return nil
}

// InsertCSVData inserts CSV data (as JSON) into the database
func (p *PostgresDB) InsertCSVData(filename string, records []map[string]string) error {
	// Convert records to JSON
	jsonData, err := json.Marshal(records)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	query := `
		INSERT INTO csv_data (filename, data)
		VALUES ($1, $2)
		RETURNING id
	`

	var id int
	err = p.DB.QueryRow(query, filename, jsonData).Scan(&id)
	if err != nil {
		return fmt.Errorf("failed to insert data: %w", err)
	}

	return nil
}

// GetAllCSVData retrieves all CSV data from the database
func (p *PostgresDB) GetAllCSVData() ([]map[string]interface{}, error) {
	query := `
		SELECT id, filename, data, created_at
		FROM csv_data
		ORDER BY created_at DESC
	`

	rows, err := p.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query data: %w", err)
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id int
		var filename string
		var data []byte
		var createdAt time.Time

		if err := rows.Scan(&id, &filename, &data, &createdAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		var jsonData interface{}
		if err := json.Unmarshal(data, &jsonData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal data: %w", err)
		}

		results = append(results, map[string]interface{}{
			"id":         id,
			"filename":   filename,
			"data":       jsonData,
			"created_at": createdAt,
		})
	}

	return results, nil
}

// GetCSVDataByID retrieves CSV data by ID
func (p *PostgresDB) GetCSVDataByID(id int) (map[string]interface{}, error) {
	query := `
		SELECT id, filename, data, created_at
		FROM csv_data
		WHERE id = $1
	`

	var filename string
	var data []byte
	var createdAt time.Time

	err := p.DB.QueryRow(query, id).Scan(&id, &filename, &data, &createdAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("record not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query data: %w", err)
	}

	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return map[string]interface{}{
		"id":         id,
		"filename":   filename,
		"data":       jsonData,
		"created_at": createdAt,
	}, nil
}

// Close closes the database connection
func (p *PostgresDB) Close() error {
	return p.DB.Close()
}
