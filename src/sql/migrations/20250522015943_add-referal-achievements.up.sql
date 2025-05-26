CREATE TABLE referral_achievements (
    id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
    referrer_id UUID NOT NULL CONSTRAINT fk_referrer_identity REFERENCES identities (id),
    referee_id UUID NOT NULL CONSTRAINT fk_referee_identity REFERENCES identities (id),
    achievement_type TEXT NOT NULL,
    meta JSONB,
    created_at TIMESTAMP DEFAULT NOW()
)
