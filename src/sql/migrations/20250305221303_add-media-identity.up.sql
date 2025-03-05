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
            'avatar', NEW.avatar_id,
            'cover', NEW.cover_id
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

-- Media
CREATE TABLE media (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    identity_id uuid NOT NULL REFERENCES identities(id) ON DELETE CASCADE,
    filename text,
    url text,
    created_at timestamp with time zone NOT NULL DEFAULT now()
);