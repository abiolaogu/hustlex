-- =============================================================================
-- HustleX Database Initialization Script
-- =============================================================================
-- This script runs when PostgreSQL container starts for the first time

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";      -- UUID generation
CREATE EXTENSION IF NOT EXISTS "pg_trgm";        -- Trigram matching for search
CREATE EXTENSION IF NOT EXISTS "btree_gin";      -- GIN index support
CREATE EXTENSION IF NOT EXISTS "btree_gist";     -- GIST index support

-- Create additional schemas for organization
CREATE SCHEMA IF NOT EXISTS audit;
CREATE SCHEMA IF NOT EXISTS analytics;

-- Grant permissions
GRANT ALL PRIVILEGES ON SCHEMA audit TO hustlex;
GRANT ALL PRIVILEGES ON SCHEMA analytics TO hustlex;

-- Create audit log table for compliance
CREATE TABLE IF NOT EXISTS audit.activity_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID,
    action VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID,
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Index for audit queries
CREATE INDEX IF NOT EXISTS idx_audit_user_id ON audit.activity_log(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_action ON audit.activity_log(action);
CREATE INDEX IF NOT EXISTS idx_audit_entity ON audit.activity_log(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_audit_created_at ON audit.activity_log(created_at DESC);

-- Create analytics tables for reporting
CREATE TABLE IF NOT EXISTS analytics.daily_metrics (
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL UNIQUE,
    new_users INTEGER DEFAULT 0,
    active_users INTEGER DEFAULT 0,
    total_transactions INTEGER DEFAULT 0,
    transaction_volume BIGINT DEFAULT 0,  -- in kobo
    new_gigs INTEGER DEFAULT 0,
    completed_gigs INTEGER DEFAULT 0,
    gig_volume BIGINT DEFAULT 0,          -- in kobo
    new_circles INTEGER DEFAULT 0,
    contributions_count INTEGER DEFAULT 0,
    contributions_volume BIGINT DEFAULT 0, -- in kobo
    loans_disbursed INTEGER DEFAULT 0,
    loans_volume BIGINT DEFAULT 0,         -- in kobo
    loans_repaid INTEGER DEFAULT 0,
    repayment_volume BIGINT DEFAULT 0,     -- in kobo
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply trigger to analytics tables
DROP TRIGGER IF EXISTS update_daily_metrics_updated_at ON analytics.daily_metrics;
CREATE TRIGGER update_daily_metrics_updated_at
    BEFORE UPDATE ON analytics.daily_metrics
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Create function for audit logging (can be called from application)
CREATE OR REPLACE FUNCTION audit.log_activity(
    p_user_id UUID,
    p_action VARCHAR(50),
    p_entity_type VARCHAR(50),
    p_entity_id UUID,
    p_old_values JSONB DEFAULT NULL,
    p_new_values JSONB DEFAULT NULL,
    p_ip_address INET DEFAULT NULL,
    p_user_agent TEXT DEFAULT NULL
)
RETURNS UUID AS $$
DECLARE
    log_id UUID;
BEGIN
    INSERT INTO audit.activity_log (
        user_id, action, entity_type, entity_id,
        old_values, new_values, ip_address, user_agent
    ) VALUES (
        p_user_id, p_action, p_entity_type, p_entity_id,
        p_old_values, p_new_values, p_ip_address, p_user_agent
    ) RETURNING id INTO log_id;
    
    RETURN log_id;
END;
$$ LANGUAGE plpgsql;

-- Create function to get user transaction summary
CREATE OR REPLACE FUNCTION get_user_transaction_summary(p_user_id UUID)
RETURNS TABLE (
    total_deposits BIGINT,
    total_withdrawals BIGINT,
    total_transfers_sent BIGINT,
    total_transfers_received BIGINT,
    gig_earnings BIGINT,
    gig_spending BIGINT,
    savings_contributions BIGINT,
    loan_received BIGINT,
    loan_repaid BIGINT
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        COALESCE(SUM(CASE WHEN type = 'deposit' AND status = 'completed' THEN amount ELSE 0 END), 0) as total_deposits,
        COALESCE(SUM(CASE WHEN type = 'withdrawal' AND status = 'completed' THEN amount ELSE 0 END), 0) as total_withdrawals,
        COALESCE(SUM(CASE WHEN type = 'transfer_out' AND status = 'completed' THEN amount ELSE 0 END), 0) as total_transfers_sent,
        COALESCE(SUM(CASE WHEN type = 'transfer_in' AND status = 'completed' THEN amount ELSE 0 END), 0) as total_transfers_received,
        COALESCE(SUM(CASE WHEN type = 'gig_payment' AND status = 'completed' THEN amount ELSE 0 END), 0) as gig_earnings,
        COALESCE(SUM(CASE WHEN type = 'escrow_lock' AND status = 'completed' THEN amount ELSE 0 END), 0) as gig_spending,
        COALESCE(SUM(CASE WHEN type = 'savings_contribution' AND status = 'completed' THEN amount ELSE 0 END), 0) as savings_contributions,
        COALESCE(SUM(CASE WHEN type = 'loan_disbursement' AND status = 'completed' THEN amount ELSE 0 END), 0) as loan_received,
        COALESCE(SUM(CASE WHEN type = 'loan_repayment' AND status = 'completed' THEN amount ELSE 0 END), 0) as loan_repaid
    FROM transactions
    WHERE user_id = p_user_id;
END;
$$ LANGUAGE plpgsql;

-- Print success message
DO $$
BEGIN
    RAISE NOTICE 'HustleX database initialized successfully!';
END $$;
