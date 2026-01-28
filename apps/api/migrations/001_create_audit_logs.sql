-- Migration: Create Audit Logs Table
-- Description: SOC 2/PCI DSS compliant audit logging infrastructure
-- Author: HustleX Security Team
-- Date: 2024

-- ============================================================================
-- UP Migration
-- ============================================================================

-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create enum types for audit events
CREATE TYPE audit_event_type AS ENUM (
    'ACCESS',
    'DATA_CHANGE',
    'AUTHENTICATION',
    'AUTHORIZATION',
    'CONFIGURATION',
    'SECURITY_ALERT',
    'TRANSACTION'
);

CREATE TYPE audit_event_action AS ENUM (
    'C', -- Create
    'R', -- Read
    'U', -- Update
    'D', -- Delete
    'E'  -- Execute
);

CREATE TYPE audit_event_outcome AS ENUM (
    'SUCCESS',
    'FAILURE',
    'ERROR'
);

-- Create the main audit_logs table
CREATE TABLE audit_logs (
    -- Primary Key
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    -- Event Identification
    event_type audit_event_type NOT NULL,
    event_action audit_event_action NOT NULL,
    event_outcome audit_event_outcome NOT NULL,

    -- Timestamps (RFC 3339 format)
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Service/Application Context
    service VARCHAR(100) NOT NULL,
    component VARCHAR(100),
    environment VARCHAR(50),
    version VARCHAR(50),

    -- Actor Information (Who performed the action)
    actor_user_id VARCHAR(100),
    actor_username VARCHAR(255),
    actor_role VARCHAR(100),
    actor_ip_address INET,
    actor_user_agent TEXT,
    actor_session_id VARCHAR(255),

    -- Target Information (What was affected)
    target_type VARCHAR(100),
    target_id VARCHAR(255),
    target_name VARCHAR(255),

    -- Request Context
    request_id VARCHAR(100),
    correlation_id VARCHAR(100),

    -- Event Details
    message TEXT NOT NULL,
    error_code VARCHAR(50),
    error_message TEXT,

    -- Change Tracking (for DATA_CHANGE events)
    old_value JSONB,
    new_value JSONB,

    -- Additional Metadata
    metadata JSONB DEFAULT '{}'::jsonb,

    -- Retention Management
    retention_until TIMESTAMPTZ,

    -- Audit Trail Protection (immutability)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    checksum VARCHAR(64) -- SHA-256 of row for integrity verification
);

-- ============================================================================
-- Indexes for Query Performance
-- ============================================================================

-- Time-based queries (most common pattern)
CREATE INDEX idx_audit_logs_timestamp ON audit_logs (timestamp DESC);

-- User activity queries
CREATE INDEX idx_audit_logs_actor_user_id ON audit_logs (actor_user_id) WHERE actor_user_id IS NOT NULL;

-- Event type filtering
CREATE INDEX idx_audit_logs_event_type ON audit_logs (event_type);

-- Outcome filtering (for security review)
CREATE INDEX idx_audit_logs_event_outcome ON audit_logs (event_outcome);

-- Target queries
CREATE INDEX idx_audit_logs_target ON audit_logs (target_type, target_id) WHERE target_type IS NOT NULL;

-- Request tracing
CREATE INDEX idx_audit_logs_request_id ON audit_logs (request_id) WHERE request_id IS NOT NULL;
CREATE INDEX idx_audit_logs_correlation_id ON audit_logs (correlation_id) WHERE correlation_id IS NOT NULL;

-- Security incident investigation
CREATE INDEX idx_audit_logs_ip_address ON audit_logs (actor_ip_address) WHERE actor_ip_address IS NOT NULL;

-- Composite index for common query patterns
CREATE INDEX idx_audit_logs_user_time ON audit_logs (actor_user_id, timestamp DESC) WHERE actor_user_id IS NOT NULL;
CREATE INDEX idx_audit_logs_type_time ON audit_logs (event_type, timestamp DESC);

-- Retention management
CREATE INDEX idx_audit_logs_retention ON audit_logs (retention_until) WHERE retention_until IS NOT NULL;

-- JSONB metadata index for flexible queries
CREATE INDEX idx_audit_logs_metadata ON audit_logs USING GIN (metadata);

-- ============================================================================
-- Partitioning for Large-Scale Deployments (Optional)
-- ============================================================================

-- If you need partitioning for high-volume logging, uncomment and modify:
--
-- CREATE TABLE audit_logs (
--     ...
-- ) PARTITION BY RANGE (timestamp);
--
-- CREATE TABLE audit_logs_2024_q1 PARTITION OF audit_logs
--     FOR VALUES FROM ('2024-01-01') TO ('2024-04-01');
-- CREATE TABLE audit_logs_2024_q2 PARTITION OF audit_logs
--     FOR VALUES FROM ('2024-04-01') TO ('2024-07-01');
-- etc.

-- ============================================================================
-- Row-Level Security (Optional - for multi-tenant scenarios)
-- ============================================================================

-- Enable RLS
ALTER TABLE audit_logs ENABLE ROW LEVEL SECURITY;

-- Policy for security team (full access)
CREATE POLICY audit_logs_security_team ON audit_logs
    FOR ALL
    TO security_team
    USING (true);

-- Policy for regular users (own events only)
CREATE POLICY audit_logs_user_own ON audit_logs
    FOR SELECT
    TO authenticated
    USING (actor_user_id = current_setting('app.current_user_id', true));

-- ============================================================================
-- Triggers for Data Integrity
-- ============================================================================

-- Function to calculate row checksum
CREATE OR REPLACE FUNCTION calculate_audit_checksum()
RETURNS TRIGGER AS $$
BEGIN
    NEW.checksum = encode(
        sha256(
            (NEW.id::text || NEW.event_type::text || NEW.event_action::text ||
             NEW.event_outcome::text || NEW.timestamp::text || NEW.service ||
             COALESCE(NEW.actor_user_id, '') || COALESCE(NEW.target_id, '') ||
             NEW.message || COALESCE(NEW.metadata::text, '{}'))::bytea
        ),
        'hex'
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to set checksum on insert
CREATE TRIGGER audit_logs_checksum_trigger
    BEFORE INSERT ON audit_logs
    FOR EACH ROW
    EXECUTE FUNCTION calculate_audit_checksum();

-- Function to prevent updates (immutability)
CREATE OR REPLACE FUNCTION prevent_audit_update()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'Audit logs cannot be modified';
END;
$$ LANGUAGE plpgsql;

-- Trigger to prevent updates
CREATE TRIGGER audit_logs_immutable_trigger
    BEFORE UPDATE ON audit_logs
    FOR EACH ROW
    EXECUTE FUNCTION prevent_audit_update();

-- Function to prevent deletes (except for retention policy)
CREATE OR REPLACE FUNCTION prevent_audit_delete()
RETURNS TRIGGER AS $$
BEGIN
    -- Only allow delete if called from retention cleanup job
    IF current_setting('app.retention_cleanup', true) = 'true' THEN
        RETURN OLD;
    END IF;
    RAISE EXCEPTION 'Audit logs cannot be deleted manually';
END;
$$ LANGUAGE plpgsql;

-- Trigger to prevent manual deletes
CREATE TRIGGER audit_logs_no_delete_trigger
    BEFORE DELETE ON audit_logs
    FOR EACH ROW
    EXECUTE FUNCTION prevent_audit_delete();

-- ============================================================================
-- Helper Functions
-- ============================================================================

-- Function to query audit logs with pagination
CREATE OR REPLACE FUNCTION query_audit_logs(
    p_actor_user_id VARCHAR DEFAULT NULL,
    p_event_type audit_event_type DEFAULT NULL,
    p_event_outcome audit_event_outcome DEFAULT NULL,
    p_target_type VARCHAR DEFAULT NULL,
    p_target_id VARCHAR DEFAULT NULL,
    p_start_time TIMESTAMPTZ DEFAULT NULL,
    p_end_time TIMESTAMPTZ DEFAULT NULL,
    p_limit INT DEFAULT 100,
    p_offset INT DEFAULT 0
)
RETURNS TABLE (
    id UUID,
    event_type audit_event_type,
    event_action audit_event_action,
    event_outcome audit_event_outcome,
    timestamp TIMESTAMPTZ,
    service VARCHAR,
    actor_user_id VARCHAR,
    target_type VARCHAR,
    target_id VARCHAR,
    message TEXT,
    metadata JSONB,
    total_count BIGINT
) AS $$
DECLARE
    total BIGINT;
BEGIN
    -- Get total count
    SELECT COUNT(*) INTO total
    FROM audit_logs al
    WHERE (p_actor_user_id IS NULL OR al.actor_user_id = p_actor_user_id)
      AND (p_event_type IS NULL OR al.event_type = p_event_type)
      AND (p_event_outcome IS NULL OR al.event_outcome = p_event_outcome)
      AND (p_target_type IS NULL OR al.target_type = p_target_type)
      AND (p_target_id IS NULL OR al.target_id = p_target_id)
      AND (p_start_time IS NULL OR al.timestamp >= p_start_time)
      AND (p_end_time IS NULL OR al.timestamp <= p_end_time);

    -- Return results
    RETURN QUERY
    SELECT
        al.id,
        al.event_type,
        al.event_action,
        al.event_outcome,
        al.timestamp,
        al.service,
        al.actor_user_id,
        al.target_type,
        al.target_id,
        al.message,
        al.metadata,
        total
    FROM audit_logs al
    WHERE (p_actor_user_id IS NULL OR al.actor_user_id = p_actor_user_id)
      AND (p_event_type IS NULL OR al.event_type = p_event_type)
      AND (p_event_outcome IS NULL OR al.event_outcome = p_event_outcome)
      AND (p_target_type IS NULL OR al.target_type = p_target_type)
      AND (p_target_id IS NULL OR al.target_id = p_target_id)
      AND (p_start_time IS NULL OR al.timestamp >= p_start_time)
      AND (p_end_time IS NULL OR al.timestamp <= p_end_time)
    ORDER BY al.timestamp DESC
    LIMIT p_limit
    OFFSET p_offset;
END;
$$ LANGUAGE plpgsql;

-- Function for retention cleanup
CREATE OR REPLACE FUNCTION cleanup_expired_audit_logs()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    -- Set flag to allow deletion
    PERFORM set_config('app.retention_cleanup', 'true', true);

    -- Delete expired logs
    DELETE FROM audit_logs
    WHERE retention_until IS NOT NULL
      AND retention_until < NOW();

    GET DIAGNOSTICS deleted_count = ROW_COUNT;

    -- Reset flag
    PERFORM set_config('app.retention_cleanup', 'false', true);

    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- Comments for Documentation
-- ============================================================================

COMMENT ON TABLE audit_logs IS 'SOC 2/PCI DSS compliant audit log table for security and compliance tracking';

COMMENT ON COLUMN audit_logs.event_type IS 'Type of audit event (ACCESS, DATA_CHANGE, AUTHENTICATION, etc.)';
COMMENT ON COLUMN audit_logs.event_action IS 'CRUD action: C=Create, R=Read, U=Update, D=Delete, E=Execute';
COMMENT ON COLUMN audit_logs.event_outcome IS 'Result of the action: SUCCESS, FAILURE, ERROR';
COMMENT ON COLUMN audit_logs.actor_user_id IS 'User ID who performed the action';
COMMENT ON COLUMN audit_logs.actor_ip_address IS 'IP address of the actor';
COMMENT ON COLUMN audit_logs.target_id IS 'ID of the resource that was affected';
COMMENT ON COLUMN audit_logs.old_value IS 'Previous value for DATA_CHANGE events (JSONB)';
COMMENT ON COLUMN audit_logs.new_value IS 'New value for DATA_CHANGE events (JSONB)';
COMMENT ON COLUMN audit_logs.checksum IS 'SHA-256 checksum for integrity verification';
COMMENT ON COLUMN audit_logs.retention_until IS 'Date until which the log must be retained';

-- ============================================================================
-- DOWN Migration
-- ============================================================================

-- To rollback, run:
-- DROP TABLE IF EXISTS audit_logs CASCADE;
-- DROP TYPE IF EXISTS audit_event_type;
-- DROP TYPE IF EXISTS audit_event_action;
-- DROP TYPE IF EXISTS audit_event_outcome;
-- DROP FUNCTION IF EXISTS calculate_audit_checksum();
-- DROP FUNCTION IF EXISTS prevent_audit_update();
-- DROP FUNCTION IF EXISTS prevent_audit_delete();
-- DROP FUNCTION IF EXISTS query_audit_logs();
-- DROP FUNCTION IF EXISTS cleanup_expired_audit_logs();
