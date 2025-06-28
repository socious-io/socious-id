ALTER TABLE verification_credentials RENAME TO credentials;
CREATE TYPE credential_type AS ENUM (
    'BADGES',
    'KYC'
);
ALTER TABLE credentials ADD COLUMN type credential_type NOT NULL DEFAULT 'KYC';
ALTER TABLE impact_points ADD COLUMN claimed_at TIMESTAMP DEFAULT NULL;