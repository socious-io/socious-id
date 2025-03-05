-- Media
CREATE TABLE media (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    identity_id uuid NOT NULL REFERENCES identities(id) ON DELETE CASCADE,
    filename text,
    url text,
    created_at timestamp with time zone NOT NULL DEFAULT now()
);

-- Organizations
CREATE TABLE organizations (
    id UUID DEFAULT public.uuid_generate_v4() PRIMARY KEY,
    shortname text NOT NULL,
    name text,
    bio text,
    description text,
    email text,
    phone text,
    
    city text,
    country text,
    address text,
    website text,
    
    mission text,
    culture text,
    
    logo jsonb,
    cover jsonb,
    
    status status_type DEFAULT 'INACTIVE',
    
    verified_impact boolean NOT NULL DEFAULT FALSE,
    verified boolean NOT NULL DEFAULT FALSE,
    
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone NOT NULL DEFAULT NOW()
);

CREATE TABLE org_members (
    id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
    user_id uuid NOT NULL,
    org_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_org FOREIGN KEY (org_id) REFERENCES organizations(id) ON DELETE CASCADE
);

-- Identity
CREATE TYPE identity_type AS ENUM (
    'users',
    'organizations'
);

CREATE TABLE identities (
    id uuid PRIMARY KEY,
    type identity_type NOT NULL,
    meta jsonb,
    created_at timestamp  NOT NULL DEFAULT NOW(),
    updated_at timestamp  NOT NULL DEFAULT NOW()
);

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
            'avatar', NEW.avatar,
            'cover', NEW.cover
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
            'logo', NEW.logo,
            'cover', NEW.cover,
            'status', NEW.status,
            'verified_impact', NEW.verified_impact,
            'verified', NEW.verified
        ), NOW(), NOW())
        ON CONFLICT (id) DO UPDATE
        SET meta = EXCLUDED.meta, updated_at = NOW();
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER trigger_users_sync
AFTER INSERT OR UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION sync_identities();

CREATE OR REPLACE TRIGGER trigger_organizations_sync
AFTER INSERT OR UPDATE ON organizations
FOR EACH ROW EXECUTE FUNCTION sync_identities();

-- Triggering the updates
UPDATE users SET id=id;
UPDATE organizations SET id=id;