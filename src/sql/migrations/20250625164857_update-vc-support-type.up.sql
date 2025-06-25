CREATE TYPE verification_credential_type AS ENUM (
    'BADGES',
    'KYC',
    'EDUCATION',
    'EXPERIENCE'
);
ALTER TABLE verification_credentials ADD COLUMN type verification_credential_type NOT NULL DEFAULT 'KYC';