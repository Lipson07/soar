DROP TRIGGER IF EXISTS update_security_settings_updated_at ON security_settings;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS security_audit_logs;
DROP TABLE IF EXISTS user_sessions;
DROP TABLE IF EXISTS security_settings;