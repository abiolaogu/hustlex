-- HustleX Pro Database Initialization
-- This script runs when PostgreSQL container first starts

-- Create additional databases
CREATE DATABASE hustlex_n8n;

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE hustlex TO hustlex;
GRANT ALL PRIVILEGES ON DATABASE hustlex_n8n TO hustlex;

-- Enable extensions in main database
\c hustlex
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
