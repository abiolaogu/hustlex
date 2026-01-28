-- HustleX Pro Database Schema
-- Multi-tenant platform for gig marketplace, fintech, and diaspora services

-- =============================================================================
-- EXTENSIONS
-- =============================================================================
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- =============================================================================
-- ENUMS
-- =============================================================================

-- User related
CREATE TYPE user_status AS ENUM ('active', 'inactive', 'suspended', 'pending_verification');
CREATE TYPE user_role AS ENUM ('consumer', 'provider', 'admin', 'super_admin');
CREATE TYPE verification_status AS ENUM ('unverified', 'pending', 'verified', 'rejected');
CREATE TYPE kyc_level AS ENUM ('none', 'basic', 'intermediate', 'full');

-- Wallet related
CREATE TYPE wallet_type AS ENUM ('main', 'savings', 'escrow', 'multi_currency');
CREATE TYPE wallet_status AS ENUM ('active', 'frozen', 'closed');
CREATE TYPE transaction_type AS ENUM ('credit', 'debit', 'transfer', 'payment', 'refund', 'withdrawal', 'deposit', 'fee', 'fx_conversion');
CREATE TYPE transaction_status AS ENUM ('pending', 'processing', 'completed', 'failed', 'reversed', 'cancelled');
CREATE TYPE currency_code AS ENUM ('NGN', 'GBP', 'USD', 'EUR', 'CAD', 'GHS', 'KES');

-- Service related
CREATE TYPE service_status AS ENUM ('draft', 'active', 'paused', 'archived');
CREATE TYPE booking_status AS ENUM ('pending', 'quoted', 'confirmed', 'paid', 'processing', 'in_progress', 'completed', 'cancelled', 'refunded', 'failed', 'disputed');
CREATE TYPE payment_method AS ENUM ('wallet', 'card', 'bank_transfer', 'mobile_money', 'ussd');

-- Remittance related
CREATE TYPE remittance_status AS ENUM ('pending', 'quoted', 'initiated', 'processing', 'in_transit', 'delivered', 'completed', 'failed', 'cancelled', 'refunded', 'on_hold');
CREATE TYPE remittance_purpose AS ENUM ('family_support', 'education', 'medical', 'investment', 'property_purchase', 'business', 'gift', 'salary', 'other');
CREATE TYPE delivery_method AS ENUM ('bank_transfer', 'mobile_wallet', 'cash_pickup', 'home_delivery', 'hustlex_wallet');
CREATE TYPE recurrence_type AS ENUM ('none', 'weekly', 'bi_weekly', 'monthly', 'quarterly');

-- Beneficiary related
CREATE TYPE beneficiary_type AS ENUM ('family', 'friend', 'business', 'self', 'charity', 'other');
CREATE TYPE beneficiary_status AS ENUM ('active', 'inactive', 'pending_verification', 'blocked');

-- Savings related
CREATE TYPE savings_circle_status AS ENUM ('forming', 'active', 'paused', 'completed', 'dissolved');
CREATE TYPE circle_member_status AS ENUM ('invited', 'active', 'defaulted', 'withdrawn', 'completed');

-- Notification related
CREATE TYPE notification_type AS ENUM ('push', 'sms', 'email', 'in_app');
CREATE TYPE notification_status AS ENUM ('pending', 'sent', 'delivered', 'failed', 'read');

-- =============================================================================
-- USERS & PROFILES
-- =============================================================================

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE,
    phone VARCHAR(20) UNIQUE,
    password_hash VARCHAR(255),
    status user_status NOT NULL DEFAULT 'pending_verification',
    role user_role NOT NULL DEFAULT 'consumer',
    email_verified BOOLEAN DEFAULT FALSE,
    phone_verified BOOLEAN DEFAULT FALSE,
    last_login_at TIMESTAMPTZ,
    failed_login_attempts INT DEFAULT 0,
    locked_until TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100),
    display_name VARCHAR(100),
    avatar_url TEXT,
    date_of_birth DATE,
    gender VARCHAR(20),

    -- Location
    address_line1 VARCHAR(255),
    address_line2 VARCHAR(255),
    city VARCHAR(100),
    state VARCHAR(100),
    country VARCHAR(100) DEFAULT 'Nigeria',
    postal_code VARCHAR(20),

    -- Diaspora specific
    country_of_residence VARCHAR(100),
    diaspora_region VARCHAR(50), -- 'UK', 'US', 'EU', 'CA', etc.

    -- KYC
    kyc_level kyc_level DEFAULT 'none',
    bvn_verified BOOLEAN DEFAULT FALSE,
    nin_verified BOOLEAN DEFAULT FALSE,
    verification_status verification_status DEFAULT 'unverified',

    -- Settings
    preferred_currency currency_code DEFAULT 'NGN',
    language VARCHAR(10) DEFAULT 'en',
    timezone VARCHAR(50) DEFAULT 'Africa/Lagos',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =============================================================================
-- WALLETS
-- =============================================================================

CREATE TABLE wallets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    wallet_type wallet_type NOT NULL DEFAULT 'main',
    currency currency_code NOT NULL DEFAULT 'NGN',
    status wallet_status NOT NULL DEFAULT 'active',

    available_balance DECIMAL(20, 4) NOT NULL DEFAULT 0,
    pending_balance DECIMAL(20, 4) NOT NULL DEFAULT 0,
    locked_balance DECIMAL(20, 4) NOT NULL DEFAULT 0,

    -- Limits
    daily_transaction_limit DECIMAL(20, 4) DEFAULT 500000,
    monthly_transaction_limit DECIMAL(20, 4) DEFAULT 5000000,

    -- Security
    pin_hash VARCHAR(255),
    pin_attempts INT DEFAULT 0,
    pin_locked_until TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(user_id, wallet_type, currency)
);

CREATE TABLE multi_currency_wallets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    base_wallet_id UUID NOT NULL REFERENCES wallets(id),
    primary_currency currency_code NOT NULL DEFAULT 'NGN',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE currency_balances (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    multi_currency_wallet_id UUID NOT NULL REFERENCES multi_currency_wallets(id) ON DELETE CASCADE,
    currency currency_code NOT NULL,
    available_balance DECIMAL(20, 4) NOT NULL DEFAULT 0,
    pending_balance DECIMAL(20, 4) NOT NULL DEFAULT 0,
    locked_balance DECIMAL(20, 4) NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(multi_currency_wallet_id, currency)
);

-- =============================================================================
-- TRANSACTIONS
-- =============================================================================

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    wallet_id UUID NOT NULL REFERENCES wallets(id),
    reference VARCHAR(50) UNIQUE NOT NULL,
    external_reference VARCHAR(100),

    type transaction_type NOT NULL,
    status transaction_status NOT NULL DEFAULT 'pending',

    amount DECIMAL(20, 4) NOT NULL,
    fee DECIMAL(20, 4) DEFAULT 0,
    currency currency_code NOT NULL DEFAULT 'NGN',

    -- For transfers
    source_wallet_id UUID REFERENCES wallets(id),
    destination_wallet_id UUID REFERENCES wallets(id),

    -- Balance tracking
    balance_before DECIMAL(20, 4),
    balance_after DECIMAL(20, 4),

    description TEXT,
    metadata JSONB DEFAULT '{}',

    -- Timestamps
    initiated_at TIMESTAMPTZ DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    failed_at TIMESTAMPTZ,
    failure_reason TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =============================================================================
-- SERVICES & BOOKINGS
-- =============================================================================

CREATE TABLE service_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    icon VARCHAR(50),
    parent_id UUID REFERENCES service_categories(id),
    is_active BOOLEAN DEFAULT TRUE,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE services (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    provider_id UUID NOT NULL REFERENCES users(id),
    category_id UUID NOT NULL REFERENCES service_categories(id),

    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255),
    description TEXT NOT NULL,

    -- Pricing
    base_price DECIMAL(20, 4) NOT NULL,
    currency currency_code DEFAULT 'NGN',
    pricing_type VARCHAR(20) DEFAULT 'fixed', -- 'fixed', 'hourly', 'negotiable'

    -- Location
    service_area TEXT[], -- Array of supported areas
    is_remote_available BOOLEAN DEFAULT FALSE,

    -- Media
    images TEXT[],
    video_url TEXT,

    -- Status & Visibility
    status service_status DEFAULT 'draft',
    is_featured BOOLEAN DEFAULT FALSE,

    -- Stats
    view_count INT DEFAULT 0,
    booking_count INT DEFAULT 0,
    average_rating DECIMAL(3, 2) DEFAULT 0,
    review_count INT DEFAULT 0,

    -- Tags
    tags TEXT[],

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    reference VARCHAR(50) UNIQUE NOT NULL,

    consumer_id UUID NOT NULL REFERENCES users(id),
    provider_id UUID NOT NULL REFERENCES users(id),
    service_id UUID REFERENCES services(id),

    status booking_status NOT NULL DEFAULT 'pending',

    -- Amounts
    service_amount DECIMAL(20, 4) NOT NULL,
    platform_fee DECIMAL(20, 4) DEFAULT 0,
    provider_fee DECIMAL(20, 4) DEFAULT 0,
    total_amount DECIMAL(20, 4) NOT NULL,
    currency currency_code DEFAULT 'NGN',

    -- Schedule
    scheduled_date DATE,
    scheduled_time_start TIME,
    scheduled_time_end TIME,

    -- Location
    service_address TEXT,
    service_city VARCHAR(100),
    service_state VARCHAR(100),
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),

    -- Notes
    consumer_notes TEXT,
    provider_notes TEXT,

    -- Payment
    payment_method payment_method,
    payment_reference VARCHAR(100),
    paid_at TIMESTAMPTZ,

    -- Escrow
    escrow_wallet_id UUID REFERENCES wallets(id),
    escrow_released_at TIMESTAMPTZ,

    -- Completion
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    cancellation_reason TEXT,
    cancelled_by UUID REFERENCES users(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =============================================================================
-- DIASPORA - BENEFICIARIES
-- =============================================================================

CREATE TABLE beneficiaries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),

    -- Basic Info
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100),
    nickname VARCHAR(50),
    relationship beneficiary_type NOT NULL,
    relationship_description VARCHAR(100),

    -- Contact
    phone_primary VARCHAR(20) NOT NULL,
    phone_secondary VARCHAR(20),
    email VARCHAR(255),

    -- Location
    address_line1 VARCHAR(255),
    address_line2 VARCHAR(255),
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL,
    country VARCHAR(100) NOT NULL DEFAULT 'Nigeria',
    postal_code VARCHAR(20),

    -- Bank Details
    bank_name VARCHAR(100),
    bank_code VARCHAR(20),
    account_number VARCHAR(20),
    account_name VARCHAR(255),

    -- Mobile Wallet
    mobile_wallet_provider VARCHAR(50),
    mobile_wallet_number VARCHAR(20),

    -- Preferred Delivery
    preferred_delivery_method delivery_method DEFAULT 'bank_transfer',
    preferred_currency currency_code DEFAULT 'NGN',

    -- Status
    status beneficiary_status DEFAULT 'active',
    verification_status verification_status DEFAULT 'unverified',
    verified_at TIMESTAMPTZ,

    -- Usage
    is_favorite BOOLEAN DEFAULT FALSE,
    last_transfer_at TIMESTAMPTZ,
    transfer_count INT DEFAULT 0,
    total_transferred DECIMAL(20, 4) DEFAULT 0,

    -- Notes
    notes TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =============================================================================
-- DIASPORA - REMITTANCES
-- =============================================================================

CREATE TABLE remittances (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    beneficiary_id UUID NOT NULL REFERENCES beneficiaries(id),
    source_wallet_id UUID REFERENCES wallets(id),

    reference VARCHAR(50) UNIQUE NOT NULL,
    external_reference VARCHAR(100),
    status remittance_status NOT NULL DEFAULT 'pending',
    status_message TEXT,

    -- Currencies & Amounts
    source_currency currency_code NOT NULL,
    target_currency currency_code NOT NULL,
    source_amount DECIMAL(20, 4) NOT NULL,
    target_amount DECIMAL(20, 4) NOT NULL,
    fx_rate DECIMAL(20, 8) NOT NULL,
    fx_quote_id VARCHAR(50),

    -- Fees
    transfer_fee DECIMAL(20, 4) DEFAULT 0,
    fx_fee DECIMAL(20, 4) DEFAULT 0,
    total_fee DECIMAL(20, 4) DEFAULT 0,
    total_source_amount DECIMAL(20, 4) NOT NULL,

    -- Purpose & Delivery
    purpose remittance_purpose NOT NULL,
    purpose_description TEXT,
    delivery_method delivery_method NOT NULL,

    -- Recurring
    is_recurring BOOLEAN DEFAULT FALSE,
    recurrence_type recurrence_type DEFAULT 'none',
    recurrence_start_date DATE,
    recurrence_end_date DATE,
    next_recurrence_date DATE,
    parent_remittance_id UUID REFERENCES remittances(id),

    -- Tracking
    estimated_delivery TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,

    -- Payment
    payment_method VARCHAR(50),
    payment_reference VARCHAR(100),
    payment_provider VARCHAR(50),
    paid_at TIMESTAMPTZ,

    -- Compliance
    compliance_status VARCHAR(50) DEFAULT 'pending',
    compliance_notes TEXT,
    aml_checked BOOLEAN DEFAULT FALSE,
    aml_checked_at TIMESTAMPTZ,

    -- Notifications
    sender_notified BOOLEAN DEFAULT FALSE,
    beneficiary_notified BOOLEAN DEFAULT FALSE,

    -- Cancellation
    cancelled_at TIMESTAMPTZ,
    cancellation_reason TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =============================================================================
-- DIASPORA - BOOKINGS
-- =============================================================================

CREATE TABLE diaspora_bookings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    beneficiary_id UUID REFERENCES beneficiaries(id),
    service_provider_id UUID REFERENCES users(id),

    type VARCHAR(50) NOT NULL, -- 'service', 'remittance', 'bill_payment', 'airtime', 'gift_card'
    status booking_status NOT NULL DEFAULT 'pending',
    reference VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,

    -- Currency & Amounts
    source_currency currency_code NOT NULL,
    target_currency currency_code NOT NULL,
    source_amount DECIMAL(20, 4) NOT NULL,
    target_amount DECIMAL(20, 4) NOT NULL,
    fee DECIMAL(20, 4) DEFAULT 0,
    fee_currency currency_code,
    total_source_amount DECIMAL(20, 4) NOT NULL,

    -- FX Rate Locking
    fx_rate DECIMAL(20, 8),
    fx_quote_id VARCHAR(50),
    fx_rate_locked_at TIMESTAMPTZ,
    fx_rate_expires_at TIMESTAMPTZ,
    fx_rate_locked BOOLEAN DEFAULT FALSE,

    -- Service Details
    service_date DATE,
    service_address TEXT,
    service_city VARCHAR(100),
    service_state VARCHAR(100),
    service_notes TEXT,

    -- Verification
    requires_verification BOOLEAN DEFAULT FALSE,
    verification_code VARCHAR(10),
    verified_at TIMESTAMPTZ,
    verified_by_phone VARCHAR(20),

    -- Payment
    payment_method VARCHAR(50),
    payment_reference VARCHAR(100),
    paid_at TIMESTAMPTZ,

    -- Fulfillment
    fulfilled_at TIMESTAMPTZ,
    fulfillment_notes TEXT,
    proof_of_delivery TEXT,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Cancellation
    cancelled_at TIMESTAMPTZ,
    cancellation_reason TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =============================================================================
-- SAVINGS CIRCLES (Ajo/Esusu)
-- =============================================================================

CREATE TABLE savings_circles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    creator_id UUID NOT NULL REFERENCES users(id),

    name VARCHAR(100) NOT NULL,
    description TEXT,

    -- Configuration
    contribution_amount DECIMAL(20, 4) NOT NULL,
    currency currency_code DEFAULT 'NGN',
    contribution_frequency recurrence_type NOT NULL,
    max_members INT NOT NULL,

    -- Status
    status savings_circle_status DEFAULT 'forming',

    -- Schedule
    start_date DATE NOT NULL,
    next_contribution_date DATE,
    current_cycle INT DEFAULT 0,
    total_cycles INT,

    -- Rules
    late_fee_percentage DECIMAL(5, 2) DEFAULT 5.00,
    allow_early_withdrawal BOOLEAN DEFAULT FALSE,

    -- Stats
    total_contributed DECIMAL(20, 4) DEFAULT 0,
    total_disbursed DECIMAL(20, 4) DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE circle_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    circle_id UUID NOT NULL REFERENCES savings_circles(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id),

    payout_position INT NOT NULL,
    status circle_member_status DEFAULT 'invited',

    -- Contributions
    total_contributed DECIMAL(20, 4) DEFAULT 0,
    missed_contributions INT DEFAULT 0,
    late_contributions INT DEFAULT 0,

    -- Payout
    payout_received BOOLEAN DEFAULT FALSE,
    payout_amount DECIMAL(20, 4),
    payout_date DATE,

    joined_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(circle_id, user_id),
    UNIQUE(circle_id, payout_position)
);

-- =============================================================================
-- NOTIFICATIONS
-- =============================================================================

CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    type notification_type NOT NULL,
    status notification_status DEFAULT 'pending',

    title VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    data JSONB DEFAULT '{}',

    -- Delivery
    channel VARCHAR(50), -- 'fcm', 'apns', 'twilio', 'sendgrid'
    external_id VARCHAR(100),

    -- Status
    sent_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    read_at TIMESTAMPTZ,
    failed_at TIMESTAMPTZ,
    failure_reason TEXT,

    -- Scheduling
    scheduled_for TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =============================================================================
-- AUDIT LOGS
-- =============================================================================

CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),

    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(100) NOT NULL,
    entity_id UUID,

    old_values JSONB,
    new_values JSONB,

    ip_address INET,
    user_agent TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =============================================================================
-- FX RATES CACHE
-- =============================================================================

CREATE TABLE fx_rates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    source_currency currency_code NOT NULL,
    target_currency currency_code NOT NULL,

    mid_rate DECIMAL(20, 8) NOT NULL,
    buy_rate DECIMAL(20, 8) NOT NULL,
    sell_rate DECIMAL(20, 8) NOT NULL,
    spread_bps INT DEFAULT 0,

    provider VARCHAR(50),
    fetched_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(source_currency, target_currency)
);

-- =============================================================================
-- FX QUOTES
-- =============================================================================

CREATE TABLE fx_quotes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    source_currency currency_code NOT NULL,
    target_currency currency_code NOT NULL,

    mid_rate DECIMAL(20, 8) NOT NULL,
    buy_rate DECIMAL(20, 8) NOT NULL,
    sell_rate DECIMAL(20, 8) NOT NULL,
    spread_bps INT NOT NULL,

    source_amount DECIMAL(20, 4) NOT NULL,
    target_amount DECIMAL(20, 4) NOT NULL,
    fee DECIMAL(20, 4) DEFAULT 0,
    total_source DECIMAL(20, 4) NOT NULL,

    valid_until TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    used_for_id UUID,
    used_for_type VARCHAR(50),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =============================================================================
-- REVIEWS & RATINGS
-- =============================================================================

CREATE TABLE reviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    reviewer_id UUID NOT NULL REFERENCES users(id),
    reviewed_id UUID NOT NULL REFERENCES users(id),
    booking_id UUID REFERENCES bookings(id),

    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    title VARCHAR(255),
    comment TEXT,

    -- Response
    response TEXT,
    responded_at TIMESTAMPTZ,

    -- Moderation
    is_visible BOOLEAN DEFAULT TRUE,
    reported BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =============================================================================
-- INDEXES
-- =============================================================================

-- Users
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_status ON users(status);

-- Profiles
CREATE INDEX idx_profiles_user_id ON profiles(user_id);
CREATE INDEX idx_profiles_country_of_residence ON profiles(country_of_residence);
CREATE INDEX idx_profiles_diaspora_region ON profiles(diaspora_region);

-- Wallets
CREATE INDEX idx_wallets_user_id ON wallets(user_id);
CREATE INDEX idx_wallets_user_type_currency ON wallets(user_id, wallet_type, currency);

-- Transactions
CREATE INDEX idx_transactions_wallet_id ON transactions(wallet_id);
CREATE INDEX idx_transactions_reference ON transactions(reference);
CREATE INDEX idx_transactions_status ON transactions(status);
CREATE INDEX idx_transactions_created_at ON transactions(created_at);

-- Services
CREATE INDEX idx_services_provider_id ON services(provider_id);
CREATE INDEX idx_services_category_id ON services(category_id);
CREATE INDEX idx_services_status ON services(status);
CREATE INDEX idx_services_service_area ON services USING GIN(service_area);

-- Bookings
CREATE INDEX idx_bookings_consumer_id ON bookings(consumer_id);
CREATE INDEX idx_bookings_provider_id ON bookings(provider_id);
CREATE INDEX idx_bookings_service_id ON bookings(service_id);
CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_bookings_scheduled_date ON bookings(scheduled_date);

-- Beneficiaries
CREATE INDEX idx_beneficiaries_user_id ON beneficiaries(user_id);
CREATE INDEX idx_beneficiaries_status ON beneficiaries(status);
CREATE INDEX idx_beneficiaries_phone_primary ON beneficiaries(phone_primary);

-- Remittances
CREATE INDEX idx_remittances_user_id ON remittances(user_id);
CREATE INDEX idx_remittances_beneficiary_id ON remittances(beneficiary_id);
CREATE INDEX idx_remittances_status ON remittances(status);
CREATE INDEX idx_remittances_reference ON remittances(reference);
CREATE INDEX idx_remittances_created_at ON remittances(created_at);

-- Diaspora Bookings
CREATE INDEX idx_diaspora_bookings_user_id ON diaspora_bookings(user_id);
CREATE INDEX idx_diaspora_bookings_status ON diaspora_bookings(status);
CREATE INDEX idx_diaspora_bookings_reference ON diaspora_bookings(reference);

-- Savings Circles
CREATE INDEX idx_savings_circles_creator_id ON savings_circles(creator_id);
CREATE INDEX idx_savings_circles_status ON savings_circles(status);
CREATE INDEX idx_circle_members_circle_id ON circle_members(circle_id);
CREATE INDEX idx_circle_members_user_id ON circle_members(user_id);

-- Notifications
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_status ON notifications(status);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);

-- Audit Logs
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);

-- FX
CREATE INDEX idx_fx_rates_currencies ON fx_rates(source_currency, target_currency);
CREATE INDEX idx_fx_quotes_valid_until ON fx_quotes(valid_until);

-- =============================================================================
-- TRIGGERS
-- =============================================================================

-- Updated at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply updated_at trigger to all tables
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_profiles_updated_at BEFORE UPDATE ON profiles FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_wallets_updated_at BEFORE UPDATE ON wallets FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_multi_currency_wallets_updated_at BEFORE UPDATE ON multi_currency_wallets FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_currency_balances_updated_at BEFORE UPDATE ON currency_balances FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_transactions_updated_at BEFORE UPDATE ON transactions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_service_categories_updated_at BEFORE UPDATE ON service_categories FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_services_updated_at BEFORE UPDATE ON services FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_bookings_updated_at BEFORE UPDATE ON bookings FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_beneficiaries_updated_at BEFORE UPDATE ON beneficiaries FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_remittances_updated_at BEFORE UPDATE ON remittances FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_diaspora_bookings_updated_at BEFORE UPDATE ON diaspora_bookings FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_savings_circles_updated_at BEFORE UPDATE ON savings_circles FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_circle_members_updated_at BEFORE UPDATE ON circle_members FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_notifications_updated_at BEFORE UPDATE ON notifications FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_reviews_updated_at BEFORE UPDATE ON reviews FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =============================================================================
-- SEED DATA
-- =============================================================================

-- Insert default service categories
INSERT INTO service_categories (name, slug, description, icon, sort_order) VALUES
('Home Services', 'home-services', 'Services for your home', 'home', 1),
('Personal Care', 'personal-care', 'Beauty and wellness services', 'spa', 2),
('Professional Services', 'professional-services', 'Business and professional services', 'briefcase', 3),
('Technology', 'technology', 'Tech support and digital services', 'laptop', 4),
('Events', 'events', 'Event planning and services', 'calendar', 5),
('Education', 'education', 'Tutoring and educational services', 'school', 6),
('Health & Fitness', 'health-fitness', 'Health and fitness services', 'heart', 7),
('Transportation', 'transportation', 'Delivery and transport services', 'truck', 8);

-- Insert sub-categories
INSERT INTO service_categories (name, slug, description, parent_id, sort_order) VALUES
('Cleaning', 'cleaning', 'House cleaning services', (SELECT id FROM service_categories WHERE slug = 'home-services'), 1),
('Plumbing', 'plumbing', 'Plumbing repairs and installation', (SELECT id FROM service_categories WHERE slug = 'home-services'), 2),
('Electrical', 'electrical', 'Electrical repairs and installation', (SELECT id FROM service_categories WHERE slug = 'home-services'), 3),
('Hair Styling', 'hair-styling', 'Hair styling and barbing', (SELECT id FROM service_categories WHERE slug = 'personal-care'), 1),
('Makeup', 'makeup', 'Makeup and beauty services', (SELECT id FROM service_categories WHERE slug = 'personal-care'), 2),
('Catering', 'catering', 'Food and catering services', (SELECT id FROM service_categories WHERE slug = 'events'), 1),
('Photography', 'photography', 'Photography and videography', (SELECT id FROM service_categories WHERE slug = 'events'), 2);
