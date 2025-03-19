DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_enum 
        WHERE enumtypid = 'status_type'::regtype 
        AND enumlabel = 'NOT_ACTIVE'
    ) THEN
        ALTER TYPE status_type ADD VALUE 'NOT_ACTIVE';
    END IF;
END $$;

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

-- Redefine sync_identities to exclude status
CREATE OR REPLACE FUNCTION sync_identities()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_TABLE_NAME = 'users' THEN
        INSERT INTO identities (id, type, meta, created_at, updated_at)
        VALUES (NEW.id, 'users', jsonb_build_object(
            'username', NEW.username,
            'first_name', NEW.first_name,
            'last_name', NEW.last_name,
            'email', NEW.email,
            'city', NEW.city,
            'country', NEW.country,
            'address', NEW.address,
            'avatar', NEW.avatar_id,
            'cover', NEW.cover_id
        ), NOW(), NOW())
        ON CONFLICT (id) DO UPDATE
        SET meta = EXCLUDED.meta, updated_at = NOW();
    ELSIF TG_TABLE_NAME = 'organizations' THEN
        INSERT INTO identities (id, type, meta, created_at, updated_at)
        VALUES (NEW.id, 'organizations', jsonb_build_object(
            'shortname', NEW.shortname,
            'name', NEW.name,
            'bio', NEW.bio,
            'description', NEW.description,
            'email', NEW.email,
            'phone', NEW.phone,
            'city', NEW.city,
            'country', NEW.country,
            'address', NEW.address,
            'website', NEW.website,
            'mission', NEW.mission,
            'culture', NEW.culture,
            'logo', NEW.logo_id,
            'cover', NEW.cover_id,
            'verified_impact', NEW.verified_impact,
            'verified', NEW.verified
        ), NOW(), NOW())
        ON CONFLICT (id) DO UPDATE
        SET meta = EXCLUDED.meta, updated_at = NOW();
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

UPDATE organizations SET id=id;