-- Unique Impact Points
CREATE EXTENSION IF NOT EXISTS pgcrypto;

ALTER TABLE impact_points ADD COLUMN unique_tag TEXT;

UPDATE impact_points
SET unique_tag = 'tag_' || gen_random_uuid() || '_' || encode(gen_random_bytes(24), 'base64')
WHERE unique_tag IS NULL;

ALTER TABLE impact_points ADD CONSTRAINT unique_impact_tag UNIQUE (unique_tag);
ALTER TABLE impact_points ALTER COLUMN unique_tag SET NOT NULL;

-- Unique Referral Achievements
CREATE UNIQUE INDEX idx_unique_achievements  ON referral_achievements (referrer_id, referee_id, achievement_type);