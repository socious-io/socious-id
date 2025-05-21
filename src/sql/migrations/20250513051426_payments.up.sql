CREATE TABLE wallets (
  id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
  identity_id UUID NOT NULL REFERENCES identities(id) ON DELETE CASCADE,
  chain text NOT NULL, --TODO: Define enum
  chain_id text,
  address text NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  
  UNIQUE(identity_id, chain)
);

-- CREATE TABLE cards (
--   id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
--   identity_id UUID NOT NULL REFERENCES identities(id) ON DELETE CASCADE,
--   customer text,
--   card text,
--   created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
--   updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
-- );

ALTER TABLE users ADD COLUMN stripe_customer_id VARCHAR(255);
ALTER TABLE organizations ADD COLUMN stripe_customer_id VARCHAR(255);



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
            'cover', NEW.cover_id,
            'stripe_customer_id', NEW.stripe_customer_id
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
            'verified', NEW.verified,
            'stripe_customer_id', NEW.stripe_customer_id
        ), NOW(), NOW())
        ON CONFLICT (id) DO UPDATE
        SET meta = EXCLUDED.meta, updated_at = NOW();
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

UPDATE users SET id=id;
UPDATE organizations SET id=id;