-- HustleX Pro Database Schema Rollback

-- Drop triggers
DROP TRIGGER IF EXISTS update_reviews_updated_at ON reviews;
DROP TRIGGER IF EXISTS update_notifications_updated_at ON notifications;
DROP TRIGGER IF EXISTS update_circle_members_updated_at ON circle_members;
DROP TRIGGER IF EXISTS update_savings_circles_updated_at ON savings_circles;
DROP TRIGGER IF EXISTS update_diaspora_bookings_updated_at ON diaspora_bookings;
DROP TRIGGER IF EXISTS update_remittances_updated_at ON remittances;
DROP TRIGGER IF EXISTS update_beneficiaries_updated_at ON beneficiaries;
DROP TRIGGER IF EXISTS update_bookings_updated_at ON bookings;
DROP TRIGGER IF EXISTS update_services_updated_at ON services;
DROP TRIGGER IF EXISTS update_service_categories_updated_at ON service_categories;
DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions;
DROP TRIGGER IF EXISTS update_currency_balances_updated_at ON currency_balances;
DROP TRIGGER IF EXISTS update_multi_currency_wallets_updated_at ON multi_currency_wallets;
DROP TRIGGER IF EXISTS update_wallets_updated_at ON wallets;
DROP TRIGGER IF EXISTS update_profiles_updated_at ON profiles;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS fx_quotes;
DROP TABLE IF EXISTS fx_rates;
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS circle_members;
DROP TABLE IF EXISTS savings_circles;
DROP TABLE IF EXISTS diaspora_bookings;
DROP TABLE IF EXISTS remittances;
DROP TABLE IF EXISTS beneficiaries;
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS services;
DROP TABLE IF EXISTS service_categories;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS currency_balances;
DROP TABLE IF EXISTS multi_currency_wallets;
DROP TABLE IF EXISTS wallets;
DROP TABLE IF EXISTS profiles;
DROP TABLE IF EXISTS users;

-- Drop enums
DROP TYPE IF EXISTS notification_status;
DROP TYPE IF EXISTS notification_type;
DROP TYPE IF EXISTS circle_member_status;
DROP TYPE IF EXISTS savings_circle_status;
DROP TYPE IF EXISTS beneficiary_status;
DROP TYPE IF EXISTS beneficiary_type;
DROP TYPE IF EXISTS recurrence_type;
DROP TYPE IF EXISTS delivery_method;
DROP TYPE IF EXISTS remittance_purpose;
DROP TYPE IF EXISTS remittance_status;
DROP TYPE IF EXISTS payment_method;
DROP TYPE IF EXISTS booking_status;
DROP TYPE IF EXISTS service_status;
DROP TYPE IF EXISTS currency_code;
DROP TYPE IF EXISTS transaction_status;
DROP TYPE IF EXISTS transaction_type;
DROP TYPE IF EXISTS wallet_status;
DROP TYPE IF EXISTS wallet_type;
DROP TYPE IF EXISTS kyc_level;
DROP TYPE IF EXISTS verification_status;
DROP TYPE IF EXISTS user_role;
DROP TYPE IF EXISTS user_status;
