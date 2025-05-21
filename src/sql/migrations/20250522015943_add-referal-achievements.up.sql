CREATE TABLE referral_achievements (
    id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
    identity_id UUID NOT NULL CONSTRAINT fk_identity REFERENCES identities (id),
    type TEXT NOT NULL,
    meta JSONB,
    created_at TIMESTAMP DEFAULT NOW()
)
