-- PostgreSQL setup script for csv2json-api

-- Create database (run as postgres superuser)
CREATE DATABASE csv2json;

-- Connect to the database
\c csv2json;

-- Create csv_data table
CREATE TABLE IF NOT EXISTS csv_data (
    id SERIAL PRIMARY KEY,
    filename VARCHAR(255),
    data JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_csv_data_created_at ON csv_data(created_at);
CREATE INDEX IF NOT EXISTS idx_csv_data_filename ON csv_data(filename);

-- Grant privileges (adjust username as needed)
-- GRANT ALL PRIVILEGES ON DATABASE csv2json TO your_username;
-- GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO your_username;

-- Verify table creation
\dt

-- Sample query to view all data
-- SELECT id, filename, created_at FROM csv_data ORDER BY created_at DESC;
