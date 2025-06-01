ALTER TYPE status_type_old2 RENAME TO user_status_type;
ALTER TYPE status_type RENAME TO organization_status_type;

CREATE TABLE oauth_connects (
    id UUID PRIMARY KEY,
    identity_id UUID REFERENCES identities(id),
    status user_status_type NOT NULL DEFAULT 'NOT_ACTIVE',
    provider TEXT NOT NULL,
    matrix_unique_id TEXT NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT,
    redirect_url TEXT,
    is_confirmed BOOLEAN DEFAULT TRUE,
    meta JSONB,
    expired_at TIMESTAMP,
    created_at TIMESTAMP  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP  NOT NULL DEFAULT NOW(),
    UNIQUE (matrix_unique_id, provider)
);