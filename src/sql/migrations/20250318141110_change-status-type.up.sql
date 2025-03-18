ALTER TYPE status_type RENAME TO status_type_old2;
CREATE TYPE status_type AS ENUM (
    'ACTIVE',
    'NOT_ACTIVE',
    'SUSPENDED'
);

ALTER TABLE organizations
    ALTER COLUMN status DROP DEFAULT,
    ALTER COLUMN status TYPE status_type USING status::text::status_type,
    ALTER COLUMN status SET DEFAULT 'NOT_ACTIVE';

-- Synchronize Verification With Status
CREATE OR REPLACE FUNCTION sync_organization_verification()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.verified IS TRUE OR NEW.verified_impact IS TRUE THEN
        NEW.status = 'ACTIVE';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER trigger_sync_organization_verification
BEFORE INSERT OR UPDATE ON organizations
FOR EACH ROW EXECUTE FUNCTION sync_organization_verification();

UPDATE organizations SET id=id;