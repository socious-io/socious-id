ALTER TABLE verification_credentials
ALTER COLUMN body TYPE jsonb USING body::jsonb;